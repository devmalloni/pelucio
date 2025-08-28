package pelucio

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/devmalloni/pelucio/x/xmap"
	"github.com/devmalloni/pelucio/x/xtime"
	"github.com/devmalloni/pelucio/x/xuuid"

	"github.com/gofrs/uuid/v5"
)

var (
	ErrExternalIDAlreadyInUse = errors.New("external id already in use")
	ErrRequiredAccountID      = errors.New("account ID cannot be empty")
	ErrTransactionNil         = errors.New("transaction cannot be nil")
)

type Pelucio struct {
	readWriter ReadWriter
	clock      xtime.Clock
}

func NewPelucio(opts ...PelucionOpt) *Pelucio {
	p := &Pelucio{
		clock: xtime.StdClock{},
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *Pelucio) CreateAccount(ctx context.Context,
	externalID,
	name string,
	normalSide EntrySide,
	metadata json.RawMessage) (*Account, error) {
	account := NewAccount(externalID, name, normalSide, metadata, p.clock)

	_, err := p.readWriter.ReadAccountByExternalID(ctx, externalID)
	if err != nil && err != ErrNotFound {
		return nil, err
	}

	if err == nil {
		return nil, ErrExternalIDAlreadyInUse
	}

	err = p.readWriter.WriteAccount(ctx, account, false)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (p *Pelucio) DeleteAccount(ctx context.Context, accountID uuid.UUID) error {
	account, err := p.readWriter.ReadAccount(ctx, accountID)
	if err != nil {
		return err
	}

	err = account.Delete(p.clock)
	if err != nil {
		return err
	}

	err = p.readWriter.WriteAccount(ctx, account, true)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pelucio) UpdateAccount(ctx context.Context, accountID uuid.UUID, name string, metadata json.RawMessage) error {
	account, err := p.readWriter.ReadAccount(ctx, accountID)
	if err != nil {
		return err
	}

	account.UpdateData(name, metadata, p.clock)

	err = p.readWriter.WriteAccount(ctx, account, true)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pelucio) FindAccounts(ctx context.Context, query ReadAccountFilter) ([]*Account, error) {
	return p.readWriter.ReadAccounts(ctx, query)
}

func (p *Pelucio) FindAccountByID(ctx context.Context, accountID uuid.UUID) (*Account, error) {
	return p.readWriter.ReadAccount(ctx, accountID)
}

func (p *Pelucio) FindAccountByExternalID(ctx context.Context, externalID string) (*Account, error) {
	return p.readWriter.ReadAccountByExternalID(ctx, externalID)
}

func (p *Pelucio) FindTransactionByID(ctx context.Context, id uuid.UUID) (*Transaction, error) {
	return p.readWriter.ReadTransaction(ctx, id)
}

func (p *Pelucio) FindTransactioByExternalID(ctx context.Context, externalID string) (*Transaction, error) {
	return p.readWriter.ReadTransactionByExternalID(ctx, externalID)
}

func (p *Pelucio) BalanceOf(ctx context.Context, accountID uuid.UUID) (Balance, error) {
	account, err := p.readWriter.ReadAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return account.Balance, nil
}

func (p *Pelucio) BalanceOfAccountFromLedger(ctx context.Context, accountID uuid.UUID) (Balance, error) {
	if xuuid.IsNilOrEmpty(accountID) {
		return nil, ErrRequiredAccountID
	}

	account, err := p.readWriter.ReadAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	entries, err := p.readWriter.ReadEntriesOfAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	err = account.ComputeFromEntries(entries, p.clock)
	if err != nil {
		return nil, err
	}

	return account.Balance, nil
}

func (p *Pelucio) ExecuteTransaction(ctx context.Context, transaction *Transaction) error {
	if transaction == nil {
		return ErrTransactionNil
	}

	// check for existing external id
	_, err := p.readWriter.ReadTransactionByExternalID(ctx, transaction.ExternalID)
	if err != nil && err != ErrNotFound {
		return err
	}
	if err == nil {
		return ErrExternalIDAlreadyInUse
	}

	accounts, err := p.FindAccounts(ctx, ReadAccountFilter{
		AccountIDs: xuuid.ToStrings(transaction.Accounts()...),
	})
	if err != nil {
		return err
	}

	accountMap := xmap.ToMap(accounts, func(a *Account) uuid.UUID {
		return a.ID
	})

	err = transaction.ApplyToAccounts(accountMap, p.clock)
	if err != nil {
		return err
	}

	accounts = xmap.Values(accountMap)

	return p.readWriter.WriteTransaction(ctx, transaction, accounts...)
}

func (p *Pelucio) RevertTransaction(ctx context.Context, originalTransactionID uuid.UUID, externalID string) error {
	originalTransaction, err := p.readWriter.ReadTransaction(ctx, originalTransactionID)
	if err != nil {
		return err
	}

	if externalID == "" {
		externalID = fmt.Sprintf("%s-%s", originalTransaction.ExternalID, "revert")
	}
	revertTransaction := originalTransaction.
		Reverse(externalID,
			fmt.Sprintf("reverted transaction %s", originalTransactionID), p.clock)

	err = p.ExecuteTransaction(ctx, revertTransaction)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pelucio) EntriesOfAccount(ctx context.Context, accountID uuid.UUID) ([]*Entry, error) {
	if xuuid.IsNilOrEmpty(accountID) {
		return nil, ErrRequiredAccountID
	}

	return p.readWriter.ReadEntriesOfAccount(ctx, accountID)
}

func (p *Pelucio) FindEntries(ctx context.Context, filter ReadEntryFilter) ([]*Entry, error) {
	return p.readWriter.ReadEntries(ctx, filter)
}
