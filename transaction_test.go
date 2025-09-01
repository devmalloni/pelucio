package pelucio

import (
	"math/big"
	"testing"

	"github.com/devmalloni/pelucio/x/xtime"
	"github.com/devmalloni/pelucio/x/xuuid"

	"github.com/gofrs/uuid/v5"
)

func TestTransaction_Balances(t *testing.T) {

}

func TestTransaction_BalancesByAccount(t *testing.T) {

}

func TestTransaction_IsBalanced(t *testing.T) {

}

func TestTransaction_ApplyToAccounts(t *testing.T) {
	transaction := &Transaction{
		ID: xuuid.MustParseString("e6a4f6b1-649a-4bcb-a8ea-43dc82d5e7fd"),
	}
	err := transaction.ApplyToAccounts(nil, xtime.DefaultClock)
	if err != ErrNoAccountProvided {
		t.Fatalf("expected ErrNoAccountProvided")
	}

	firstAccount := NewAccount(xtime.DefaultClock, WithNormalSide(Credit))
	accounts := map[uuid.UUID]*Account{
		firstAccount.ID: nil,
	}
	accounts[firstAccount.ID] = firstAccount

	err = transaction.ApplyToAccounts(accounts, xtime.DefaultClock)
	if err != ErrEntriesNotFound {
		t.Fatalf("expected ErrEntriesNotFound")
	}

	transaction.Entries = append(transaction.Entries, &Entry{
		ID:          xuuid.MustParseString("962fec36-3a62-4c48-ab10-c5d3f62135c0"),
		AccountID:   xuuid.MustParseString("212287c7-cfff-4fff-b94a-33e1a431fa51"),
		AccountSide: Credit,
		EntrySide:   Credit,
		Amount:      big.NewInt(1),
	})

	err = transaction.ApplyToAccounts(accounts, xtime.DefaultClock)
	if err != ErrTransactionIsNotBalanced {
		t.Fatalf("expected ErrTransactionIsNotBalanced")
	}

	transaction.Entries = append(transaction.Entries, &Entry{
		ID:            xuuid.MustParseString("54233be4-0ae1-4286-8367-260ca9089a5d"),
		TransactionID: transaction.ID,
		AccountID:     firstAccount.ID,
		AccountSide:   Credit,
		EntrySide:     Debit,
		Amount:        big.NewInt(1),
	})

	err = transaction.ApplyToAccounts(accounts, xtime.DefaultClock)
	if err != ErrAccountNotFound {
		t.Fatalf("expected ErrAccountNotFound")
	}
	transaction.Entries[0].AccountID = firstAccount.ID

	err = transaction.ApplyToAccounts(accounts, xtime.DefaultClock)
	if err != ErrEntryTransactionMismatch {
		t.Fatalf("expected ErrEntryTransactionMismatch. got %v", err)
	}

	transaction.Entries[0].TransactionID = transaction.ID

	transaction.Entries[0].Amount = big.NewInt(1)
	err = transaction.ApplyToAccounts(accounts, xtime.DefaultClock)
	if err != nil {
		t.Fatalf("expected error nil. got %v", err)
	}
}

func TestTransaction_ApplyToAccounts_ErrOnApply(t *testing.T) {
	transaction := &Transaction{
		ID: xuuid.MustParseString("e6a4f6b1-649a-4bcb-a8ea-43dc82d5e7fd"),
	}
	err := transaction.ApplyToAccounts(nil, xtime.DefaultClock)
	if err != ErrNoAccountProvided {
		t.Fatalf("expected ErrNoAccountProvided")
	}

	firstAccount := NewAccount(xtime.DefaultClock, WithNormalSide(Credit))
	accounts := map[uuid.UUID]*Account{
		firstAccount.ID: nil,
	}
	accounts[firstAccount.ID] = firstAccount
	transaction.Entries = append(transaction.Entries, &Entry{
		ID:            xuuid.MustParseString("962fec36-3a62-4c48-ab10-c5d3f62135c0"),
		TransactionID: transaction.ID,
		AccountID:     firstAccount.ID,
		AccountSide:   Credit,
		EntrySide:     Credit,
		Amount:        big.NewInt(0),
	})

	err = transaction.ApplyToAccounts(accounts, xtime.DefaultClock)
	if err != ErrNotPositiveAmount {
		t.Fatalf("expected ErrNotPositiveAmount. got %v", err)
	}
}

