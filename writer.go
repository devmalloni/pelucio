package pelucio

import (
	"context"
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
)

var (
	ErrNotFound         = errors.New("record not found")
	ErrDuplicatedRecord = errors.New("duplicated record")
)

type (
	ReadAccountFilter struct {
		FromDate    *time.Time `json:"from_date,omitempty"`
		ToDate      *time.Time `json:"to_date,omitempty"`
		AccountIDs  []string   `json:"account_i_ds,omitempty"`
		ExternalIDs []string   `json:"external_i_ds,omitempty"`
	}

	ReadTransactionFilter struct {
		FromDate    *time.Time `json:"from_date,omitempty"`
		ToDate      *time.Time `json:"to_date,omitempty"`
		AccountIDs  []string   `json:"account_i_ds,omitempty"`
		ExternalIDs []string   `json:"external_i_ds,omitempty"`
	}

	ReadEntryFilter struct {
		FromDate       *time.Time `json:"from_date,omitempty"`
		ToDate         *time.Time `json:"to_date,omitempty"`
		AccountIDs     []string   `json:"account_i_ds,omitempty"`
		TransactionIDs []string   `json:"transaction_i_ds,omitempty"`
	}

	Writer interface {
		WriteAccount(ctx context.Context, account *Account, allowUpdate bool) error
		WriteTransaction(ctx context.Context, transaction *Transaction, account ...*Account) error
	}

	Reader interface {
		ReadAccount(ctx context.Context, accountID uuid.UUID) (*Account, error)
		ReadAccountByExternalID(ctx context.Context, externalID string) (*Account, error)
		ReadAccounts(ctx context.Context, filter ReadAccountFilter) ([]*Account, error)

		ReadTransaction(ctx context.Context, transactionID uuid.UUID) (*Transaction, error)
		ReadTransactions(ctx context.Context, filter ReadTransactionFilter) ([]*Transaction, error)

		ReadEntriesOfAccount(ctx context.Context, accountID uuid.UUID) ([]*Entry, error)
		ReadEntries(ctx context.Context, filter ReadEntryFilter) ([]*Entry, error)
	}

	ReadWriter interface {
		Reader
		Writer
	}
)
