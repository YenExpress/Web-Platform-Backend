package auth

import (
	"YenExpress/config"
	"YenExpress/service/patient/guard"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterPatient godoc
// @Summary      Create user account for patient
// @Description  save user details to database
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      201  {object}  DefaultResponse
// @Failure      400  {object}  DefaultResponse
// @Failure      409  {object} 	DefaultResponse
// @Failure      500  {object}  DefaultResponse
// @Router       /patient/auth/register/ [post]
func Register(c *gin.Context) {

	func() {

		var input *CreatePatientDTO

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, DefaultResponse{Message: err.Error()})
			return
		}

		var user *Patient
		err := config.DB.Where("email = ?", input.Email).First(&user).Error
		if err == nil {
			c.JSON(http.StatusConflict, DefaultResponse{Message: "User Account Already Exists"})
			return
		}

		user = &Patient{
			FirstName: input.FirstName, LastName: input.LastName,
			Email: input.Email, Password: input.Password,
			Sex:      input.Sex,
			Location: input.Location, PhoneNumber: input.PhoneNumber,
		}

		_, err = user.SaveNew()

		if err != nil {
			c.JSON(http.StatusInternalServerError, DefaultResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, DefaultResponse{Message: "User Account created"})
		return

	}()
}

// LoginPatient godoc
// @Summary      Enable sign in and authorization for patient with valid credentials
// @Description  Validate patient credentials, authenticate and authorize with JWT provision
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  LoginResponse
// @Failure      400  {object}  DefaultResponse
// @Failure      401  {object}  DefaultResponse
// @Failure      403  {object} 	DefaultResponse
// @Failure      429  {object}  DefaultResponse
// @Failure      409  {object}  DefaultResponse
// @Router        /patient/auth/login/ [post]
func Login(c *gin.Context) {

	credentials, allowedToLogin := RateLimitLogin(c)

	if allowedToLogin {
		func() {

			var user *Patient
			err := config.DB.Where("Email = ?", credentials.Email).First(&user).Error
			if err != nil {
				c.JSON(http.StatusUnauthorized, DefaultResponse{Message: "Incorrect Email"})
				return
			}

			err = user.ValidatePwd(credentials.Password)
			if err != nil {
				guard.LoginLimiter.UpdateRequest(credentials.Email, credentials.IPAddress)
				c.JSON(http.StatusForbidden, DefaultResponse{Message: "Incorrect Password"})
				return
			}
			if PatientLoginManager.checkActiveSession(user.ID) {
				PatientLoginManager.authConcurrentSignin(user.Email)
				c.JSON(http.StatusConflict, DefaultResponse{Message: "One Time Password for sign in sent to user email address to enable account concurrency!!"})
				return
			}
			sessionID := PatientLoginManager.createNewSession(*user, credentials.IPAddress)
			IDToken, AccessToken, RefreshToken := PatientLoginManager.generateAuthTokens(*user, sessionID)

			c.JSON(http.StatusOK, LoginResponse{
				IDToken, AccessToken, RefreshToken, "Bearer"})
			return

		}()
	}
}

// SendOTPToPatientMail godoc
// @Summary      initiate email validation for patient sign up or login concurrency
// @Description  Send OTP to specified patient email address for authentication and registration
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      202  {object}  DefaultResponse
// @Failure      401  {object}  DefaultResponse
// @Failure      429  {object}  DefaultResponse
// @Router        /patient/auth/confirm-email/process/ [post]
func GenerateOTPforAuth(c *gin.Context) {

	cred, allowedToGenerate := RateLimitOTPGeneration(c)
	if allowedToGenerate {

		func() {

			if cred.process == "signup" {
				var user *Patient
				err := config.DB.Where("email = ?", cred.email).First(&user).Error
				if err == nil {
					c.JSON(http.StatusConflict, DefaultResponse{Message: "User Account Already Exists"})
					return
				}
				PatientRegistrationManager.authNewEmail(cred.email)
				guard.CreateOTPLimiter.UpdateRequest(cred.email, cred.ipAddress)
				c.JSON(http.StatusAccepted, DefaultResponse{Message: "OTP sent to mail address provided to confirm ownership and proceed with registration"})
				return

			} else if cred.process == "signin" {
				PatientLoginManager.authConcurrentSignin(cred.email)
				guard.CreateOTPLimiter.UpdateRequest(cred.email, cred.ipAddress)
				c.JSON(http.StatusAccepted, DefaultResponse{Message: "One Time Password for sign in sent to user email address to enable account concurrency!!"})
				return
			}
		}()
	}
}

// ValidatePatientWithOTPSentToMail godoc
// @Summary      Enable patient sign up or login concurrency via validation
// @Description  validate OTP to specified patient email address for authentication and registration
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  DefaultResponse
// @Success      200  {object}  LoginResponse
// @Failure      401  {object}  DefaultResponse
// @Failure      429  {object}  DefaultResponse
// @Router        /patient/auth/verify-otp/process/ [post]
func ValidateOneTimePass(c *gin.Context) {

	credentials, allowedToValidate := RateLimitOTPValidation(c)
	if allowedToValidate {
		func() {
			if credentials.process == "signup" {
				var user *Patient
				err := config.DB.Where("email = ?", credentials.email).First(&user).Error
				if err == nil {
					c.JSON(http.StatusConflict, DefaultResponse{Message: "User Account Already Exists"})
					return
				}
				err = PatientRegistrationManager.enableSignUpwithMail(credentials.email, credentials.otp)
				if err != nil {
					guard.EmailValidationLimiter.UpdateRequest(credentials.email, credentials.ipAddress)
					c.JSON(http.StatusUnauthorized, DefaultResponse{Message: err.Error()})
					return
				}
				c.JSON(http.StatusOK, DefaultResponse{Message: "Email Validation Successful"})
				return

			} else if credentials.process == "signin" {
				var user *Patient
				err := config.DB.Where("email = ?", credentials.email).First(&user).Error
				if err != nil {
					c.JSON(http.StatusUnauthorized, DefaultResponse{Message: "User Unknown"})
					return
				}
				idTok, accTok, refTok, err := PatientLoginManager.enableConcurrentSignin(*user, credentials.ipAddress, credentials.otp)
				if err != nil {
					guard.EmailValidationLimiter.UpdateRequest(credentials.email, credentials.ipAddress)
					c.JSON(http.StatusUnauthorized, DefaultResponse{Message: err.Error()})
					return
				}
				c.JSON(http.StatusOK, LoginResponse{idTok, accTok, refTok, "Bearer"})
				return
			}

		}()
	}

}

// LogoutPatient godoc
// @Summary      Enable sign out and session delete for patient with valid credentials
// @Description  Log patient out with server wipe of session data
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  LoginResponse
// @Failure      401  {object}  DefaultResponse
// @Failure      403  {object} 	DefaultResponse
// @Router        /patient/auth/logout/ [delete]
func Logout(c *gin.Context) {

	func() {

		token := c.Request.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		load, _ := guard.GetTokenPayload(token, "access_token")
		PatientLoginManager.endSession(load.UserId, load.SessionID)

		c.JSON(http.StatusOK, DefaultResponse{Message: "Account Successfully Logged out"})
		return

	}()
}
