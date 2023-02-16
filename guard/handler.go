package guard

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidatePatientAccessToken(token string) (*Patient, error) {
	user, err := PatientJWTMaker.VerifyPatientToken(token, "access_token")
	if err != nil {
		return nil, err
	}
	return user, err
}

func ValidatePatientRefreshToken(token string) error {
	_, err := PatientJWTMaker.VerifyPatientToken(token, "refresh_token")
	return err

}

func GetPatientFromToken(c *gin.Context) (*Patient, error) {
	token := c.Request.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
		user, err := ValidatePatientAccessToken(token)
		if err == nil {

			return user, nil
		} else {
			return &Patient{}, err
		}
	}
	return &Patient{}, ErrInvalidToken
}
