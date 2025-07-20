package pelucio

import (
	"pelucio/x/xuuid"

	"github.com/gofrs/uuid/v5"
)

type (
	Ledger struct {
		ID                 uuid.UUID
		BalancesOfAccounts map[uuid.UUID]Balance
		Isbalanced         bool
	}
)

func ComputeLedger(transactions []*Transaction) (*Ledger, error) {
	l := &Ledger{
		ID: xuuid.New(),
	}

	balances, isBalanced := l.compute(transactions)
	l.BalancesOfAccounts = balances
	l.Isbalanced = isBalanced

	return l, nil
}

func (p *Ledger) compute(transactions []*Transaction) (map[uuid.UUID]Balance, bool) {
	balancesPerAccount := make(map[uuid.UUID]Balance)
	isBalanced := true
	for _, transaction := range transactions {
		transaction.BalancesByAccount(balancesPerAccount)
		isBalanced = isBalanced && transaction.IsBalanced()
	}

	return balancesPerAccount, isBalanced
}
