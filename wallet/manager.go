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

type ComplexTransferItem struct {
	Currency     WalletCurrency
	Amount       *big.Int
	FromWalletID uuid.UUID
	ToWalletID   uuid.UUID
}

func NewManager(d managerDependencies) *Manager {
	return &Manager{
		l: sync.Mutex{},
		d: d,
	}
}

// Transfer one amount from one wallet to another.
//
// If from wallet doesn't have the amount, an error will be thrown
func (p *Manager) ComplexTransfer(ctx context.Context, toTransfer []*ComplexTransferItem) error {
	p.l.Lock()
	defer p.l.Unlock()

	walletMap := make(map[uuid.UUID]*Wallet)
	wallets := []*Wallet{}

	var records []*WalletRecord
	for _, t := range toTransfer {
		var err error
		fromWallet, found := walletMap[t.FromWalletID]
		if !found {
			fromWallet, err = p.d.WalletPersister().FindWalletByID(ctx, t.FromWalletID)
			if err != nil {
				return err
			}
			walletMap[t.FromWalletID] = fromWallet
			wallets = append(wallets, fromWallet)
		}

		toWallet, found := walletMap[t.ToWalletID]
		if !found {
			toWallet, err = p.d.WalletPersister().FindWalletByID(ctx, t.ToWalletID)
			if err != nil {
				return err
			}
			walletMap[t.ToWalletID] = toWallet
			wallets = append(wallets, toWallet)
		}

		rsub := fromWallet.Sub(t.Amount, t.Currency)
		radd := toWallet.Add(t.Amount, t.Currency)
		records = append(records, rsub, radd)
	}

	t := NewTransaction(records...)

	err := t.Apply()
	if err != nil {
		return err
	}

	err = p.d.WalletPersister().SaveWallet(ctx, wallets, records, []*WalletTransaction{t})
	if err != nil {
		return err
	}

	return nil
}

// Multi currency transfers the same amount of many currencies from one wallet to another
func (p *Manager) MultiTransfer(ctx context.Context, fromWalletID uuid.UUID, toWalletID uuid.UUID, amounts map[WalletCurrency]*big.Int) error {
	p.l.Lock()
	defer p.l.Unlock()

	fromWallet, err := p.d.WalletPersister().FindWalletByID(ctx, fromWalletID)
	if err != nil {
		return err
	}

	toWallet, err := p.d.WalletPersister().FindWalletByID(ctx, toWalletID)
	if err != nil {
		return err
	}

	var records []*WalletRecord
	for currency, amount := range amounts {
		rsub := fromWallet.Sub(amount, currency)
		radd := toWallet.Add(amount, currency)

		records = append(records, rsub, radd)
	}

	t := NewTransaction(records...)

	err = t.Apply()
	if err != nil {
		return err
	}

	err = p.d.WalletPersister().SaveWallet(ctx, []*Wallet{fromWallet, toWallet}, records, []*WalletTransaction{t})
	if err != nil {
		return err
	}

	return nil
}

// Transfer one amount from one wallet to another.
//
// If from wallet doesn't have the amount, an error will be thrown
func (p *Manager) Transfer(ctx context.Context, fromWalletID uuid.UUID, toWalletID uuid.UUID, amount *big.Int, currency WalletCurrency) error {
	p.l.Lock()
	defer p.l.Unlock()

	fromWallet, err := p.d.WalletPersister().FindWalletByID(ctx, fromWalletID)
	if err != nil {
		return err
	}

	toWallet, err := p.d.WalletPersister().FindWalletByID(ctx, toWalletID)
	if err != nil {
		return err
	}

	rsub := fromWallet.Sub(amount, currency)
	radd := toWallet.Add(amount, currency)

	t := NewTransaction(rsub, radd)

	err = t.Apply()
	if err != nil {
		return err
	}

	err = p.d.WalletPersister().SaveWallet(ctx, []*Wallet{fromWallet, toWallet}, []*WalletRecord{radd, rsub}, []*WalletTransaction{t})
	if err != nil {
		return err
	}

	return nil
}

// Removes amount from the given wallet
//
// If no amount is available in the given wallet, an error will be thrown
func (p *Manager) Burn(ctx context.Context, fromWalletID uuid.UUID, amount *big.Int, currency WalletCurrency) error {
	p.l.Lock()
	defer p.l.Unlock()

	fromWallet, err := p.d.WalletPersister().FindWalletByID(ctx, fromWalletID)
	if err != nil {
		return err
	}

	rsub := fromWallet.Sub(amount, currency)

	t := NewTransaction(rsub)

	err = t.Apply()
	if err != nil {
		return err
	}

	err = p.d.WalletPersister().SaveWallet(ctx, []*Wallet{fromWallet}, []*WalletRecord{rsub}, []*WalletTransaction{t})
	if err != nil {
		return err
	}

	return nil
}

// Add funds to the given wallet
func (p *Manager) Mint(ctx context.Context, toWalletID uuid.UUID, amount *big.Int, currency WalletCurrency) error {
	p.l.Lock()
	defer p.l.Unlock()

	toWallet, err := p.d.WalletPersister().FindWalletByID(ctx, toWalletID)
	if err != nil {
		return err
	}

	radd := toWallet.Add(amount, currency)

	t := NewTransaction(radd)
	err = t.Apply()
	if err != nil {
		return err
	}

	err = p.d.WalletPersister().SaveWallet(ctx, []*Wallet{toWallet}, []*WalletRecord{radd}, []*WalletTransaction{t})
	if err != nil {
		return err
	}

	return nil
}

