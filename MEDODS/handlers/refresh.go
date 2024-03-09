package handlers

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handlers) refreshTokens(c *gin.Context) {
	//Get token from cookie
	tokensStr, err := c.Cookie("token")
	tokens := strings.Split(tokensStr, "\\||\\")
	if err != nil {
		logrus.Errorf("Cant get cookie token. Err:%s", err.Error())
		return
	}

	jwtToken, refreshToken, refreshTokenGUID, guid, err := h.service.RefreshTokens.RefreshTokens(tokens[1], tokens[2])
	if err != nil {
		logrus.Errorf("Cant refresh token. Err:%s", err.Error())
		return
	}

	tokensStr = fmt.Sprintf(jwtToken + "\\||\\" + refreshToken + "\\||\\" + refreshTokenGUID)

	// Empty name deletes all cookies
	c.SetCookie("", "", -1, "/", "localhost", true, true)

	//Create Http only cookie for save jwt and refresh tokens
	c.SetCookie("token", tokensStr, 3600*24*7, "/", "localhost", true, true)

	c.Header("HX-Redirect", "/")
	c.Set(userGUID, guid)
}
