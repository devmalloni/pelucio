package pelucio

import (
	"encoding/json"

	"github.com/gofrs/uuid/v5"
)

type AccountOption func(*Account)

func WithID(id uuid.UUID) func(*Account) {
	return func(a *Account) {
		a.ID = id
	}
}

func WithName(name string) func(*Account) {
	return func(a *Account) {
		a.Name = name
	}
}

func WithNormalSide(normalSide EntrySide) func(*Account) {
	return func(a *Account) {
		a.NormalSide = normalSide
	}
}

func WithMetadata(metadata json.RawMessage) func(*Account) {
	return func(a *Account) {
		a.Metadata = metadata
	}
}

func WithExternalID(externalID string) func(*Account) {
	return func(a *Account) {
		a.ExternalID = externalID
	}
}
