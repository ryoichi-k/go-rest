package usecase

import (
	"go-rest/model"
	"go-rest/repository"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/golang-jwt/jwt/v4"
)

type IUserUsecase interface {
	SignUp(user model.User) (model.UserResponse, error)
	Login(user model.User) (string, error)
}

type userUsecase struct {
	ur repository.IUserRepository
}

// constructor
// 外部でインスタンス化されるurを受け取るようにしている
func NewUserUsecase(ur repository.IUserRepository) IUserUsecase {
	return &userUsecase{ur}
}

func (uu *userUsecase) SignUp(user model.User) (model.UserResponse, error) {
	//平文のパスワードをハッシュ化
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return model.UserResponse{}, err
	}
	newUser := model.User{Email: user.Email, Password: string(hash)}
	if err := uu.ur.CreateUser(&newUser); err != nil {
		return model.UserResponse{}, err
	}
	// エラー発生していない場合は新規ユーザー情報がレスポンスモデルに代入される
	resUser := model.UserResponse{
		ID:    newUser.ID,
		Email: newUser.Email,
	}
	return resUser, nil
}

func (uu *userUsecase) Login(user model.User) (string, error) {
	// userリクエストをもとにユーザー情報がDBにあるか
	storedUser := model.User{}
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		return "", err
	}
	// パスワードの検証(hash化されたパスと平文のパスワードを比較)
	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}
	//パスワードが一致、jwtトークンを生成する
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		"exp":     time.Now().Add(time.Hour * 12).Unix(), //jwtトークンの有効期限設定
	})

	//jwtトークン生成
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil

}
