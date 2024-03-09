package handlers

import (
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handlers) homeGet(c *gin.Context) {
	userGuid := c.GetString(userGUID)

	//Get user name for show on home page
	userName, err := h.service.Authorization.GetUserName(userGuid)
	if err != nil {
		logrus.Errorf("Cant get uses name. Err:%s", err.Error())
	}
	str := ""

	if userGuid != "" {
		str = fmt.Sprintf("Hello %s", userName)
	} else {
		str = "Please log in or registration"
	}

	data := map[string]string{
		"UserText": str,
	}

	tmpl, _ := template.ParseFiles("../templates/home.html")
	tmpl.Execute(c.Writer, data)
}
