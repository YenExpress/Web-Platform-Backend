package controllers

import (
	"YenExpress/config"
	"YenExpress/service/dto"
	"YenExpress/service/guard"
	mid "YenExpress/service/patient/middlewares"
	"YenExpress/service/patient/models"
	pro "YenExpress/service/patient/providers"

	"time"

	"net/http"

	"github.com/gin-gonic/gin"
)

// Create user account for patient
// Save user details to database
func Register(c *gin.Context) {

	func() {

		var input *models.RegistrationDTO

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, dto.DefaultResponse{Message: err.Error()})
			return
		}

		var user *dto.Patient
		err := config.DB.Where("email = ?", input.Email).First(&user).Error
		if err == nil {
			c.JSON(http.StatusConflict, dto.DefaultResponse{Message: "User Account Already Exists"})
			return
		}

		user = &dto.Patient{
			FirstName: input.FirstName, LastName: input.LastName,
			Email: input.Email, Password: input.Password,
			Sex:      input.Sex,
			Location: input.Location, PhoneNumber: input.PhoneNumber,
		}

		_, err = user.SaveNew()

		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.DefaultResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, dto.DefaultResponse{Message: "User Account created"})
		return

	}()
}

// Enable sign in and authorization for patient with valid credentials
// Validate patient credentials, authenticate and authorize with JWT provision
func Login(c *gin.Context) {

	func() {

		var credentials *dto.LoginCredentials
		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, dto.DefaultResponse{Message: err.Error()})
		}

		var user *dto.Patient
		err := config.DB.Where("Email = ?", credentials.Email).First(&user).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.DefaultResponse{Message: "Incorrect Email"})
			return
		}

		err = user.ValidatePwd(credentials.Password)
		if err != nil {
			c.JSON(http.StatusForbidden, dto.DefaultResponse{Message: "Incorrect Password"})
			return
		}

		IDToken, AccessToken, RefreshToken := pro.LoginService.GenerateAuthTokens(*user)

		c.JSON(http.StatusOK, dto.LoginResponse{
			IDToken: IDToken, AccessToken: AccessToken,
			RefreshToken: RefreshToken, TokenType: "Bearer"})
		return

	}()
}

// Initiate email validation for patient sign up
// Send OTP to specified patient email address for new account registration
func SendNewMailOTP(c *gin.Context) {

	cred, allowedToGenerate := mid.RateLimitOTPGeneration(c)
	if allowedToGenerate {

		func() {

			var user *dto.Patient
			err := config.DB.Where("email = ?", cred.Email).First(&user).Error
			if err == nil {
				c.JSON(http.StatusConflict, dto.DefaultResponse{Message: "User Account Already Exists"})
				return
			}
			pro.RegistrationService.AuthNewEmail(cred.Email)
			pro.CreateOTPLimiter.UpdateRequest(cred.Email, cred.IPAddress)
			c.JSON(http.StatusAccepted, dto.DefaultResponse{Message: "OTP sent to mail address provided to confirm ownership and proceed with registration"})
			return

		}()
	}
}

// Enable patient sign up via email validation with OTP
// Verify OTP sent to specified patient email address for new account registration
func ConfirmNewMail(c *gin.Context) {

	credentials, allowedToValidate := mid.RateLimitOTPValidation(c)
	if allowedToValidate {
		func() {

			var user *dto.Patient
			err := config.DB.Where("email = ?", credentials.Email).First(&user).Error
			if err == nil {
				c.JSON(http.StatusConflict, dto.DefaultResponse{Message: "User Account Already Exists"})
				return
			}
			err = pro.RegistrationService.EnableSignUpwithMail(credentials.Email, credentials.OTP)
			if err != nil {
				pro.EmailValidationLimiter.UpdateRequest(credentials.Email, credentials.IPAddress)
				c.JSON(http.StatusUnauthorized, dto.DefaultResponse{Message: err.Error()})
				return
			}
			c.JSON(http.StatusOK, dto.DefaultResponse{Message: "Email Validation Successful"})
			return
		}()
	}

}

// Refresh Expired Access Token
// Create New Access Token for Patient Authentication with Refresh Token
func Refresh(c *gin.Context) {

	func() {

		user_id, session_id, err := GetIDsFromRequest(c, "refresh_token")
		if err != nil {
			c.JSON(500, dto.DefaultResponse{Message: err.Error()})
			return
		}
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

// Enable sign out for patient with valid credentials
// Log patient out with server
func Logout(c *gin.Context) {

	func() {

		c.JSON(http.StatusOK, dto.DefaultResponse{Message: "Account Successfully Logged out"})
		return

	}()
}
