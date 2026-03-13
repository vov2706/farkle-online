package handlers

import (
	"app/http/inputs"
	"app/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	service     *services.AuthService
	userService *services.UserService
}

func NewAuthHandler(service *services.AuthService, userService *services.UserService) *AuthHandler {
	return &AuthHandler{service: service, userService: userService}
}

func (handler *AuthHandler) Login(c fiber.Ctx) error {
	userData, err := validateUserData(c)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user, err := handler.userService.GetUserByUsername(userData.Username)

	const dummyHash = "2a$10$7zFqzDbD3RrlkMTczbXG9OWZ0FLOXjIxXzSZ.QZxkVXjXcx7QZQiC"

	if err != nil {
		// prevents timing attacks
		handler.service.IsValidPassword(dummyHash, userData.Password)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid username or password",
		})
	}

	if !handler.service.IsValidPassword(user.Password, userData.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid username or password",
		})
	}

	token, err := handler.service.CreateToken(user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Logged in",
		"token":   token,
	})
}

func (handler *AuthHandler) Register(c fiber.Ctx) error {
	userData, err := validateUserData(c)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	_, err = handler.userService.GetUserByUsername(userData.Username)

	// if the user already exists
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User already exists",
		})
	}

	hashedPassword := handler.service.GeneratePassword(userData.Password)
	user, err := handler.userService.CreateUser(userData.Username, hashedPassword)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user, err = handler.userService.GetUserByUsername(userData.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to load created user",
		})
	}

	token, err := handler.service.CreateToken(user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Success registration",
		"token":   token,
	})
}

func validateUserData(c fiber.Ctx) (*inputs.LoginInput, error) {
	input := new(inputs.LoginInput)

	if err := c.Bind().Body(input); err != nil {
		return nil, err
	}

	// Validation logic
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		errs := err.(validator.ValidationErrors)

		for _, e := range errs {
			return nil, e
		}
	}

	return input, nil
}
