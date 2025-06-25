package wallet

import (
	"context"
	"math/big"
	"sync"

	uuid "github.com/satori/go.uuid"
)

type (
	managerDependencies interface {
		PersisterProvider
	}
	ManagerProvider interface {
		WalletManager() *Manager
	}
	Manager struct {
		l sync.Mutex
		d managerDependencies
	}
)

func NewManager(d managerDependencies) *Manager {
	return &Manager{
		l: sync.Mutex{},
		d: d,
	}
}

func (p *Manager) WalletByID(ctx context.Context, walletID uuid.UUID) (*Wallet, error) {
	return p.d.WalletPersister().FindWalletByID(ctx, walletID)
}

func (p *Manager) WalletByExternalID(ctx context.Context, externalID string) (*Wallet, error) {
	return p.d.WalletPersister().FindWalletByExternalID(ctx, externalID)
}

func (p *Manager) WalletRecordsByID(ctx context.Context, walletID uuid.UUID) ([]*WalletRecord, error) {
	return p.d.WalletPersister().FindWalletRecords(ctx, walletID)
}

func (p *Manager) CreateWallet(ctx context.Context, id *uuid.UUID) (*Wallet, error) {
	wallet := &Wallet{
		ID:            uuid.NewV4(),
		Balance:       make(map[WalletCurrency]*big.Int),
		LockedBalance: make(map[WalletCurrency]*big.Int),
		Version:       uuid.NewV4(),
	}

	if id != nil {
		wallet.ID = *id
	}

	err := p.d.WalletPersister().SaveWallet(ctx, []*Wallet{wallet}, nil, nil)

	return wallet, err
}

func (p *Manager) GetWallets(ctx context.Context) ([]*Wallet, error) {
	return p.d.WalletPersister().FindWallets(ctx)
}
