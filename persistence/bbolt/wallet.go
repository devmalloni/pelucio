package bbolt

import (
	"context"
	"encoding/json"
	"pelucio/config"
	"pelucio/wallet"

	uuid "github.com/satori/go.uuid"
	bolt "go.etcd.io/bbolt"
)

type (
	Persister struct {
		db *bolt.DB
	}
	persisterDependencies interface {
		config.Provider
	}
)

func NewPersister(d persisterDependencies) (*Persister, error) {
	// add zerolog
	db, err := bolt.Open(d.Config().PATH(), 0600, nil)
	if err != nil {
		return nil, err
	}

	return &Persister{
		db: db,
	}, nil
}

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

func (p *Persister) FindWalletByExternalID(ctx context.Context, id string) (*wallet.Wallet, error) {
	var res *wallet.Wallet
	err := p.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Wallets"))
		err := b.ForEach(func(k, v []byte) error {
			var wa *wallet.Wallet
			err := json.Unmarshal(v, &wa)
			if err != nil {
				return err
			}

			if wa.ExternalID == id {
				res = wa
			}

			return nil
		})

		return err
	})

	return res, err
}

func (p *Persister) FindWalletRecords(ctx context.Context, walletID uuid.UUID) ([]*wallet.WalletRecord, error) {
	var res []*wallet.WalletRecord
	err := p.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("WalletRecords"))
		v := b.Get(walletID.Bytes())

		if v != nil {
			err := json.Unmarshal(v, &res)
			if err != nil {
				return err
			}
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

func (p *Persister) FindWallets(ctx context.Context) ([]*wallet.Wallet, error) {
	var res []*wallet.Wallet
	err := p.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Wallets"))
		err := b.ForEach(func(k, v []byte) error {
			var wa *wallet.Wallet
			err := json.Unmarshal(v, &wa)
			if err != nil {
				return err
			}

			res = append(res, wa)
			return nil
		})

		return err
	})

	return res, err
}
