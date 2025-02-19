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
	WalletCurrency          string
	WalletRecordKind        uint8
	WalletTransactionStatus uint8

	Wallet struct {
		ID         uuid.UUID `json:"id,omitempty"`
		ExternalID string    `json:"external_id,omitempty"`

		Balance       map[WalletCurrency]*big.Int `json:"balance,omitempty"`
		LockedBalance map[WalletCurrency]*big.Int `json:"lockedBalance,omitempty"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		Version   uuid.UUID  `json:"version,omitempty"`
	}

	WalletRecord struct {
		ID       uuid.UUID
		WalletID uuid.UUID
		Currency WalletCurrency
		Value    *big.Int
		Kind     WalletRecordKind

		Wallet    *Wallet
		CreatedAt time.Time
	}

	WalletTransaction struct {
		ID            uuid.UUID
		WalletRecords []*WalletRecord
		Status        WalletTransactionStatus
		CreatedAt     time.Time
	}

	WalletResponse struct {
		ID         uuid.UUID `json:"id,omitempty"`
		ExternalID string    `json:"external_id,omitempty"`

		Balance       map[WalletCurrency]*string `json:"balance,omitempty"`
		LockedBalance map[WalletCurrency]*string `json:"lockedBalance,omitempty"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		Version   uuid.UUID  `json:"version,omitempty"`
	}
)

func (p *Wallet) SetBalance(balance, lockedBalance map[WalletCurrency]*big.Int) error {
	for _, b := range balance {
		if b.Cmp(big.NewInt(0)) == -1 {
			return errors.New("negative balance not allowed")
		}
	}

	for _, b := range lockedBalance {
		if b.Cmp(big.NewInt(0)) == -1 {
			return errors.New("negative balance not allowed")
		}
	}

	p.Balance = balance
	p.LockedBalance = lockedBalance
	p.Version = uuid.NewV4()

	return nil
}

func (p *Wallet) Record(amount *big.Int, kind WalletRecordKind, currency WalletCurrency) *WalletRecord {
	return &WalletRecord{
		ID:        uuid.NewV4(),
		WalletID:  p.ID,
		Currency:  currency,
		Value:     amount,
		Kind:      kind,
		CreatedAt: time.Now(),
		Wallet:    p,
	}
}

func (p *Wallet) Add(amount *big.Int, currency WalletCurrency) *WalletRecord {
	return p.Record(amount, Sum, currency)
}

func (p *Wallet) Sub(amount *big.Int, currency WalletCurrency) *WalletRecord {
	return p.Record(amount, Sub, currency)
}

func (p *Wallet) Lock(amount *big.Int, currency WalletCurrency) *WalletRecord {
	return p.Record(amount, Lock, currency)
}

func (p *Wallet) Unlock(amount *big.Int, currency WalletCurrency) *WalletRecord {
	return p.Record(amount, Unlock, currency)
}

func (p *Wallet) Apply(records ...*WalletRecord) error {
	balances := p.Balance
	lockedBalances := p.LockedBalance

	for i := range records {
		currency := records[i].Currency
		_, found := balances[currency]
		if !found {
			balances[currency] = p.BalanceOf(currency)
		}
		currentBalance := balances[currency]

		_, found = lockedBalances[currency]
		if !found {
			lockedBalances[currency] = p.LockedBalanceOf(currency)
		}
		currentLockedBalance := lockedBalances[currency]

		switch records[i].Kind {
		case Sum:
			currentBalance.Add(currentBalance, records[i].Value)
		case Sub:
			currentBalance.Sub(currentBalance, records[i].Value)
		case Lock:
			currentBalance.Sub(currentBalance, records[i].Value)
			currentLockedBalance.Add(currentLockedBalance, records[i].Value)
		case Unlock:
			currentLockedBalance.Sub(currentLockedBalance, records[i].Value)
			currentBalance.Add(currentBalance, records[i].Value)
		default:
			panic("")
		}
	}

	err := p.SetBalance(balances, lockedBalances)
	if err != nil {
		return err
	}

	return nil
}

func (p *Wallet) BalanceOf(currency WalletCurrency) *big.Int {
	return p.balanceOf(p.Balance, currency)
}

func (p *Wallet) LockedBalanceOf(currency WalletCurrency) *big.Int {
	return p.balanceOf(p.LockedBalance, currency)
}

func (p *Wallet) balanceOf(currencyMap map[WalletCurrency]*big.Int, currency WalletCurrency) *big.Int {
	c, found := currencyMap[currency]
	if !found {
		return big.NewInt(0)
	}

	return new(big.Int).Set(c)
}

func (p *Wallet) ToWalletResponse() *WalletResponse {
	balance := make(map[WalletCurrency]*string)
	lockedBalance := make(map[WalletCurrency]*string)

	for k, v := range p.Balance {
		vv := v.String()
		balance[k] = &vv
	}

	for k, v := range p.LockedBalance {
		vv := v.String()
		lockedBalance[k] = &vv
	}

	return &WalletResponse{
		ID:            p.ID,
		Version:       p.Version,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
		Balance:       balance,
		LockedBalance: lockedBalance,
		ExternalID:    p.ExternalID,
	}
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
