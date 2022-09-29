package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/madhurikadam/app-transcation/internal/domain"
)

type (
	TranscationService struct {
		repo Repo
	}

	Repo interface {
		CreateAccount(ctx context.Context, account domain.Account) error
		GetAccount(ctx context.Context, id string) (*domain.Account, error)

		CreateTranscation(ctx context.Context, transcation domain.Transcation) error
	}
)

var (
	ErrInvalidDocumentNumber = fmt.Errorf("invalid document id")
)

func New(repo Repo) TranscationService {
	return TranscationService{
		repo: repo,
	}
}

// CreateAccount create account with document number
func (t *TranscationService) CreateAccount(ctx context.Context, documentNumber string) (*domain.Account, error) {
	if documentNumber == "" {
		return nil, ErrInvalidDocumentNumber
	}

	now := time.Now().UTC()

	account := domain.Account{
		ID:             uuid.NewString(),
		DocumentNumber: documentNumber,
		CreatedAt:      now,
		UpdatedAt:      &now,
	}

	err := t.repo.CreateAccount(ctx, account)
	if err != nil {
		log.Error("failed to create account", err)
		return nil, err
	}

	return &account, nil
}

// GetAccount get account details via account id
func (t *TranscationService) GetAccount(ctx context.Context, accountID string) (*domain.Account, error) {
	if accountID == "" {
		return nil, fmt.Errorf("invalid account id")
	}

	account, err := t.repo.GetAccount(ctx, accountID)
	if err != nil {
		log.WithField("account_id", accountID).Error("failed to get account", err)
		return nil, err
	}

	return account, nil
}

// CreateTranscation add new transcation for given account id
func (t *TranscationService) CreateTranscation(ctx context.Context, transcation domain.Transcation) (*domain.Transcation, error) {
	if transcation.AccountID == "" {
		return nil, fmt.Errorf("invalid account id")
	}

	if err := validateOpTypeID(transcation.OperationTypeID); err != nil {
		return nil, err
	}

	if err := validateOpAndAmount(transcation.OperationTypeID, transcation.Amount); err != nil {
		return nil, err
	}

	transcation.ID = uuid.NewString()
	transcation.EventAt = time.Now().UTC()

	// TODO validation on the balance available and transcation amount
	// before doing any transcation on account.

	err := t.repo.CreateTranscation(ctx, transcation)
	if err != nil {
		log.WithField("account_id", transcation.AccountID).Error("failed to create transcation", err)
		return nil, err
	}

	return &transcation, nil
}

func validateOpTypeID(opTypeID int) error {
	switch opTypeID {
	case 1, 2, 3, 4:
		return nil
	default:
		return fmt.Errorf("invalid operation type id")
	}
}

func validateOpAndAmount(opTypeID int, amount float64) error {
	if (opTypeID == 1 || opTypeID == 2 || opTypeID == 3) && amount >= 0 {
		return fmt.Errorf("invalid amount for operation type id %d", opTypeID)
	}

	if opTypeID == 4 && amount <= 0 {
		return fmt.Errorf("invalid amount for operation type id %d", opTypeID)
	}

	return nil
}
