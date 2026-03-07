package handlers

import (
	"app/models"
	"app/repositories"

	"github.com/gofiber/fiber/v3"
)

type CurrencyHandler struct {
	repo *repositories.CurrencyRepository
}

func NewCurrencyHandler(repo *repositories.CurrencyRepository) *CurrencyHandler {
	return &CurrencyHandler{repo: repo}
}

func (handler *CurrencyHandler) Index(c fiber.Ctx) error {
	currencies, err := handler.repo.GetCurrencies()

	if err != nil {
		return c.JSON(fiber.Map{
			"data": make([]models.Currency, 0),
		})
	}

	return c.JSON(fiber.Map{
		"data": currencies,
	})
}
