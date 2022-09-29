package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
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
