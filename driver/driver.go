package driver

import (
	"pelucio/config"
	"pelucio/persistence/bbolt"
	"pelucio/wallet"
)

type (
	Driver struct {
		persister     wallet.WalletPersister
		walletManager *wallet.Manager
		walletHandler *wallet.Handler
		config        *config.Config
	}
)

func (p *Driver) Config() *config.Config {
	if p.config == nil {
		p.config = &config.Config{}
	}
	return p.config
}

func (p *Driver) WalletPersister() wallet.WalletPersister {
	if p.persister == nil {
		persister, err := bbolt.NewPersister(p)
		if err != nil {
			panic(err)
		}
		p.persister = persister
	}

	return p.persister
}

func (p *Driver) WalletManager() *wallet.Manager {
	if p.walletManager == nil {
		p.walletManager = wallet.NewManager(p)
	}
	return p.walletManager
}

func (p *Driver) WalletHandler() *wallet.Handler {
	if p.walletHandler == nil {
		p.walletHandler = wallet.NewHandler(p)
	}
	return p.walletHandler
}
