package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/madhurikadam/app-transcation/internal/domain"
	"github.com/madhurikadam/app-transcation/internal/service/mocks"

	"github.com/stretchr/testify/suite"
)

var (
	errTestFoo = fmt.Errorf("error foo")
)

type ServiceTestSuite struct {
	suite.Suite

	repo *mocks.MockRepo

	svc TranscationService
}

func TestService(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) SetupTest() {
	s.repo = mocks.NewMockRepo(gomock.NewController(s.T()))

	s.svc = New(s.repo)
}

func (s *ServiceTestSuite) TestCreateAccount() {
	ctx := context.Background()
	documentNumber := "12345678"

	tests := []struct {
		name           string
		mocks          func()
		documentNumber string
		expErr         bool
		expError       error
	}{
		{
			name:           "invalid document number",
			mocks:          func() {},
			documentNumber: "",
			expErr:         true,
			expError:       ErrInvalidDocumentNumber,
		},
		{
			name: "failed to create account in database",
			mocks: func() {
				s.repo.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(errTestFoo)
			},
			documentNumber: documentNumber,
			expErr:         true,
			expError:       errTestFoo,
		},
		{
			name: "create account with success",
			mocks: func() {
				s.repo.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(nil)
			},
			documentNumber: documentNumber,
		},
	}

	for _, tt := range tests {
		tt := tt

		s.Run(tt.name, func() {
			s.SetupTest()
			tt.mocks()

			account, err := s.svc.CreateAccount(ctx, tt.documentNumber)
			if tt.expErr {
				s.Require().Error(err)
				s.Require().Equal(tt.expError, err)

				return
			}

			s.Require().NoError(err)
			s.Require().Equal(tt.documentNumber, account.DocumentNumber)
			s.NotNil(account.ID)
		})
	}
}

func (s *ServiceTestSuite) TestGetAccount() {
	ctx := context.Background()
	accountID := "12345678"
	testAcc := &domain.Account{
		ID:              accountID,
		DocumentNumber:  accountID,
		WithdrawalLimit: 0,
	}

	tests := []struct {
		name            string
		mocks           func()
		accountID       string
		expErr          bool
		expError        error
		expectedAccount *domain.Account
	}{
		{
			name:      "invalid account id",
			mocks:     func() {},
			accountID: "",
			expErr:    true,
			expError:  ErrInvalidAccountID,
		},
		{
			name: "failed to get account from database",
			mocks: func() {
				s.repo.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(nil, errTestFoo)
			},
			accountID: accountID,
			expErr:    true,
			expError:  errTestFoo,
		},
		{
			name: "get account with success",
			mocks: func() {
				s.repo.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(testAcc, nil)
			},
			accountID:       accountID,
			expectedAccount: testAcc,
		},
	}

	for _, tt := range tests {
		tt := tt

		s.Run(tt.name, func() {
			s.SetupTest()
			tt.mocks()

			account, err := s.svc.GetAccount(ctx, tt.accountID)
			if tt.expErr {
				s.Require().Error(err)
				s.Require().Equal(tt.expError, err)

				return
			}

			s.Require().NoError(err)
			s.Require().Equal(tt.expectedAccount, account)
			s.NotNil(account.ID)
		})
	}
}

