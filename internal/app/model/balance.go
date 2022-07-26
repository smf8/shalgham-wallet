package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Profile struct {
	ID          int64   `json:"id"`
	PhoneNumber string  `json:"phone_number"`
	Balance     float64 `json:"balance"`
}

type Transaction struct {
	ID        int     `json:"id"`
	ProfileID int64   `json:"profile_id"`
	Amount    float64 `json:"amount"`
}

type ProfileRepo interface {
	Create(*Profile) error
	UpdateBalance(string, float64) error
}

type SQLProfileRepo struct {
	DB *gorm.DB
}

func (p *SQLProfileRepo) Create(profile *Profile) error {
	return p.DB.Create(profile).Error
}

func (p *SQLProfileRepo) UpdateBalance(phoneNumber string, amount float64) error {
	return p.DB.Transaction(func(tx *gorm.DB) error {
		profile := &Profile{PhoneNumber: phoneNumber}

		err := tx.Model(profile).Clauses(
			clause.Returning{Columns: []clause.Column{{Name: "id"}}},
		).
			Where("phone_number = ?", profile.PhoneNumber).
			UpdateColumn("balance", profile.Balance+amount).Error

		if err != nil {
			return err
		}

		transaction := &Transaction{Amount: amount, ProfileID: profile.ID}

		return tx.Create(transaction).Error
	})
}
