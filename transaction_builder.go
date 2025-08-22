package pelucio

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/devmalloni/pelucio/x/xtime"
	"github.com/devmalloni/pelucio/x/xuuid"

	"github.com/gofrs/uuid/v5"
)

type TransactionBuilder struct {
	externalID  string
	description string
	metadata    json.RawMessage
	entries     []*Entry

	defaultClock xtime.Clock
}

func NewTransaction(clock xtime.Clock) *TransactionBuilder {
	if clock == nil {
		clock = xtime.DefaultClock
	}

	return &TransactionBuilder{
		defaultClock: clock,
	}
}

func (p *TransactionBuilder) WithExternalID(externalID string) *TransactionBuilder {
	p.externalID = externalID
	return p
}

func (p *TransactionBuilder) WithDescription(description string) *TransactionBuilder {
	p.description = description
	return p
}

func (p *TransactionBuilder) WithMetadata(metadata json.RawMessage) *TransactionBuilder {
	p.metadata = metadata
	return p
}

func (p *TransactionBuilder) AddEntry(accountID uuid.UUID,
	entrySide EntrySide,
	accountSide EntrySide,
	amount *big.Int,
	currency Currency) *TransactionBuilder {
	p.entries = append(p.entries, &Entry{
		ID:          xuuid.New(),
		AccountID:   accountID,
		EntrySide:   entrySide,
		AccountSide: accountSide,
		Amount:      amount,
		Currency:    currency,
		CreatedAt:   p.defaultClock.Now(),
	})
	return p
}

func (p *TransactionBuilder) Build() (*Transaction, error) {
	transactionID := xuuid.New()

	for _, entry := range p.entries {
		entry.TransactionID = transactionID
	}

	if len(p.entries) == 0 {
		return nil, errors.New("transaction must have at least one entry")
	}

	if p.externalID == "" {
		p.externalID = xuuid.New().String()
	}

	t := &Transaction{
		ID:          transactionID,
		ExternalID:  p.externalID,
		Description: p.description,
		Metadata:    p.metadata,
		Entries:     p.entries,
		CreatedAt:   p.defaultClock.Now(),
	}

	if !t.IsBalanced() {
		return nil, errors.New("transaction is not balanced")
	}

	return t, nil
}

func (p *TransactionBuilder) MustBuild() *Transaction {
	t, err := p.Build()
	if err != nil {
		panic(err)
	}

	return t
}
