package pelucio

import (
	"math/big"
	"testing"

	"github.com/devmalloni/pelucio/x/xuuid"

	"github.com/gofrs/uuid/v5"
)

func TestLedger_Compute(t *testing.T) {
	cashAccount := xuuid.MustParseString("f66263a4-8d54-413d-a133-32a65edec78a")
	clientAccount := xuuid.MustParseString("57557e00-fe68-47a6-a953-dcdc074c2caf")
	anotherClientAccount := xuuid.MustParseString("0c2f9b02-b6e3-43f8-990b-39608b28d365")

	transactions := []*Transaction{
		Deposit("deposit-1", cashAccount, clientAccount, big.NewInt(10), "BRL"),
		Withdraw("withdraw-1", cashAccount, clientAccount, big.NewInt(5), "BRL"),
		Deposit("deposit-2", cashAccount, clientAccount, big.NewInt(50), "USD"),
		Deposit("deposit-3", cashAccount, anotherClientAccount, big.NewInt(10), "BRL"),
		Withdraw("withdraw-2", cashAccount, anotherClientAccount, big.NewInt(3), "BRL"),
		Deposit("deposit-4", cashAccount, anotherClientAccount, big.NewInt(10), "EUR"),
		Transfer("transfer-1", anotherClientAccount, clientAccount, big.NewInt(5), "BRL"),
	}

	ledger, err := ComputeLedger(transactions)

	if err != nil {
		t.Fatalf("expected no error. got %v", err)
	}

	expectedBalancesPerAccount := map[uuid.UUID]Balance{
		cashAccount: {
			"BRL": big.NewInt(12),
			"EUR": big.NewInt(10),
			"USD": big.NewInt(50),
		},
		clientAccount: {
			"BRL": big.NewInt(10),
			"USD": big.NewInt(50),
		},
		anotherClientAccount: {
			"BRL": big.NewInt(2),
			"EUR": big.NewInt(10),
		},
	}

	if len(expectedBalancesPerAccount) != len(ledger.BalancesOfAccounts) {
		t.Errorf("expected %v balances per account; got %v", len(expectedBalancesPerAccount), len(ledger.BalancesOfAccounts))
	}

	for accountID, balances := range expectedBalancesPerAccount {
		if len(balances) != len(ledger.BalancesOfAccounts[accountID]) {
			t.Fatalf("expected len %v for account id %v; got %v", len(balances), accountID, len(ledger.BalancesOfAccounts[accountID]))
		}

		for currency, balance := range balances {
			if balance.Cmp(ledger.BalancesOfAccounts[accountID][currency]) != 0 {
				t.Errorf("on currency %v account %v expected %v but got %v", currency, accountID, balance, ledger.BalancesOfAccounts[accountID][currency])
			}
		}

	}

	if !ledger.Isbalanced {
		t.Errorf("expected ledger to be balanced, but got false.")
	}
}
