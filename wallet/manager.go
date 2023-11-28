package wallet

import (
	"context"
	"math/big"
	"sync"

	uuid "github.com/satori/go.uuid"
)

type (
	managerDependencies struct {
		WalletPersister
	}
	ManagerProvider interface {
		WalletManager() *Manager
	}
	Manager struct {
		l sync.Mutex
		d managerDependencies
	}
)

// Transfer one amount from one wallet to another.
//
// If from wallet doesn't have the amount, an error will be thrown
func (p *Manager) Transfer(ctx context.Context, fromWalletID uuid.UUID, toWalletID uuid.UUID, amount *big.Int) error {
	p.l.Lock()
	defer p.l.Unlock()

	fromWallet, err := p.d.WalletPersister.FindWalletByID(ctx, fromWalletID)
	if err != nil {
		return err
	}

	toWallet, err := p.d.WalletPersister.FindWalletByID(ctx, toWalletID)
	if err != nil {
		return err
	}

	rsub := fromWallet.Sub(amount)
	radd := toWallet.Add(amount)

	t := NewTransaction(rsub, radd)

	err = t.Apply()
	if err != nil {
		return err
	}

	err = p.d.WalletPersister.SaveWallet(ctx, []*Wallet{fromWallet, toWallet}, []*WalletRecord{radd, rsub}, []*WalletTransaction{t})
	if err != nil {
		return err
	}

	return nil
}

// Removes amount from the given wallet
//
// If no amount is available in the given wallet, an error will be thrown
func (p *Manager) Burn(ctx context.Context, fromWalletID uuid.UUID, amount *big.Int) error {
	p.l.Lock()
	defer p.l.Unlock()

	fromWallet, err := p.d.WalletPersister.FindWalletByID(ctx, fromWalletID)
	if err != nil {
		return err
	}

	rsub := fromWallet.Sub(amount)

	t := NewTransaction(rsub)

	err = t.Apply()
	if err != nil {
		return err
	}

	err = p.d.WalletPersister.SaveWallet(ctx, []*Wallet{fromWallet}, []*WalletRecord{rsub}, []*WalletTransaction{t})
	if err != nil {
		return err
	}

	return nil
}

// Add funds to the given wallet
func (p *Manager) Mint(ctx context.Context, toWalletID uuid.UUID, amount *big.Int) error {
	p.l.Lock()
	defer p.l.Unlock()

	toWallet, err := p.d.WalletPersister.FindWalletByID(ctx, toWalletID)
	if err != nil {
		return err
	}

	radd := toWallet.Add(amount)

	t := NewTransaction(radd)
	err = t.Apply()
	if err != nil {
		return err
	}

	err = p.d.WalletPersister.SaveWallet(ctx, []*Wallet{toWallet}, []*WalletRecord{radd}, []*WalletTransaction{t})
	if err != nil {
		return err
	}

	return nil
}

// Remove funds from the given wallet
//
// While burn transactions also do the same,
// Lock transaction is created with the meaning of unlock it at any time.
func (p *Manager) Lock(ctx context.Context, fromWalletID uuid.UUID, amount *big.Int) error {
	p.l.Lock()
	defer p.l.Unlock()

	fromWallet, err := p.d.WalletPersister.FindWalletByID(ctx, fromWalletID)
	if err != nil {
		return err
	}

	rlock := fromWallet.Lock(amount)

	t := NewTransaction(rlock)
	err = t.Apply()
	if err != nil {
		return err
	}

	err = p.d.WalletPersister.SaveWallet(ctx, []*Wallet{fromWallet}, []*WalletRecord{rlock}, []*WalletTransaction{t})
	if err != nil {
		return err
	}

	return nil
}

// Unlock a given amount from wallet.
// Unlock can not exceed current locked amount.
func (p *Manager) Unlock(ctx context.Context, fromWalletID uuid.UUID, amount *big.Int) error {
	p.l.Lock()
	defer p.l.Unlock()

	fromWallet, err := p.d.WalletPersister.FindWalletByID(ctx, fromWalletID)
	if err != nil {
		return err
	}

	runlock := fromWallet.Unlock(amount)

	t := NewTransaction(runlock)
	err = t.Apply()
	if err != nil {
		return err
	}

	err = p.d.WalletPersister.SaveWallet(ctx, []*Wallet{fromWallet}, []*WalletRecord{runlock}, []*WalletTransaction{t})
	if err != nil {
		return err
	}

	return nil
}

// Executes a mint and a lock in the same transaction.
//
// Useful when you want to add a fund to a wallet but not want
// the user to spend it until something occur.
func (p *Manager) MintAndLock(ctx context.Context, toWalletID uuid.UUID, amount *big.Int) error {
	p.l.Lock()
	defer p.l.Unlock()

	toWallet, err := p.d.WalletPersister.FindWalletByID(ctx, toWalletID)
	if err != nil {
		return err
	}

	radd := toWallet.Add(amount)
	rlock := toWallet.Lock(amount)

	t := NewTransaction(radd, rlock)
	err = t.Apply()
	if err != nil {
		return err
	}

	err = p.d.WalletPersister.SaveWallet(ctx, []*Wallet{toWallet}, []*WalletRecord{radd, rlock}, []*WalletTransaction{t})
	if err != nil {
		return err
	}

	return nil
}

// Executes an Unlock and burn in the same transaction
//
// Useful when you want unlock some amount and immediatelly burns,
// making this amount unavailable to the given user.
func (p *Manager) UnlockAndBurn(ctx context.Context, fromWalletID uuid.UUID, amount *big.Int) error {
	p.l.Lock()
	defer p.l.Unlock()

	fromWallet, err := p.d.WalletPersister.FindWalletByID(ctx, fromWalletID)
	if err != nil {
		return err
	}

	runlock := fromWallet.Unlock(amount)
	rsub := fromWallet.Sub(amount)

	t := NewTransaction(runlock, rsub)

	err = t.Apply()
	if err != nil {
		return err
	}

	err = p.d.WalletPersister.SaveWallet(ctx, []*Wallet{fromWallet}, []*WalletRecord{runlock, rsub}, []*WalletTransaction{t})
	if err != nil {
		return err
	}

	return nil
}
