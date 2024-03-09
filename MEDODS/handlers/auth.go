package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TokensStruct struct {
	JwtToken     string
	RefreshToken string
}

// FUNCTION FOR RETURN ERROR ALERT THROW HTMX
func generateErrorAllert(errorStatus int, errorTitle string, errorMessage string, err error, c *gin.Context) {
	if err != nil {
		logrus.Infof("ErrorStatus : %d | Error Title :%s | Error Message :%s | Error :%s", errorStatus, errorTitle, errorMessage, err.Error())
	} else {
		logrus.Infof("ErrorStatus : %d | Error Title :%s | Error Message :%s | Error :nill", errorStatus, errorTitle, errorMessage)
	}

	htmlStr := fmt.Sprintf("<div class='alert alert-danger' role='alert'>%s! %s</div>", errorTitle, errorMessage)
	tmpl, _ := template.New("t").Parse(htmlStr)
	tmpl.Execute(c.Writer, nil)
}

// Post Function sign up
func (h *Handlers) signUpPost(c *gin.Context) {
	username := c.Request.PostFormValue("inputUserName")
	password := c.Request.PostFormValue("inputPassword")
	reapeatPassword := c.Request.PostFormValue("repeatPassword")

	//Check data
	if password != reapeatPassword {
		generateErrorAllert(http.StatusBadRequest, "Registration Failed", "Incorrectly entered repeated password", nil, c)
		return
	}
	if password == "" {
		generateErrorAllert(http.StatusBadRequest, "Registration Failed", "Password is empty", nil, c)
		return
	}
	if username == "" {
		generateErrorAllert(http.StatusBadRequest, "Registration Failed", "Username is empty", nil, c)
		return
	}
	if len(username) > 75 {
		generateErrorAllert(http.StatusBadRequest, "Registration Failed", "Username is to long", nil, c)
		return
	}
	if len(password) > 500 {
		generateErrorAllert(http.StatusBadRequest, "Registration Failed", "Password is to long", nil, c)
		return
	}

	_, err := h.service.Authorization.CreateUser(username, password)
	if err != nil {
		generateErrorAllert(http.StatusBadRequest, "Registration Failed", "This user is registred", err, c)
		return
	}

	htmlStr := "<div class='alert alert-success' role='alert'>Succes registration</div>"
	tmpl, _ := template.New("t").Parse(htmlStr)
	tmpl.Execute(c.Writer, nil)
	c.HTML(http.StatusOK, "", gin.H{
		"Message": "Succes registration",
	})

}

// Get Function sign up
func (h *Handlers) signUpGet(c *gin.Context) {
	c.HTML(http.StatusOK, "signUp.html", gin.H{
		"ErrorTitle":   "",
		"ErrorMessage": "",
	})
}

// Post function sign in
func (h *Handlers) signInPost(c *gin.Context) {
	username := c.PostForm("inputUserName")
	password := c.PostForm("inputPassword")

	//Check Data
	if password == "" {
		generateErrorAllert(http.StatusBadRequest, "Authorization Failed", "Password is empty", nil, c)
		return
	}
	if username == "" {
		generateErrorAllert(http.StatusBadRequest, "Authorization Failed", "Username is empty", nil, c)
		return
	}

	//Generate Tokens
	jwtToken, refreshToken, refreshTokenGUID, err := h.service.Authorization.GenerateToken(username, password)
	if err != nil {
		generateErrorAllert(http.StatusInternalServerError, "Authorization Failed", "No such user", err, c)
		return
	}

	//Prepare tokens to save in cookie
	tokensStr := fmt.Sprintf(jwtToken + "\\||\\" + refreshToken + "\\||\\" + refreshTokenGUID)

	// Empty name deletes all cookies
	c.SetCookie("", "", -1, "/", "localhost", true, true)

	//Create Http only cookie for save jwt and refresh tokens
	c.SetCookie("token", tokensStr, 3600*24*7, "/", "localhost", true, true)

	c.Header("HX-Redirect", "/")
}

// Get function sign in
func (handler *Handlers) signInGet(c *gin.Context) {
	c.HTML(http.StatusOK, "signIn.html", gin.H{
		"ErrorTitle":   "",
		"ErrorMessage": "",
	})
}
