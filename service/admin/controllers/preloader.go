package controllers

import (
	pro "YenExpress/service/admin/providers"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetIDsFromRequest(c *gin.Context, variety string) (user_id uint, session_id string) {
	authCred := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(authCred, "Bearer ")
	load, err := pro.JWTMaker.GetPayloadFromToken(token, variety)
	if err != nil {

		return 0, ""
	} else {
		user_id, session_id = load.UserId, load.SessionID
	}
	return
}
