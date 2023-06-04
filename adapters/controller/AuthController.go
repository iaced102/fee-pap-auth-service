package controller

import (
	"aBet/crypt"
	"aBet/library"
	"aBet/model"
	"aBet/usecase/service"
	"bytes"
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Claims struct {
	Username string `json:"userName"`
	Id       string `json:"id"`
	UserType int    `json:"userType"`
	jwt.StandardClaims
}

var (
	signKey    *rsa.PrivateKey
	serverPort int
)

type authController struct {
	usersService service.UsersService
}

func checkExistUserName(uC *authController, userName string) bool {

	userInfo, e := uC.usersService.GetUsersByName(userName)
	fmt.Println(userInfo)
	if e != nil {
		return false
	}
	return userInfo.Id != ""
}
func convertByteToString(b []byte) string {
	s := make([]string, len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s, ",")
}
func EncryptPass(Password string, filepath string) string {
	Password = strings.ReplaceAll(Password, " ", "")
	PublicKey := crypt.GetPublicKey(filepath)
	EncryptPass := crypt.EncryptByPublicKey(PublicKey, []byte(Password))
	return EncryptPass
}
func DecryptPass(encryptPass string, filepath string) string {
	PrivateKey := crypt.GetPrivateKey(filepath)
	DecryptPass := crypt.DecryptByPrivateKey(PrivateKey, encryptPass)
	return DecryptPass
}

func GenerateToken(accountUser model.Users) (string, error) {
	secret := crypt.GetPrivateKeyByte("crypt/public3.pem") //[]byte("mysecretkey")
	fmt.Println(accountUser.Id, "iddddddddd")
	// Create a new set of claims for the token
	claims := &Claims{
		Username: accountUser.UserName,
		Id:       accountUser.Id,
		UserType: accountUser.UserType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	// Create a new token with the claims and sign it with the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secret)
	if err != nil {
		fmt.Println("Error creating token: ", err)
		return "", nil
	}

	fmt.Println("Token: ", signedToken)
	return signedToken, nil
}

func (uC *authController) LoginUserAccount(c *Context) error {
	var AccountParams model.Users
	c.Bind(&AccountParams)
	authenUserName := checkExistUserName(uC, AccountParams.UserName)
	if !authenUserName {
		return c.Output(http.StatusBadRequest, "User name not found", nil)
	}
	AccountParams.Password = library.HashStringSha256(AccountParams.Password)
	fmt.Println(AccountParams, AccountParams.Id, "id loginnnnnnnnnnnnn")
	usersInfo, e := uC.usersService.LoginUserAccount(AccountParams.UserName, AccountParams.Password)
	if e != nil || usersInfo.Id == "" {
		return c.Output(http.StatusBadRequest, "User name or password wrong", e)
	}

	token, errorToken := GenerateToken(usersInfo)
	if errorToken != nil {
		return c.Output(http.StatusBadRequest, "RenderToken Gone wrong", e)
	}

	return c.Output(http.StatusOK, map[string]interface{}{
		"message": "Login Successfully",
		"profile": usersInfo,
		"token":   token,
		"Expires": time.Now().Add(time.Minute * 59).Unix(),
	}, e)
}

// AddUsers implements SystemConfigController
func (uC *authController) AddUsers(c *Context) error {

	var AccountParams model.Users
	c.Bind(&AccountParams)
	authenUserName := checkExistUserName(uC, AccountParams.UserName)
	if authenUserName {
		return c.Output(http.StatusBadRequest, "Create account fail : duplicate user name", nil)
	}
	AccountParams.Id = uuid.NewString()
	AccountParams.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	AccountParams.CryptPassword = EncryptPass(AccountParams.Password, "crypt/pubkeyv2.pem")

	AccountParams.Password = library.HashStringSha256(AccountParams.Password)
	us, e := uC.usersService.AddUsers(AccountParams)
	if e != nil {
		return c.Output(http.StatusBadRequest, "Create account fail", e)
	}
	return c.Output(http.StatusOK, map[string]interface{}{
		"message":     "Create Account Successfully",
		"accountInfo": us,
	}, e)
}

// DeleteUsers implements SystemConfigController
func (uC *authController) DeleteUsers(c *Context) error {
	var UsersInfo model.Users
	c.Bind(&UsersInfo)
	_, e := uC.usersService.DeleteUsers(UsersInfo)
	if e != nil {
		return c.Output(http.StatusBadRequest, "Delete Account Fail", e)
	}
	return c.Output(http.StatusOK, "Ok", e)
}

// EditUsers implements SystemConfigController
func (uC *authController) UpdateUsers(c *Context) error {
	userId := c.AuthObject.GetUserId()
	currentUserInfo := model.Users{
		Id: userId,
	}
	LcurrentUserInfo, e := uC.usersService.GetByIDUsers(userId)
	currentUserInfo = LcurrentUserInfo[0]
	if e != nil {
		return c.Output(http.StatusBadRequest, nil, errors.New("can not find this user"))
	}
	AccountParams := model.Users{}
	c.Bind(&AccountParams)
	if AccountParams.Password != "" {

		AccountParams.CryptPassword = EncryptPass(AccountParams.Password, "crypt/pubkeyv2.pem")

		AccountParams.Password = library.HashStringSha256(AccountParams.Password)
	} else {
		AccountParams.Password = currentUserInfo.Password
		AccountParams.CryptPassword = currentUserInfo.CryptPassword

	}
	AccountParams.Id = currentUserInfo.Id
	AccountParams.UserType = currentUserInfo.UserType
	AccountParams.UserName = currentUserInfo.UserName
	AccountParams.CreatedAt = currentUserInfo.CreatedAt
	if AccountParams.FirstName == "" {
		AccountParams.FirstName = currentUserInfo.FirstName
	}
	if AccountParams.LastName == "" {
		AccountParams.LastName = currentUserInfo.LastName
	}
	if AccountParams.Email == "" {
		AccountParams.Email = currentUserInfo.Email
	}
	if AccountParams.CustomField == "" {
		AccountParams.CustomField = currentUserInfo.CustomField
	}

	// AccountParams.CryptPassword = EncryptPass(AccountParams.Password, "crypt/pubkeyv2.pem")

	// AccountParams.Password = library.HashStringSha256(AccountParams.Password)
	us, e := uC.usersService.EditUsers(AccountParams)
	if e != nil {
		return c.Output(http.StatusBadRequest, "Create account fail", e)
	}
	return c.Output(http.StatusOK, map[string]interface{}{
		"message":     "Create Account Successfully",
		"accountInfo": us,
	}, e)
}

func (uC *authController) TestJWT(c *Context) error {
	// var AccountParams model.Users
	// c.Bind(&AccountParams)
	// AccountParams.CryptPassword = EncryptPass(AccountParams.Password, "crypt/pubkeyv2.pem")

	// AccountParams.Password = library.HashStringSha256(AccountParams.Password)
	// us, e := uC.usersService.EditUsers(AccountParams)
	// if e != nil {
	// 	return c.Output(http.StatusBadRequest, "Create account fail", e)
	// }
	// return c.Output(http.StatusOK, map[string]interface{}{
	// 	"message":     "Create Account Successfully",
	// 	"accountInfo": us,
	// }, e)
	return c.Output(http.StatusOK, nil, nil)
}

func (uC *authController) GetDetailUsers(c *Context) error {
	var UsersAccountParams model.Users
	c.Bind(&UsersAccountParams)
	u, err := uC.usersService.GetByIDUsers(UsersAccountParams.Id)
	for i := 0; i < len(u); i++ {
		// u[i].CryptPassword = DecryptPass(u[i].CryptPassword, "crypt/privkeyv2.pem")
	}
	return c.Output(http.StatusOK, u, err)
}
func (uC *authController) ResetPassword(c *Context) error {
	// userName := c.Param("userName")
	user := model.Users{}
	c.Bind(&user)
	user, e := uC.usersService.GetUsersByName(user.UserName)
	if e != nil {
		return c.Output(http.StatusBadRequest, nil, errors.New("can not find user"))
	}
	// newUUID, _ := exec.Command("uuidgen").Output()
	newPassword := uuid.NewString()
	fmt.Println(user.Id)
	fmt.Println(newPassword)
	user.Password = newPassword
	user.CryptPassword = EncryptPass(user.Password, "crypt/pubkeyv2.pem")
	user.Password = library.HashStringSha256(user.Password)
	_, er := uC.usersService.EditUsers(user)
	if er != nil {
		c.Output(http.StatusBadRequest, nil, errors.New("can not resetpassword"))
	}
	from := os.Getenv("EMAIL")
	password := os.Getenv("APPLICATION_EMAIL_PASSWORD")

	// Receiver email address.
	to := []string{
		user.Email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	fmt.Println("", from, password, smtpHost)
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("template.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Name        string
		Newpassword string
	}{
		Name:        user.FirstName + " " + user.LastName,
		Newpassword: newPassword,
	})

	// Sending email.
	fmt.Println(smtpHost+":"+smtpPort, auth, from, to)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return c.Output(http.StatusBadRequest, nil, errors.New("can not send email"))
	}
	fmt.Println("Email Sent!")
	return c.Output(http.StatusOK, nil, errors.New("check mail"))
}

type AuthController interface {
	LoginUserAccount(c *Context) error
	AddUsers(c *Context) error
	UpdateUsers(c *Context) error
	DeleteUsers(c *Context) error
	GetDetailUsers(c *Context) error
	TestJWT(c *Context) error
	ResetPassword(c *Context) error
}

func NewAuthController(us service.UsersService) AuthController {
	return &authController{
		usersService: us,
	}
}
