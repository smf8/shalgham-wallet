package model

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrRecordNotFound = errors.New("record not found")

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
	FindByPhone(ctx context.Context, phoneNumber string) (*Profile, error)
	Create(context.Context, *Profile) error
	UpdateBalance(context.Context, string, float64) error
}

type SQLProfileRepo struct {
	DB *gorm.DB
}

func (p *SQLProfileRepo) Create(ctx context.Context, profile *Profile) error {
	return p.DB.WithContext(ctx).Create(profile).Error
}

func (p *SQLProfileRepo) UpdateBalance(ctx context.Context, phoneNumber string, amount float64) error {
	return p.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

func (p *SQLProfileRepo) FindByPhone(ctx context.Context, phoneNumber string) (*Profile, error) {
	result := &Profile{}

	err := p.DB.WithContext(ctx).Where("phone_number = ?", phoneNumber).First(result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to find profile: %w", ErrRecordNotFound)
		}

		return nil, err
	}

	return result, nil
}
