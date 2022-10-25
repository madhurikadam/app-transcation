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

		CreateCreditTranscation(ctx context.Context, transcation domain.Transcation, dbTxList []domain.DebitTx) error
		CreateDebitTranscation(ctx context.Context, transcation domain.Transcation) error
		// GetCreditBalance(ctx context.Context) (float64, error)
		ListDebitTx(ctx context.Context) ([]domain.Transcation, error)
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
		transcation.Balance = transcation.Amount
		transcation.Amount = -transcation.Amount
		if err := t.repo.CreateDebitTranscation(ctx, transcation); err != nil {
			return nil, err
		}

		return &transcation, nil
	}

	creditBalance, dTxList, err := t.dispatchTx(ctx, transcation)
	if err != nil {
		return nil, err
	}
	transcation.Balance = creditBalance
	if err := t.repo.CreateCreditTranscation(ctx, transcation, dTxList); err != nil {
		return nil, err
	}

	return &transcation, nil
}

func (t *TranscationService) dispatchTx(ctx context.Context, transcation domain.Transcation) (float64, []domain.DebitTx, error) {
	var balance float64
	dTxList := make([]domain.DebitTx, 0)
	dList, err := t.repo.ListDebitTx(ctx)
	if err != nil {
		return balance, dTxList, err
	}

	if len(dList) <= 0 {
		balance = transcation.Amount
		return balance, dTxList, nil
	}
	txAmount := transcation.Amount
	for _, val := range dList {
		if -val.Balance <= txAmount {
			dTxList = append(dTxList, domain.DebitTx{
				ID:     val.ID,
				Amount: 0,
			})

			txAmount = txAmount + val.Balance
		} else {
			//set balance of debit tx to tx.amount - txAmount
			dTxList = append(dTxList, domain.DebitTx{
				ID:     val.ID,
				Amount: val.Balance + txAmount,
			})
			txAmount = 0
		}
	}

	return txAmount, dTxList, nil
}

func validateOpTypeID(opTypeID int) error {
	switch opTypeID {
	case 1, 2, 3, 4:
		return nil
	default:
		return ErrInvalidOperationTypeID
	}
}