func (s *ServiceTestSuite) TestCreateTranscation() {
	ctx := context.Background()
	accountID := "12345678"

	tests := []struct {
		name       string
		mocks      func()
		input      domain.Transcation
		expErr     bool
		expError   error
		expectedOp domain.Transcation
	}{
		{
			name:  "invalid account id",
			mocks: func() {},
			input: domain.Transcation{
				AccountID: "",
			},
			expErr:   true,
			expError: ErrInvalidAccountID,
		},
		{
			name: "invalid operation type",
			mocks: func() {
			},
			input: domain.Transcation{
				AccountID:       accountID,
				OperationTypeID: 0,
			},
			expErr:   true,
			expError: ErrInvalidOperationTypeID,
		},
		{
			name: "failed to debit create transcation in database",
			mocks: func() {
				s.repo.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(&domain.Account{WithdrawalLimit: 400}, nil)
				s.repo.EXPECT().CreateDebitTranscation(gomock.Any(), gomock.Any()).Return(errTestFoo)
			},
			input: domain.Transcation{
				AccountID:       accountID,
				OperationTypeID: 2,
				Amount:          20,
			},
			expErr:   true,
			expError: errTestFoo,
		},
		{
			name: "debit amount is greater than debit limit",
			mocks: func() {
				s.repo.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(&domain.Account{WithdrawalLimit: 400}, nil)
			},
			input: domain.Transcation{
				AccountID:       accountID,
				OperationTypeID: 2,
				Amount:          500,
			},
			expErr:   true,
			expError: fmt.Errorf("exceed withdrwal limit"),
		},
		{
			name: "credit amount is greater than creidt limit",
			mocks: func() {
				s.repo.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(&domain.Account{CreaditLimit: 400}, nil)
			},
			input: domain.Transcation{
				AccountID:       accountID,
				OperationTypeID: 4,
				Amount:          500,
			},
			expErr:   true,
			expError: fmt.Errorf("exceed credit limit"),
		},
		{
			name: "create debit transcation with success",
			mocks: func() {
				s.repo.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(&domain.Account{WithdrawalLimit: 400}, nil)
				s.repo.EXPECT().CreateDebitTranscation(gomock.Any(), gomock.Any()).Return(nil)
			},
			input: domain.Transcation{
				AccountID:       accountID,
				OperationTypeID: 2,
				Amount:          20,
			},
			expectedOp: domain.Transcation{
				AccountID:       accountID,
				OperationTypeID: 2,
				Amount:          -20,
			},
		},
		{
			name: "create credit transcation with success",
			mocks: func() {
				s.repo.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(&domain.Account{CreaditLimit: 400}, nil)
				s.repo.EXPECT().CreateCreditTranscation(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				s.repo.EXPECT().ListDebitTx(gomock.Any()).Return([]domain.Transcation{
					{
						ID:      uuid.NewString(),
						Amount:  50,
						Balance: -50,
					},
					{
						ID:      uuid.NewString(),
						Amount:  23.5,
						Balance: -23.5,
					},
					{
						ID:      uuid.NewString(),
						Amount:  18.7,
						Balance: -18.7,
					},
				}, nil)

			},
			input: domain.Transcation{
				AccountID:       accountID,
				OperationTypeID: 4,
				Amount:          20,
			},
			expectedOp: domain.Transcation{
				AccountID:       accountID,
				OperationTypeID: 4,
				Amount:          20,
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		s.Run(tt.name, func() {
			s.SetupTest()
			tt.mocks()

			tx, err := s.svc.CreateTranscation(ctx, tt.input)
			if tt.expErr {
				s.Require().Error(err)
				s.Require().Equal(tt.expError, err)

				return
			}

			s.Require().NoError(err)
			s.NotNil(tx.ID)
			s.Equal(tt.expectedOp.OperationTypeID, tx.OperationTypeID)
			s.Equal(tt.expectedOp.AccountID, tx.AccountID)
			s.Equal(tt.expectedOp.Amount, tx.Amount)
		})
	}
}
func (s *ServiceTestSuite) TestDispatchTx() {
	ctx := context.Background()
	accountID := "12345678"
	id1 := uuid.NewString()
	id2 := uuid.NewString()
	id3 := uuid.NewString()
	tests := []struct {
		name              string
		mocks             func()
		input             domain.Transcation
		expErr            bool
		expError          error
		expectedCreditBal float64
		expectedDList     []domain.DebitTx
	}{
		{
			name: "create credit transcation with success where debit balance is greater than credit balance",
			mocks: func() {
				s.repo.EXPECT().ListDebitTx(gomock.Any()).Return([]domain.Transcation{
					{
						ID:      id1,
						Amount:  50,
						Balance: -50,
					},
					{
						ID:      id2,
						Amount:  23.5,
						Balance: -23.5,
					},
					{
						ID:      id3,
						Amount:  18.7,
						Balance: -18.7,
					},
				}, nil)
			},
			input: domain.Transcation{
				AccountID:       accountID,
				OperationTypeID: 4,
				Amount:          60,
			},
			expectedCreditBal: 0,
			expectedDList: []domain.DebitTx{
				{
					ID:     id1,
					Amount: 0,
				},
				{
					ID:     id2,
					Amount: -13.5,
				},
				{
					ID:     id3,
					Amount: -18.7,
				},
			},
		},
		{
			name: "create credit transcation with success where debit balance is greater than credit balance",
			mocks: func() {
				s.repo.EXPECT().ListDebitTx(gomock.Any()).Return([]domain.Transcation{
					{
						ID:      id1,
						Amount:  50,
						Balance: 0,
					},
					{
						ID:      id2,
						Amount:  23.5,
						Balance: -13.5,
					},
					{
						ID:      id3,
						Amount:  18.7,
						Balance: -18.7,
					},
				}, nil)
			},
			input: domain.Transcation{
				AccountID:       accountID,
				OperationTypeID: 4,
				Amount:          100,
			},
			expectedCreditBal: 67.8,
			expectedDList: []domain.DebitTx{
				{
					ID:     id1,
					Amount: 0,
				},
				{
					ID:     id2,
					Amount: 0,
				},
				{
					ID:     id3,
					Amount: 0,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		s.Run(tt.name, func() {
			s.SetupTest()
			tt.mocks()

			creditBal, dList, err := s.svc.dispatchTx(ctx, tt.input)
			if tt.expErr {
				s.Require().Error(err)
				s.Require().Equal(tt.expError, err)

				return
			}

			s.Require().NoError(err)
			s.Require().Equal(tt.expectedCreditBal, creditBal)
			s.Require().ElementsMatch(tt.expectedDList, dList)
			for _, bal := range dList {
				for _, aBal := range tt.expectedDList {
					if aBal.ID == bal.ID {
						s.Require().Equal(aBal.Amount, bal.Amount)
					}
				}
			}
		})
	}
}
