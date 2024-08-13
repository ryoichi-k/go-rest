package controller

import (
	"go-rest/model"
	"go-rest/usecase"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type IUserController interface {
	SignUp(c echo.Context) error
	LogIn(c echo.Context) error
	LogOut(c echo.Context) error
	//CsrfToken(c echo.Context) error
}

type userController struct {
	uu usecase.IUserUsecase
}

// constructor
// <I>usecaseからcontrollerにDI注入のため
func NewUserController(uu usecase.IUserUsecase) IUserController {
	return &userController{uu}
}

func (uc *userController) SignUp(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	userRes, err := uc.uu.SignUp(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// code 201
	return c.JSON(http.StatusCreated, userRes)
}

func (uc *userController) LogIn(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	//jwtトークン作成
	tokenString, err := uc.uu.Login(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	//認証成功時、jwtトークンをサーバーサイドでクッキーに設定する
	cookie := new(http.Cookie) //cookie作成
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	cookie.HttpOnly = true                  //clientのjsからトークンが読めないようにしておく
	cookie.SameSite = http.SameSiteNoneMode // クロスドメインによる送受信なので
	c.SetCookie(cookie)                     // cookieをhttpレスポンスに含める
	return c.NoContent(http.StatusOK)       // okステータスをクライアントに返す
}

func (uc *userController) LogOut(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	cookie.HttpOnly = true                  //clientのjsからトークンが読めないようにしておく
	cookie.SameSite = http.SameSiteNoneMode // クロスドメインによる送受信なので
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)
}