// Remove funds from the given wallet
//
// While burn transactions also do the same,
// Lock transaction is created with the meaning of unlock it at any time.
func (p *Manager) Lock(ctx context.Context, fromWalletID uuid.UUID, amount *big.Int, currency WalletCurrency) error {
	p.l.Lock()
	defer p.l.Unlock()

	fromWallet, err := p.d.WalletPersister().FindWalletByID(ctx, fromWalletID)
	if err != nil {
		return err
	}

	rlock := fromWallet.Lock(amount, currency)

	t := NewTransaction(rlock)
	err = t.Apply()
	if err != nil {
		return err
	}

	err = p.d.WalletPersister().SaveWallet(ctx, []*Wallet{fromWallet}, []*WalletRecord{rlock}, []*WalletTransaction{t})
	if err != nil {
		return err
	}

	return nil
}

// Unlock a given amount from wallet.
// Unlock can not exceed current locked amount.
func (p *Manager) Unlock(ctx context.Context, fromWalletID uuid.UUID, amount *big.Int, currency WalletCurrency) error {
	p.l.Lock()
	defer p.l.Unlock()

	fromWallet, err := p.d.WalletPersister().FindWalletByID(ctx, fromWalletID)
	if err != nil {
		return err
	}

	runlock := fromWallet.Unlock(amount, currency)

	t := NewTransaction(runlock)
	err = t.Apply()
	if err != nil {
		return err
	}

	err = p.d.WalletPersister().SaveWallet(ctx, []*Wallet{fromWallet}, []*WalletRecord{runlock}, []*WalletTransaction{t})
	if err != nil {
		return err
	}

	return nil
}

// Executes a mint and a lock in the same transaction.
//
// Useful when you want to add a fund to a wallet but not want
// the user to spend it until something occur.
func (p *Manager) MintAndLock(ctx context.Context, toWalletID uuid.UUID, amount *big.Int, currency WalletCurrency) error {
	p.l.Lock()
	defer p.l.Unlock()

	toWallet, err := p.d.WalletPersister().FindWalletByID(ctx, toWalletID)
	if err != nil {
		return err
	}

	radd := toWallet.Add(amount, currency)
	rlock := toWallet.Lock(amount, currency)

	t := NewTransaction(radd, rlock)
	err = t.Apply()
	if err != nil {
		return err
	}

	err = p.d.WalletPersister().SaveWallet(ctx, []*Wallet{toWallet}, []*WalletRecord{radd, rlock}, []*WalletTransaction{t})
	if err != nil {
		return err
	}

	return nil
}

// Executes an Unlock and burn in the same transaction
//
// Useful when you want unlock some amount and immediatelly burns,
// making this amount unavailable to the given user.
func (p *Manager) UnlockAndBurn(ctx context.Context, fromWalletID uuid.UUID, amount *big.Int, currency WalletCurrency) error {
	p.l.Lock()
	defer p.l.Unlock()

	fromWallet, err := p.d.WalletPersister().FindWalletByID(ctx, fromWalletID)
	if err != nil {
		return err
	}

	runlock := fromWallet.Unlock(amount, currency)
	rsub := fromWallet.Sub(amount, currency)

	t := NewTransaction(runlock, rsub)

	err = t.Apply()
	if err != nil {
		return err
	}

	err = p.d.WalletPersister().SaveWallet(ctx, []*Wallet{fromWallet}, []*WalletRecord{runlock, rsub}, []*WalletTransaction{t})
	if err != nil {
		return err
	}

	return nil
}

// Trade unlocks a balance from one wallet and send it to another wallet.
// It does the same on the opposit side
func (p *Manager) Trade(ctx context.Context,
	fromWalletID uuid.UUID,
	fromAmount *big.Int,
	fromCurrency WalletCurrency,
	toWalletID uuid.UUID,
	toAmount *big.Int,
	toCurrency WalletCurrency) error {
	p.l.Lock()
	defer p.l.Unlock()

	fromWallet, err := p.d.WalletPersister().FindWalletByID(ctx, fromWalletID)
	if err != nil {
		return err
	}

	toWallet, err := p.d.WalletPersister().FindWalletByID(ctx, toWalletID)
	if err != nil {
		return err
	}

	rfunlock := fromWallet.Unlock(fromAmount, fromCurrency)
	rsub := fromWallet.Sub(fromAmount, fromCurrency)
	tadd := toWallet.Add(fromAmount, fromCurrency)

	tfunlock := toWallet.Unlock(toAmount, toCurrency)
	tsub := toWallet.Sub(toAmount, toCurrency)
	radd := fromWallet.Add(toAmount, toCurrency)

	t := NewTransaction(rfunlock, rsub, tadd, tfunlock, tsub, radd)
	err = t.Apply()
	if err != nil {
		return err
	}

	err = p.d.WalletPersister().SaveWallet(ctx, []*Wallet{fromWallet}, t.WalletRecords, []*WalletTransaction{t})
	if err != nil {
		return err
	}

	return nil
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
