package pelucio

import (
	"context"
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
)

var (
	ErrNotFound = errors.New("record not found")
)

type (
	ReadAccountFilter struct {
		FromDate        *time.Time `json:"from_date,omitempty"`
		ToDate          *time.Time `json:"to_date,omitempty"`
		AccountIDs      []string   `json:"account_i_ds,omitempty"`
		ExternalIDs     []string   `json:"external_i_ds,omitempty"`
		PaginationToken *string    `json:"pagination_token,omitempty"`
		Limit           *uint      `json:"limit,omitempty"`
	}

	ReadTransactionFilter struct {
		FromDate        *time.Time `json:"from_date,omitempty"`
		ToDate          *time.Time `json:"to_date,omitempty"`
		AccountIDs      []string   `json:"account_i_ds,omitempty"`
		ExternalIDs     []string   `json:"external_i_ds,omitempty"`
		PaginationToken *string    `json:"pagination_token,omitempty"`
		Limit           *uint      `json:"limit,omitempty"`
	}

	ReadEntryFilter struct {
		FromDate        *time.Time `json:"from_date,omitempty"`
		ToDate          *time.Time `json:"to_date,omitempty"`
		AccountIDs      []string   `json:"account_i_ds,omitempty"`
		TransactionIDs  []string   `json:"transaction_i_ds,omitempty"`
		PaginationToken *string    `json:"pagination_token,omitempty"`
		Limit           *uint      `json:"limit,omitempty"`
	}

	Writer interface {
		WriteAccount(ctx context.Context, account *Account, allowUpdate bool) error
		WriteTransaction(ctx context.Context, transaction *Transaction, account ...*Account) error
	}

	Reader interface {
		ReadAccount(ctx context.Context, accountID uuid.UUID) (*Account, error)
		ReadAccountByExternalID(ctx context.Context, externalID string) (*Account, error)
		ReadAccounts(ctx context.Context, filter ReadAccountFilter) (accounts []*Account, paginationToken *string, err error)

		ReadTransaction(ctx context.Context, transactionID uuid.UUID) (*Transaction, error)
		ReadTransactionByExternalID(ctx context.Context, externalID string) (*Transaction, error)
		ReadTransactions(ctx context.Context, filter ReadTransactionFilter) (transactions []*Transaction, paginationToken *string, err error)

		ReadEntriesOfAccount(ctx context.Context, accountID uuid.UUID) ([]*Entry, error)
		ReadEntries(ctx context.Context, filter ReadEntryFilter) (entries []*Entry, paginationToken *string, err error)
	}

	ReadWriter interface {
		Reader
		Writer
	}
)
