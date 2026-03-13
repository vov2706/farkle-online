package services

import (
	"app/config"
	"app/models"
	"app/repositories"
	"time"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (service AuthService) GetAuthUser(c fiber.Ctx, relations ...string) (*models.User, error) {
	var token = jwtware.FromContext(c)
	if token == nil {
		return nil, fiber.ErrUnauthorized
	}

	claims := token.Claims.(jwt.MapClaims)
	userId, ok := claims["sub"].(float64)

	if !ok {
		return nil, fiber.ErrUnauthorized
	}

	user, err := service.userRepo.GetUserById(uint(userId), relations...)

	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	return user, nil
}

func (service AuthService) CreateToken(u *models.User) (t string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, err = token.SignedString([]byte(config.Config("JWT_SECRET")))

	if err != nil {
		return "", err
	}

	return t, nil
}

func (service AuthService) GeneratePassword(p string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(p), 8)

	return string(hash)
}

func (service AuthService) IsValidPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}
