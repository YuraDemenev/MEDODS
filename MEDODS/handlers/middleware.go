package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const userGUID = "userGUID"

func (h *Handlers) identifyUser(c *gin.Context) {
	//Get token from cookie
	tokensStr, err := c.Cookie("token")
	tokens := strings.Split(tokensStr, "\\||\\")

	if err != nil {
		logrus.Errorf("Cant get cookie token. Err:%s", err.Error())
		return
	}

	//Parse token
	guid, err := h.service.ParseToken(tokens[0])
	if err != nil {
		//If jwt token expire
		if guid == "" && err.Error() == "Token is expired" {
			h.refreshTokens(c)

		} else {
			logrus.Errorf("Cant parse token. Err:%s", err.Error())
		}
		return
	}
	c.Set(userGUID, guid)
}
