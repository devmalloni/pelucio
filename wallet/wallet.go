package wallet

import (
	"errors"
	"math/big"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	Sum WalletRecordKind = iota + 1
	Sub
	Lock
	Unlock
)
const (
	Pending WalletTransactionStatus = iota + 1
	Complete
	Failed
)

type (
	WalletRecordKind        uint8
	WalletTransactionStatus uint8

	Wallet struct {
		ID        uuid.UUID
		AccountID uuid.UUID

		Balance       *big.Int
		LockedBalance *big.Int

		CreatedAt time.Time
		UpdatedAt *time.Time
		Version   uuid.UUID
	}

	WalletRecord struct {
		ID       uuid.UUID
		WalletID uuid.UUID
		Value    *big.Int
		Kind     WalletRecordKind
		Nonce    uint64

		Wallet    *Wallet
		CreatedAt time.Time
	}

	WalletTransaction struct {
		ID            uuid.UUID
		WalletRecords []*WalletRecord
		Status        WalletTransactionStatus
		CreatedAt     time.Time
	}
)

func (p *Wallet) SetBalance(balance, lockedBalance *big.Int) {
	p.Balance = balance
	p.Version = uuid.NewV4()
}

func (p *Wallet) Record(amount *big.Int, kind WalletRecordKind) *WalletRecord {
	return &WalletRecord{
		ID:        uuid.NewV4(),
		WalletID:  p.ID,
		Value:     amount,
		Kind:      kind,
		CreatedAt: time.Now(),
		Wallet:    p,
	}
}

func (p *Wallet) Add(amount *big.Int) *WalletRecord {
	return p.Record(amount, Sum)
}

func (p *Wallet) Sub(amount *big.Int) *WalletRecord {
	return p.Record(amount, Sub)
}

func (p *Wallet) Lock(amount *big.Int) *WalletRecord {
	return p.Record(amount, Lock)
}

func (p *Wallet) Unlock(amount *big.Int) *WalletRecord {
	return p.Record(amount, Unlock)
}

func (p *Wallet) Apply(records ...*WalletRecord) error {
	currentAmount := new(big.Int).Set(p.Balance)
	currentLockedAmount := new(big.Int).Set(p.LockedBalance)

	for i := range records {
		switch records[i].Kind {
		case Sum:
			currentAmount.Add(currentAmount, records[i].Value)
		case Sub:
			currentAmount.Sub(currentAmount, records[i].Value)
		case Lock:
			currentAmount.Sub(currentAmount, records[i].Value)
			currentLockedAmount.Add(currentLockedAmount, records[i].Value)
		case Unlock:
			currentLockedAmount.Sub(currentLockedAmount, records[i].Value)
			currentAmount.Add(currentAmount, records[i].Value)
		default:
			panic("")
		}
	}

	if currentAmount.Cmp(big.NewInt(0)) == -1 {
		return errors.New("")
	}

	if currentLockedAmount.Cmp(big.NewInt(0)) == -1 {
		return errors.New("")
	}

	p.SetBalance(currentAmount, currentLockedAmount)

	return nil
}

func (p *WalletRecord) Apply() error {
	return p.Wallet.Apply(p)
}

func NewTransaction(records ...*WalletRecord) *WalletTransaction {
	return &WalletTransaction{
		ID:            uuid.NewV4(),
		WalletRecords: records,
		Status:        Pending,
		CreatedAt:     time.Now(),
	}
}

func (p *WalletTransaction) Apply() error {
	for i := range p.WalletRecords {
		if err := p.WalletRecords[i].Apply(); err != nil {
			return err
		}
	}

	p.Status = Complete

	return nil
}
