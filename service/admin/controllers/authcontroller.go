package controllers

import (
	"YenExpress/config"
	mid "YenExpress/service/admin/middlewares"
	"YenExpress/service/admin/models"
	pro "YenExpress/service/admin/providers"
	"YenExpress/service/dto"
	"YenExpress/service/guard"

	"time"

	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateAdmin godoc
// @Summary      Create user account for admin
// @Description  save user details to database
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      201  {object}  dto.DefaultResponse
// @Failure      400  {object}  dto.DefaultResponse
// @Failure      409  {object} 	dto.DefaultResponse
// @Failure      500  {object}  dto.DefaultResponse
// @Failure      429  {object}  dto.DefaultResponse
// @Router       /admin/auth/create-account/ [post]
func CreateAdmin(c *gin.Context) {

	func() {

		var input *models.CreateDTO

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, dto.DefaultResponse{Message: err.Error()})
			return
		}

		var user *models.Admin
		err := config.DB.Where("email = ?", input.Email).First(&user).Error
		if err == nil {
			c.JSON(http.StatusConflict, dto.DefaultResponse{Message: "Admin Account Already Exists"})
			return
		}

		user = &models.Admin{
			FirstName: input.FirstName, LastName: input.LastName,
			Email: input.Email, Password: input.Password,
			Role: input.Role,
		}

		_, err = user.SaveNew()

		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.DefaultResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, dto.DefaultResponse{Message: "Admin Account created"})
		return

	}()
}

// LoginAdmin godoc
// @Summary      Enable sign in and authorization for Admin with valid credentials
// @Description  Validate Admin credentials, authenticate and authorize with JWT provision
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.LoginResponse
// @Failure      400  {object}  dto.DefaultResponse
// @Failure      401  {object}  dto.DefaultResponse
// @Failure      403  {object} 	dto.DefaultResponse
// @Failure      409  {object}  dto.DefaultResponse
// @Router        /admin/auth/login/ [post]
func Login(c *gin.Context) {

	credentials, allowedToLogin := mid.RateLimitLogin(c)

	if allowedToLogin {
		func() {

			var user *models.Admin
			err := config.DB.Where("Email = ?", credentials.Email).First(&user).Error
			if err != nil {
				c.JSON(http.StatusUnauthorized, dto.DefaultResponse{Message: "Incorrect Email"})
				return
			}

			err = user.ValidatePwd(credentials.Password)
			if err != nil {
				pro.LoginLimiter.UpdateRequest(credentials.Email, credentials.IPAddress)
				c.JSON(http.StatusForbidden, dto.DefaultResponse{Message: "Incorrect Password"})
				return
			}
			if pro.LoginService.CheckActiveSession(user.ID) {
				pro.LoginService.AuthConcurrentSignin(user.Email)
				c.JSON(http.StatusConflict, dto.DefaultResponse{Message: "One Time Password for sign in sent to user email address to enable account concurrency!!"})
				return
			}
			sessionID := pro.LoginService.CreateNewSession(*user, credentials.IPAddress)
			IDToken, AccessToken, RefreshToken := pro.LoginService.GenerateAuthTokens(*user, sessionID)

			c.JSON(http.StatusOK, dto.LoginResponse{
				IDToken: IDToken, AccessToken: AccessToken,
				RefreshToken: RefreshToken, TokenType: "Bearer"})
			return

		}()
	}
}

// SendOTPToAdminMail godoc
// @Summary      initiate email validation for admin login concurrency
// @Description  Send OTP to specified admin email address for authentication
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      202  {object}  dto.DefaultResponse
// @Failure      401  {object}  dto.DefaultResponse
// @Failure      429  {object}  dto.DefaultResponse
// @Router        /admin/auth/login/send-otp [post]
func GenerateAuthOTP(c *gin.Context) {

	cred, allowedToGenerate := mid.RateLimitOTPGeneration(c)
	if allowedToGenerate {

		func() {

			pro.LoginService.AuthConcurrentSignin(cred.Email)
			pro.CreateOTPLimiter.UpdateRequest(cred.Email, cred.IPAddress)
			c.JSON(http.StatusAccepted, dto.DefaultResponse{Message: "One Time Password for sign in sent to user email address to enable account concurrency!!"})
			return
		}()
	}
}

// ValidateAuthOTPSentToAdmin godoc
// @Summary      Enable admin login concurrency via otp confirmation
// @Description  validate OTP to specified admin email address for authentication
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.DefaultResponse
// @Success      200  {object}  dto.LoginResponse
// @Failure      401  {object}  dto.DefaultResponse
// @Failure      429  {object}  dto.DefaultResponse
// @Router        /admin/auth/login/verify [post]
func ValidateLoginOTP(c *gin.Context) {

	credentials, allowedToValidate := mid.RateLimitOTPValidation(c)
	if allowedToValidate {
		func() {
			var user *models.Admin
			err := config.DB.Where("email = ?", credentials.Email).First(&user).Error
			if err != nil {
				c.JSON(http.StatusUnauthorized, dto.DefaultResponse{Message: "Admin Does Not Exist"})
				return
			}
			idTok, accTok, refTok, err := pro.LoginService.EnableConcurrentSignin(*user, credentials.IPAddress, credentials.OTP)
			if err != nil {
				pro.EmailValidationLimiter.UpdateRequest(credentials.Email, credentials.IPAddress)
				c.JSON(http.StatusUnauthorized, dto.DefaultResponse{Message: err.Error()})
				return
			}
			c.JSON(http.StatusOK, dto.LoginResponse{IDToken: idTok,
				AccessToken: accTok, RefreshToken: refTok,
				TokenType: "Bearer"})
			return
		}()
	}

}

// RefreshAdminToken godoc
// @Summary      Refresh Expired Access Token
// @Description  Create New Access Token for Admin Authentication with Refresh Token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.RefreshTokenResponse
// @Failure      401  {object}  dto.DefaultResponse
// @Failure      429  {object}  dto.DefaultResponse
// @Router        /admin/auth/refresh/ [get]
func RefreshAuthToken(c *gin.Context) {

	func() {

		user_id, session_id := GetIDsFromRequest(c, "refresh_token")
		newAccessToken := pro.JWTMaker.CreateToken(
			guard.Bearer{
				UserId: user_id, SessionID: session_id,
				Expiration: time.Now().Add(time.Hour * 24 * 3),
				Issuer:     config.ServerDomain, Class: "access_token",
			})

		c.JSON(http.StatusOK, dto.RefreshTokenResponse{AccessToken: newAccessToken})
		return

	}()
}

// LogoutAdmin godoc
// @Summary      Enable sign out and session delete for admin with valid credentials
// @Description  Log patient out with server wipe of session data
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.LoginResponse
// @Failure      401  {object}  dto.DefaultResponse
// @Failure      429  {object}  dto.DefaultResponse
// @Router        /admin/auth/logout/ [delete]
func Logout(c *gin.Context) {

	func() {

		user_id, session_id := GetIDsFromRequest(c, "access_token")
		pro.LoginService.EndSession(user_id, session_id)

		c.JSON(http.StatusOK, dto.DefaultResponse{Message: "Account Successfully Logged out"})
		return

	}()
}
