package pelucio

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/devmalloni/pelucio/x/xtime"
	"github.com/devmalloni/pelucio/x/xuuid"

	"github.com/gofrs/uuid/v5"
)

var (
	ErrEntryAccountMismatch = errors.New("entry account ID does not match account ID")
	ErrAccountSideMismatch  = errors.New("account side does not match normal side")
	ErrBalanceNotEmpty      = errors.New("cannot delete account with non-zero balance")
	ErrInvalidNormalSide    = errors.New("invalid normal side")
	ErrInvalidID            = errors.New("invalid id")
	ErrInvalidExternalID    = errors.New("ErrInvalidExternalID")
)

type (
	Account struct {
		ID         uuid.UUID `json:"id" db:"id"`
		ExternalID string    `json:"external_id" db:"external_id"`
		Name       string    `json:"name" db:"name"`

		NormalSide EntrySide       `json:"normal_side" db:"normal_side"`
		Balance    Balance         `json:"balance" db:"balance"`
		Metadata   json.RawMessage `json:"metadata" db:"metadata"`

		Version   int64      `json:"version" db:"version"` // version must be handled by persistence layer
		CreatedAt time.Time  `json:"created_at" db:"created_at"`
		UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	}
)

func NewAccount(clock xtime.Clock, opts ...AccountOption) *Account {
	a := &Account{
		ID:         xuuid.New(),
		ExternalID: xuuid.New().String(),
		Balance:    make(Balance),
		CreatedAt:  clock.Now(),
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (p *Account) UpdateData(name string, metadata json.RawMessage, clock xtime.Clock) {
	p.Name = name
	p.Metadata = metadata
	p.UpdatedAt = clock.NilNow()
}

func (p *Account) Apply(e Entry, clock xtime.Clock) error {
	if p.Balance == nil {
		p.Balance = make(Balance)
	}

	if !xuuid.Equal(p.ID, e.AccountID) {
		return ErrEntryAccountMismatch
	}

	if e.AccountSide != p.NormalSide {
		return ErrAccountSideMismatch
	}

	err := e.Apply(p.Balance)
	if err != nil {
		return err
	}

	p.UpdatedAt = clock.NilNow()

	return nil
}

func (p *Account) ComputeFromEntries(entries []*Entry, clock xtime.Clock) error {
	if p.Balance != nil {
		p.Balance.Clear()
	}

	for _, e := range entries {
		err := p.Apply(*e, clock)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Account) Delete(clock xtime.Clock) error {
	if p.Balance != nil && p.Balance.HasBalance() {
		return ErrBalanceNotEmpty
	}

	p.DeletedAt = clock.NilNow()
	p.UpdatedAt = clock.NilNow()

	return nil
}

func (p *Account) IsValid() error {
	if xuuid.IsNilOrEmpty(p.ID) {
		return ErrInvalidID
	}

	if p.NormalSide != Debit && p.NormalSide != Credit {
		return ErrInvalidNormalSide
	}

	if p.ExternalID == "" {
		return ErrInvalidExternalID
	}

	return nil
}
