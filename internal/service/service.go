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

		CreateCreditTranscation(ctx context.Context, transcation domain.Transcation) error
		CreateDebitTranscation(ctx context.Context, transcation domain.Transcation) error
	}
)

var (
	ErrInvalidDocumentNumber  = fmt.Errorf("invalid document id")
	ErrInvalidAccountID       = fmt.Errorf("invalid account id")
	ErrInvalidOperationTypeID = fmt.Errorf("invalid operation type id")

	defaultCreditLimit    = 1000.00
	defaultWithdrwalLimit = 1000.00
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
		ID:              uuid.NewString(),
		DocumentNumber:  documentNumber,
		CreatedAt:       now,
		UpdatedAt:       &now,
		CreaditLimit:    defaultCreditLimit,
		WithdrawalLimit: defaultWithdrwalLimit,
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
		return nil, ErrInvalidAccountID
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
		return nil, ErrInvalidAccountID
	}

	if err := validateOpTypeID(transcation.OperationTypeID); err != nil {
		return nil, err
	}

	transcation.ID = uuid.NewString()
	transcation.EventAt = time.Now().UTC()

	acc, err := t.repo.GetAccount(ctx, transcation.AccountID)
	if err != nil {
		return nil, err
	}

	if (transcation.OperationTypeID == 1 || transcation.OperationTypeID == 2 || transcation.OperationTypeID == 3) && acc.WithdrawalLimit < transcation.Amount {
		return nil, fmt.Errorf("exceed withdrwal limit")
	} else if transcation.OperationTypeID == 4 && transcation.Amount > acc.CreaditLimit {
		return nil, fmt.Errorf("exceed credit limit")
	}

	if transcation.OperationTypeID == 1 || transcation.OperationTypeID == 2 || transcation.OperationTypeID == 3 {
		transcation.Amount = -transcation.Amount
		if err := t.repo.CreateDebitTranscation(ctx, transcation); err != nil {
			return nil, err
		}

		return &transcation, nil
	}

	if err := t.repo.CreateCreditTranscation(ctx, transcation); err != nil {
		return nil, err
	}

	return &transcation, nil
}

func validateOpTypeID(opTypeID int) error {
	switch opTypeID {
	case 1, 2, 3, 4:
		return nil
	default:
		return ErrInvalidOperationTypeID
	}
}