func TestTransaction_Accounts(t *testing.T) {
	transaction := Transaction{
		Entries: []*Entry{
			{
				AccountID: xuuid.MustParseString("276cb046-8f96-4239-bca3-0f8f07fadfec"),
			},
			{
				AccountID: xuuid.MustParseString("fbac3400-8207-48a1-bd0f-d641b7926282"),
			},
		},
	}

	res := transaction.Accounts()
	if len(res) != 2 {
		t.Errorf("expected len of accounts to be 2. got %v", len(res))
	}

	for i, e := range transaction.Entries {
		if !xuuid.Equal(res[i], e.AccountID) {
			t.Errorf("expected account id %v. got %v", e.AccountID, res[i])
		}
	}
}

func TestTransaction_Reverse(t *testing.T) {
	transaction := Deposit("externalID",
		xuuid.MustParseString("dd926a87-dff8-40b9-9403-d980fb2a2a0f"),
		xuuid.MustParseString("2080aa0b-2eee-46b4-ada4-02381fe4e9ed"),
		big.NewInt(1),
		"BRL")

	reversedTransaction := transaction.Reverse("reverse", "description", xtime.DefaultClock)

	if reversedTransaction.ExternalID != "reverse" {
		t.Errorf("external id expected to be %v. got %v", reversedTransaction.ExternalID, "reverse")
	}

	if reversedTransaction.Description != "description" {
		t.Errorf("description expected to be %v. got %v", reversedTransaction.Description, "description")
	}

	if !reversedTransaction.IsBalanced() {
		t.Errorf("expected reversed transaction to be balanced")
	}

	if reversedTransaction.Entries[0].EntrySide == transaction.Entries[0].EntrySide {
		t.Errorf("first entry side expected to be %v. got %v", reversedTransaction.Entries[0].EntrySide, transaction.Entries[0].EntrySide)
	}

	if reversedTransaction.Entries[1].EntrySide == transaction.Entries[1].EntrySide {
		t.Errorf("second entry side expected to be %v. got %v", reversedTransaction.Entries[1].EntrySide, transaction.Entries[1].EntrySide)
	}
}

func TestTransaction_SideByAccounts(t *testing.T) {
	transaction := Transaction{
		Entries: []*Entry{
			{
				AccountID:   xuuid.MustParseString("276cb046-8f96-4239-bca3-0f8f07fadfec"),
				AccountSide: Credit,
			},
			{
				AccountID:   xuuid.MustParseString("fbac3400-8207-48a1-bd0f-d641b7926282"),
				AccountSide: Debit,
			},
			{
				AccountID:   xuuid.MustParseString("fbac3400-8207-48a1-bd0f-d641b7926282"),
				AccountSide: Debit,
			},
		},
	}

	res := transaction.SideByAccounts()
	if len(res) != 2 {
		t.Errorf("expected len of accounts to be 2. got %v", len(res))
	}

	for _, e := range transaction.Entries {
		if res[e.AccountID] != e.AccountSide {
			t.Errorf("expected account side %v. got %v", e.AccountSide, res[e.AccountID])
		}
	}
}

func TestTransaction_SideByAccounts_AccountSideError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic as expected")
		} else if r != "there are multiple entries with different account side" {
			t.Errorf("Unexpected panic value: %v", r)
		}
	}()

	transaction := Transaction{
		Entries: []*Entry{
			{
				AccountID:   xuuid.MustParseString("fbac3400-8207-48a1-bd0f-d641b7926282"),
				AccountSide: Credit,
			},
			{
				AccountID:   xuuid.MustParseString("fbac3400-8207-48a1-bd0f-d641b7926282"),
				AccountSide: Debit,
			},
		},
	}

	transaction.SideByAccounts()

}
