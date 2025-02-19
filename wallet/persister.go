package wallet

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

type (
	WalletPersister interface {
		SaveWallet(ctx context.Context, w []*Wallet, r []*WalletRecord, t []*WalletTransaction) error
		FindWalletByID(ctx context.Context, id uuid.UUID) (*Wallet, error)
		FindWalletByExternalID(ctx context.Context, id string) (*Wallet, error)
		FindWalletRecords(ctx context.Context, walletID uuid.UUID) ([]*WalletRecord, error)
		FindWalletTransactions(ctx context.Context, walletID uuid.UUID) ([]*WalletTransaction, error)
	}

	PersisterProvider interface {
		WalletPersister() WalletPersister
	}
)
