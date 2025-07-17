package pelucio

import (
	"encoding/json"
	"errors"
	"math/big"
	"pelucio/x/xtime"
	"pelucio/x/xuuid"
	"time"

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
	if accounts == nil {
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
	accounts := make([]uuid.UUID, len(p.Entries))

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

func Transfer(externalID string, fromAccountID, toAccountID uuid.UUID, amount *big.Int, currency Currency) *Transaction {
	transactionID := xuuid.New() // This will be set later, e.g., by a database or UUID generator
	return &Transaction{
		ID:         transactionID,
		ExternalID: externalID,
		Entries: []*Entry{
			{
				ID:            xuuid.New(),
				TransactionID: transactionID,
				AccountID:     fromAccountID,
				EntrySide:     Debit,
				AccountSide:   Credit,
				Amount:        amount,
				Currency:      currency,
				CreatedAt:     time.Now(),
			},
			{
				ID:            xuuid.New(),
				TransactionID: transactionID,
				AccountID:     toAccountID,
				EntrySide:     Credit,
				AccountSide:   Credit,
				Amount:        amount,
				Currency:      currency,
				CreatedAt:     time.Now(),
			},
		},
	}
}

func Deposit(externalID string, cashAccountID, toAccountID uuid.UUID, amount *big.Int, currency Currency) *Transaction {
	transactionID := xuuid.New() // This will be set later, e.g., by a database or UUID generator
	return &Transaction{
		ID:         transactionID,
		ExternalID: externalID,
		Entries: []*Entry{
			{
				ID:            xuuid.New(),
				TransactionID: transactionID,
				AccountID:     cashAccountID,
				EntrySide:     Debit,
				AccountSide:   Debit,
				Amount:        amount,
				Currency:      currency,
				CreatedAt:     time.Now(),
			},
			{
				ID:            xuuid.New(),
				TransactionID: transactionID,
				AccountID:     toAccountID,
				EntrySide:     Credit,
				AccountSide:   Credit,
				Amount:        amount,
				Currency:      currency,
				CreatedAt:     time.Now(),
			},
		},
	}
}

func DepositWithFee(externalID string, cashAccountID, feeAccount, toAccountID uuid.UUID, feeTax, amount *big.Int, currency Currency) *Transaction {
	if feeTax == nil || feeTax.Cmp(big.NewInt(0)) <= 0 {
		// error
	}

	fee := new(big.Int).Mul(feeTax, amount)
	amountWithoutFee := new(big.Int).Sub(amount, fee)
	transactionID := xuuid.New()

	return &Transaction{
		ID:         transactionID,
		ExternalID: externalID,
		Entries: []*Entry{
			{
				ID:            xuuid.New(),
				TransactionID: transactionID,
				AccountID:     cashAccountID,
				EntrySide:     Debit,
				AccountSide:   Debit,
				Amount:        amount,
				Currency:      currency,
				CreatedAt:     time.Now(),
			},
			{
				ID:            xuuid.New(),
				TransactionID: transactionID,
				AccountID:     feeAccount,
				EntrySide:     Credit,
				AccountSide:   Credit,
				Amount:        fee,
				Currency:      currency,
				CreatedAt:     time.Now(),
			},
			{
				ID:            xuuid.New(),
				TransactionID: transactionID,
				AccountID:     toAccountID,
				EntrySide:     Credit,
				AccountSide:   Credit,
				Amount:        amountWithoutFee,
				Currency:      currency,
				CreatedAt:     time.Now(),
			},
		},
	}
}

func Withdraw(externalID string, cashAccountID, fromAccount uuid.UUID, amount *big.Int, currency Currency) *Transaction {
	transactionID := xuuid.New()
	return &Transaction{
		ID:         transactionID,
		ExternalID: externalID,
		Entries: []*Entry{
			{
				ID:            xuuid.New(),
				TransactionID: transactionID,
				AccountID:     cashAccountID,
				EntrySide:     Credit,
				AccountSide:   Debit,
				Amount:        amount,
				Currency:      currency,
				CreatedAt:     time.Now(),
			},
			{
				ID:            xuuid.New(),
				TransactionID: transactionID,
				AccountID:     fromAccount,
				EntrySide:     Debit,
				AccountSide:   Credit,
				Amount:        amount,
				Currency:      currency,
				CreatedAt:     time.Now(),
			},
		},
	}
}

func WithdrawWithFee(externalID string, cashAccountID, feeAccount, fromAccount uuid.UUID, feeTax, amount *big.Int, currency Currency) *Transaction {
	if feeTax == nil || feeTax.Cmp(big.NewInt(0)) <= 0 {
		// error
	}

	fee := new(big.Int).Mul(feeTax, amount)
	amountWithFee := new(big.Int).Add(amount, fee)
	transactionID := xuuid.New()

	return &Transaction{
		ID:         transactionID,
		ExternalID: externalID,
		Entries: []*Entry{
			{
				ID:            xuuid.New(),
				TransactionID: transactionID,
				AccountID:     cashAccountID,
				EntrySide:     Credit,
				AccountSide:   Debit,
				Amount:        amount,
				Currency:      currency,
				CreatedAt:     time.Now(),
			},
			{
				ID:            xuuid.New(),
				TransactionID: transactionID,
				AccountID:     feeAccount,
				EntrySide:     Credit,
				AccountSide:   Credit,
				Amount:        fee,
				Currency:      currency,
				CreatedAt:     time.Now(),
			},
			{
				ID:            xuuid.New(),
				TransactionID: transactionID,
				AccountID:     fromAccount,
				EntrySide:     Debit,
				AccountSide:   Credit,
				Amount:        amountWithFee,
				Currency:      currency,
				CreatedAt:     time.Now(),
			},
		},
	}
}

// func Trade(externalID,
// 	takerAccountID string,
// 	takerCurrency Currency,
// 	takerAmount *big.Int,

// 	makerAccountID string,
// 	makerCurrency Currency,
// 	makerAmount *big.Int,

// 	feeAccountID string,
// 	makerFeeCurrency Currency,
// 	makerFeePercent *big.Int,
// 	takerFeeCurrency Currency,
// 	takerFeePercent *big.Int) {

// 	transaction := &Transaction{
// 		ID:         "",
// 		ExternalID: externalID,
// 		Entries: []*Entry{
// 			{ // debit on taker
// 				ID:            "",
// 				TransactionID: "",
// 				AccountID:     takerAccountID,
// 				EntrySide:     Debit,
// 				AccountSide:   Credit,
// 				Amount:        takerAmount,
// 				Currency:      takerCurrency,
// 				CreatedAt:     time.Now(),
// 			},
// 			{ // credit on taker
// 				ID:            "",
// 				TransactionID: transactionID,
// 				AccountID:     feeAccount,
// 				EntrySide:     Credit,
// 				AccountSide:   Credit,
// 				Amount:        fee,
// 				Currency:      currency,
// 				CreatedAt:     time.Now(),
// 			},
// 			{
// 				ID:            "",
// 				TransactionID: transactionID,
// 				AccountID:     fromAccount,
// 				EntrySide:     Debit,
// 				AccountSide:   Credit,
// 				Amount:        amountWithFee,
// 				Currency:      currency,
// 				CreatedAt:     time.Now(),
// 			},
// 		},
// 	}
// 	// debit on taker
// 	// credit on taker

// 	// debit on maker
// 	// credit on maker

// 	// credit on fee account for maker
// 	// credit on fee account for taker
// }
