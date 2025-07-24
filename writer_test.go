package pelucio

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/mock"
)

type ReadWriterMock struct {
	mock.Mock
}

func (p *ReadWriterMock) WriteAccount(ctx context.Context, account *Account, allowUpdate bool) error {
	args := p.Called(account, allowUpdate)
	return args.Error(0)
}

func (p *ReadWriterMock) WriteTransaction(ctx context.Context, transaction *Transaction, account ...*Account) error {
	args := p.Called(transaction, account)
	return args.Error(0)
}

func (p *ReadWriterMock) ReadAccount(ctx context.Context, accountID uuid.UUID) (*Account, error) {
	args := p.Called(accountID)
	res, _ := args.Get(0).(*Account)
	return res, args.Error(1)
}

func (p *ReadWriterMock) ReadAccountByExternalID(ctx context.Context, externalID string) (*Account, error) {
	args := p.Called(externalID)
	res, _ := args.Get(0).(*Account)
	return res, args.Error(1)
}

func (p *ReadWriterMock) ReadAccounts(ctx context.Context, filter ReadAccountFilter) ([]*Account, error) {
	args := p.Called(filter)
	res, _ := args.Get(0).([]*Account)
	return res, args.Error(1)
}

func (p *ReadWriterMock) ReadTransaction(ctx context.Context, transactionID uuid.UUID) (*Transaction, error) {
	args := p.Called(transactionID)
	res, _ := args.Get(0).(*Transaction)
	return res, args.Error(1)
}

func (p *ReadWriterMock) ReadTransactions(ctx context.Context, filter ReadTransactionFilter) ([]*Transaction, error) {
	args := p.Called(filter)
	res, _ := args.Get(0).([]*Transaction)
	return res, args.Error(1)
}

func (p *ReadWriterMock) ReadEntriesOfAccount(ctx context.Context, accountID uuid.UUID) ([]*Entry, error) {
	args := p.Called(accountID)
	res, _ := args.Get(0).([]*Entry)
	return res, args.Error(1)
}

func (p *ReadWriterMock) ReadEntries(ctx context.Context, filter ReadEntryFilter) ([]*Entry, error) {
	args := p.Called(filter)
	res, _ := args.Get(0).([]*Entry)
	return res, args.Error(1)
}
