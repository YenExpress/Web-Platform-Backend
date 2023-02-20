package patientauth

import (
	"YenExpress/config"
	midware "YenExpress/middleware/patientware"
	"YenExpress/ratelimiter"
	"fmt"
	"net/http"
	"sync"

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
// @Router       /auth/create/patient [post]
func RegisterPatient(c *gin.Context) {

	if midware.AuthorizePatientWithAPIKey(c) {

		func() {

			var input *CreatePatientDTO

			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, midware.DefaultResponse{Message: err.Error()})
				return
			}

			var user *Patient
			err := config.DB.Where("email = ?", input.Email).First(&user).Error
			// fmt.Println(err)
			if err == nil {
				c.JSON(http.StatusConflict, midware.DefaultResponse{Message: "User Account Already Exists"})
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
				c.JSON(http.StatusInternalServerError, midware.DefaultResponse{Message: err.Error()})
				return
			}

			c.JSON(http.StatusCreated, midware.DefaultResponse{Message: "User Account created"})
			return

		}()
	}
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
// @Failure      490  {object}  DefaultResponse
// @Router       /auth/create/patient [post]
func LoginPatient(c *gin.Context) {

	if midware.AuthorizePatientWithAPIKey(c) {

		credentials, allowedToLogin := midware.RateLimitPatientLogin(c)

		if allowedToLogin {
			func() {

				var user *Patient
				err := config.DB.Where("Email = ?", credentials.Email).First(&user).Error
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusUnauthorized, midware.DefaultResponse{Message: "Incorrect Email"})
					return
				}

				err = user.ValidatePwd(credentials.Password)
				if err != nil {
					ratelimiter.PatientLoginLimiter.NoteFailure(credentials.Email, credentials.IPAddress)
					c.JSON(http.StatusForbidden, midware.DefaultResponse{Message: "Incorrect Password"})
					return
				}
				var IDToken, AccessToken, RefreshToken string
				var wg sync.WaitGroup
				wg.Add(3)

				go func(token *string) {
					*token = user.GetIDToken()
					wg.Done()
				}(&IDToken)

				go func(token *string) {
					*token = user.GetAccessToken()
					wg.Done()
				}(&AccessToken)

				go func(token *string) {
					*token = user.GetRefreshToken()
					wg.Done()
				}(&RefreshToken)

				wg.Wait()

				c.JSON(http.StatusOK, LoginResponse{
					IDToken, AccessToken, RefreshToken, "Bearer"})
				return

			}()
		}
	}
}
