package postgres

import (
	"embed"

	"github.com/Masterminds/squirrel"
)

//go:embed migrations/*
var Migrations embed.FS

func SetupPSQL() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

const (
	TableAccounts     = "accounts"
	TableTranscations = "transcations"

	ID              = "id"
	AccountID       = "account_id"
	OperationTypeID = "operation_type_id"
	DocumentNumber  = "document_number"
	CreatedAt       = "created_at"
	UpdatedAt       = "updated_at"
	EventAt         = "event_at"
	Amount          = "amount"
	CreditLimit     = "credit_limit"
	WithdrewalLimit = "withdrawal_limit"
	Balance         = "balance"
)
