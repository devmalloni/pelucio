package pelucio

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/devmalloni/pelucio/x/xtime"

	"github.com/gofrs/uuid/v5"
)

var (
	ErrNoAccountProvided        = errors.New("no accounts provided")
	ErrEntriesNotFound          = errors.New("no entries in transaction")
	ErrTransactionIsNotBalanced = errors.New("transaction is not balanced")
	ErrAccountNotFound          = errors.New("account not found")
	ErrEntryTransactionMismatch = errors.New("entry transaction ID does not match transaction ID")
)

type (
	Transaction struct {
		ID          uuid.UUID       `json:"id" db:"id"`
		ExternalID  string          `json:"external_id" db:"external_id"`
		Description string          `json:"description" db:"description"`
		Metadata    json.RawMessage `json:"metadata" db:"metadata"`

		CreatedAt time.Time `json:"created_at" db:"created_at"`

		Entries []*Entry `json:"entries" db:"-"`
	}
)

func (p Transaction) BalancesByOperation() (balanceOwn Balance, balanceOwe Balance) {
	balanceOwn = make(Balance)
	balanceOwe = make(Balance)

	balancesByAccount := p.BalancesByAccount(nil)

	sidesByAccount := p.SideByAccounts()
	for account, balance := range balancesByAccount {
		switch sidesByAccount[account] {
		case Debit:
			balanceOwn.AddBalance(balance)
		case Credit:
			balanceOwe.AddBalance(balance)
		}
	}

	return
}

func (p Transaction) BalancesByAccount(balancesOfAccounts map[uuid.UUID]Balance) map[uuid.UUID]Balance {
	if balancesOfAccounts == nil {
		balancesOfAccounts = make(map[uuid.UUID]Balance)
	}

	for _, entry := range p.Entries {
		if balancesOfAccounts[entry.AccountID] == nil {
			balancesOfAccounts[entry.AccountID] = make(Balance)
		}

		entry.UnsafeApply(balancesOfAccounts[entry.AccountID])
	}

	return balancesOfAccounts
}

func (p *Transaction) IsBalanced() bool {
	balanceOwn, balanceOwe := p.BalancesByOperation()

	return balanceOwn.IsBalanced(balanceOwe)
}

func (p *Transaction) ApplyToAccounts(accounts map[uuid.UUID]*Account) error {
	if len(accounts) == 0 {
		return ErrNoAccountProvided
	}

	if len(p.Entries) == 0 {
		return ErrEntriesNotFound
	}

	if !p.IsBalanced() {
		return ErrTransactionIsNotBalanced
	}

	for _, entry := range p.Entries {
		account, ok := accounts[entry.AccountID]
		if !ok {
			return ErrAccountNotFound
		}

		if entry.TransactionID != p.ID {
			return ErrEntryTransactionMismatch
		}

		if err := entry.Apply(account.Balance); err != nil {
			return err
		}
	}

	return nil
}

func (p *Transaction) Accounts() []uuid.UUID {
	accounts := make([]uuid.UUID, 0)

	for _, entry := range p.Entries {
		accounts = append(accounts, entry.AccountID)
	}

	return accounts
}

func (p *Transaction) SideByAccounts() map[uuid.UUID]EntrySide {
	accounts := make(map[uuid.UUID]EntrySide)

	for _, entry := range p.Entries {
		if accountSide, ok := accounts[entry.AccountID]; ok && accountSide != entry.AccountSide {
			panic("there are multiple entries with different account side")
		}

		accounts[entry.AccountID] = entry.AccountSide
	}

	return accounts
}

func (p *Transaction) Reverse(externalID, description string, clock xtime.Clock) *Transaction {
	reversed := &Transaction{
		ID:          p.ID,
		ExternalID:  externalID,
		Description: description,
		CreatedAt:   clock.Now(),
	}

	for _, entry := range p.Entries {
		reversedEntry := entry.Reverse(reversed.ID, clock)
		reversed.Entries = append(reversed.Entries, &reversedEntry)
	}

	return reversed
}
