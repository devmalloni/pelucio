package bbolt

import (
	"context"
	"encoding/json"
	"pelucio/wallet"

	uuid "github.com/satori/go.uuid"
	bolt "go.etcd.io/bbolt"
)

type (
	Persister struct {
		db *bolt.DB
	}
)

func (p *Persister) SaveWallet(ctx context.Context, w []*wallet.Wallet, r []*wallet.WalletRecord, t []*wallet.WalletTransaction) error {
	return p.db.Update(func(tx *bolt.Tx) error {
		wb := tx.Bucket([]byte("Wallets"))
		rb := tx.Bucket([]byte("WalletRecords"))
		tb := tx.Bucket([]byte("WalletTransactions"))
		for _, wallet := range w {
			walletjson, err := json.Marshal(wallet)
			if err != nil {
				return err
			}
			err = wb.Put(wallet.ID.Bytes(), walletjson)
			if err != nil {
				return err
			}

		}

		for _, data := range r {
			datajson, err := json.Marshal(data)
			if err != nil {
				return err
			}
			err = rb.Put(data.ID.Bytes(), datajson)
			if err != nil {
				return err
			}

		}

		for _, data := range t {
			datajson, err := json.Marshal(data)
			if err != nil {
				return err
			}
			err = tb.Put(data.ID.Bytes(), datajson)
			if err != nil {
				return err
			}

		}

		return nil
	})
}

func (p *Persister) FindWalletByID(ctx context.Context, id uuid.UUID) (*wallet.Wallet, error) {
	var res *wallet.Wallet
	err := p.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Wallets"))
		v := b.Get(id.Bytes())

		err := json.Unmarshal(v, &res)
		if err != nil {
			return err
		}

		return nil
	})

	return res, err
}

func (p *Persister) FindWalletRecords(ctx context.Context, walletID uuid.UUID) ([]*wallet.WalletRecord, error) {
	var res []*wallet.WalletRecord
	err := p.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("WalletRecords"))
		v := b.Get(walletID.Bytes())

		err := json.Unmarshal(v, &res)
		if err != nil {
			return err
		}

		return nil
	})

	return res, err
}

func (p *Persister) FindWalletTransactions(ctx context.Context, walletID uuid.UUID) ([]*wallet.WalletTransaction, error) {
	var res []*wallet.WalletTransaction
	err := p.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("WalletTransactions"))
		v := b.Get(walletID.Bytes())

		err := json.Unmarshal(v, &res)
		if err != nil {
			return err
		}

		return nil
	})

	return res, err
}
