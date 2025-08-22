package pelucio

import (
	"errors"
	"math/big"
	"time"

	"github.com/devmalloni/pelucio/x/xtime"
	"github.com/devmalloni/pelucio/x/xuuid"

	"github.com/gofrs/uuid/v5"
)

var (
	ErrNilBalance          = errors.New("balance cannot be nil")
	ErrNotPositiveAmount   = errors.New("amount must be higher than 0")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

const (
	Debit  EntrySide = "debit"
	Credit EntrySide = "credit"
)

const (
	OperationKindAdd OperationKind = "add"
	OperationKindSub OperationKind = "sub"
)

type (
	OperationKind string
	EntrySide     string
	Currency      string

	Balance map[Currency]*big.Int

	Entry struct {
		ID            uuid.UUID `json:"id" db:"id"`
		TransactionID uuid.UUID `json:"transaction_id" db:"transaction_id"`
		AccountID     uuid.UUID `json:"account_id" db:"account_id"`

		EntrySide   EntrySide `json:"entry_side" db:"entry_side"`
		AccountSide EntrySide `json:"account_side" db:"account_side"`
		Amount      *big.Int  `json:"amount" db:"amount"`
		Currency    Currency  `json:"currency" db:"currency"`

		CreatedAt time.Time `json:"created_at" db:"created_at"`
	}
)

func (p Entry) Apply(balance Balance) error {
	if balance == nil {
		return ErrNilBalance
	}

	switch p.OperationOnBalance() {
	case OperationKindAdd:
		if err := balance.Add(p.Currency, p.Amount); err != nil {
			return err
		}
	case OperationKindSub:
		if err := balance.Sub(p.Currency, p.Amount); err != nil {
			return err
		}
	}

	return nil
}

func (p Entry) UnsafeApply(balance Balance) {
	switch p.OperationOnBalance() {
	case OperationKindAdd:
		balance.UnsafeAdd(p.Currency, p.Amount)
	case OperationKindSub:
		balance.UnsafeSub(p.Currency, p.Amount)
	}
}

func (p Entry) Reverse(newTransactionID uuid.UUID, clock xtime.Clock) Entry {
	reversedEntrySide := Credit

	if p.EntrySide == Credit {
		reversedEntrySide = Debit
	}

	return Entry{
		ID:            xuuid.New(),
		EntrySide:     reversedEntrySide,
		AccountID:     p.AccountID,
		AccountSide:   p.AccountSide,
		Amount:        new(big.Int).Set(p.Amount),
		Currency:      p.Currency,
		CreatedAt:     clock.Now(),
		TransactionID: newTransactionID,
	}
}

func (p Entry) OperationOnBalance() OperationKind {
	if (p.AccountSide == Debit && p.EntrySide == Debit) ||
		(p.AccountSide == Credit && p.EntrySide == Credit) {
		return OperationKindAdd
	} else if (p.AccountSide == Debit && p.EntrySide == Credit) ||
		(p.AccountSide == Credit && p.EntrySide == Debit) {
		return OperationKindSub
	} else {
		panic("invalid entry side or account side")
	}
}

func (p Balance) Add(currency Currency, amount *big.Int) error {
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return ErrNotPositiveAmount
	}

	p.UnsafeAdd(currency, amount)

	return nil
}

func (p Balance) Sub(currency Currency, amount *big.Int) error {
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return ErrNotPositiveAmount
	}

	if p[currency] == nil || p[currency].Cmp(amount) < 0 {
		return ErrInsufficientBalance
	}

	p.UnsafeSub(currency, amount)

	return nil
}

func (p Balance) UnsafeAdd(currency Currency, amount *big.Int) {
	if p[currency] == nil {
		p[currency] = new(big.Int)
	}

	p[currency].Add(p[currency], amount)
}

func (p Balance) UnsafeSub(currency Currency, amount *big.Int) {
	if p[currency] == nil {
		p[currency] = new(big.Int)
	}

	p[currency].Sub(p[currency], amount)
}

func (p Balance) AddBalance(o Balance) {
	for currency, amount := range o {
		if _, ok := p[currency]; !ok {
			p[currency] = new(big.Int)
		}
		p[currency].Add(p[currency], amount)
	}
}

func (p Balance) IsBalanced(o Balance) bool {
	// normalize o balance to ensure it has all currencies from p
	// if the only differences is currencies with no balances
	balanceZeroCurrencies := func(x, z Balance) {
		for k, v := range x {
			if v.Cmp(big.NewInt(0)) == 0 {
				z[k] = big.NewInt(0)
			}
		}
	}

	balanceZeroCurrencies(p, o)
	balanceZeroCurrencies(o, p)

	// check if both balances have the same currencies
	if len(p) != len(o) {
		return false
	}

	// check if there is any currency with different amounts
	for currency, amount := range p {
		if amount.Cmp(o[currency]) != 0 {
			return false
		}
	}

	return true
}

func (p Balance) HasBalance() bool {
	if len(p) == 0 {
		return false
	}
	// check if there is any currency with different amounts
	zero := big.NewInt(0)
	for _, amount := range p {
		if amount.Cmp(zero) != 0 {
			return true
		}
	}

	return false
}

func (p Balance) Clear() {
	for k := range p {
		delete(p, k)
	}
}

func (p Balance) Decimal(precision int) map[Currency]string {
	res := make(map[Currency]string, len(p))

	for currency, amount := range p {
		res[currency] = ToString(amount, precision)
	}

	return res
}

func (p Balance) DecimalFromMap(precisions map[Currency]int, defaultPrecision int) map[Currency]string {
	res := make(map[Currency]string, len(p))

	for currency, amount := range p {
		if precision, ok := precisions[currency]; ok {
			res[currency] = ToString(amount, precision)
		} else {
			res[currency] = ToString(amount, defaultPrecision)
		}
	}

	return res
}
