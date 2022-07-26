package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/smf8/arvan-wallet/internal/app/model"
	"net/http"
)

type Transaction struct {
	ProfileRepo model.ProfileRepo
}

type Profile struct {
	ProfileRepo model.ProfileRepo
}

func (t *Transaction) ApplyTransaction(c *fiber.Ctx) error {
	request := &TransactionRequest{}

	if err := c.BodyParser(request); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	if err := t.ProfileRepo.UpdateBalance(
		request.PhoneNumber, request.Amount); err != nil {
		logrus.Errorf("update balance failed: %s", err.Error())

		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.SendStatus(http.StatusOK)
}

func (p *Profile) Create(c *fiber.Ctx) error {
	request := &ProfileCreateRequest{}

	if err := c.BodyParser(request); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	profile := &model.Profile{
		PhoneNumber: request.PhoneNumber,
		Balance:     request.Balance,
	}

	if err := p.ProfileRepo.Create(profile); err != nil {
		logrus.Errorf("profile create failed: %s", err.Error())

		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.SendStatus(http.StatusOK)
}

func (p *Profile) Get(c *fiber.Ctx) error {
	phoneNumber := c.Params("phone")

	if phoneNumber == "" {
		return c.SendStatus(http.StatusBadRequest)
	}

	profile, err := p.ProfileRepo.FindByPhone(phoneNumber)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			return c.SendStatus(http.StatusNotFound)
		}

		logrus.Errorf("profile create failed: %s", err.Error())

		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(profile)
}
