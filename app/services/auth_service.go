package services

import (
	"errors"
	"go-test/models"
	"go-test/repositories"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Signup(email string, password string) error
	Login(email string, password string) (*string, error)
	GetUserFormToken(tokenString string) (*models.User, error)
}

type AuthService struct {
	repository repositories.IAuthRepository
}

func NewAuthService(repository repositories.IAuthRepository) IAuthService {
	return &AuthService{repository: repository}
}

func (s *AuthService) Signup(email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	return s.repository.CreateUser(user)
}

func (s *AuthService) Login(email string, password string) (*string, error) {
	foundUser, err := s.repository.FindUser(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	token, err := CreateToken(foundUser.ID, foundUser.Email)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func CreateToken(userId uint, email string) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"email":  email,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func (s *AuthService) GetUserFormToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// トークンが期限切れの場合の処理
				return nil, errors.New("Token has expired")
			}
		}
		return nil, err
	}
	var user *models.User
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user, err = s.repository.FindUser(claims["email"].(string))
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}
