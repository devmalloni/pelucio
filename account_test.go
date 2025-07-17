package pelucio

import (
	"math/big"
	"pelucio/x/xtime"
	"pelucio/x/xuuid"
	"testing"
	"time"
)

var stubClock = xtime.NewStubClock(time.Now())

func TestAccount_Apply_NotSameAccountID(t *testing.T) {
	a := Account{
		ID: xuuid.MustParseString("8e7dc374-4284-4a29-ae76-80b5214424df"),
	}

	e := Entry{
		ID:        xuuid.MustParseString("a3634261-f680-440a-aeaf-9c1c2ae1cd78"),
		AccountID: xuuid.MustParseString("2531668d-ef82-4ad3-affc-97b52e2e6d4d"),
	}

	err := a.Apply(e, stubClock)
	if err != ErrEntryAccountMismatch {
		t.Fatal("Expected error. Got nil")
	}
}

func TestAccount_Apply_NotExpectedNormalSide(t *testing.T) {
	a := Account{
		ID:         xuuid.MustParseString("61d289b2-3a59-42cf-8e6d-d1157a4507b0"),
		NormalSide: Debit,
	}

	e := Entry{
		ID:          xuuid.MustParseString("b9564cca-3aa9-4e0e-b5a1-f2dc35ac8439"),
		AccountID:   xuuid.MustParseString("61d289b2-3a59-42cf-8e6d-d1157a4507b0"),
		AccountSide: Credit,
	}

	err := a.Apply(e, stubClock)
	if err != ErrAccountSideMismatch {
		t.Fatal("Expected error. Got nil")
	}
}

func TestAccount_Apply_NoBalance(t *testing.T) {
	a := Account{
		ID:         xuuid.MustParseString("78ec1d47-0061-4fb1-bb97-028313b2f71a"),
		NormalSide: Credit,
	}

	e := Entry{
		ID:          xuuid.MustParseString("1b1fbdac-b6a2-4034-a5fc-96f603d27ff2"),
		AccountID:   xuuid.MustParseString("78ec1d47-0061-4fb1-bb97-028313b2f71a"),
		AccountSide: Credit,
		EntrySide:   Debit,
		Amount:      big.NewInt(100),
	}

	err := a.Apply(e, stubClock)
	if err != ErrInsufficientBalance {
		t.Fatal("Expected error. Got nil")
	}
}

func TestAccount_UpdateData(t *testing.T) {
	a := NewAccount("bar", "name", Credit, nil, stubClock)
	a.Balance = make(Balance)
	a.Balance.Add("brl", big.NewInt(1))

	a.UpdateData("foo", nil, stubClock)

	if a.Name != "foo" {
		t.Fatalf("name is expected to be foo, got %v", a.Name)
	}

	if a.Balance["brl"].Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("expected balance to not change")
	}
}

func TestAccount_ComputeFromEntries(t *testing.T) {
	a := &Account{
		ID:         xuuid.MustParseString("0e77b405-4ba6-410a-aa85-2e4c15d791e5"),
		NormalSide: Debit,
		Balance: Balance{
			"brl": big.NewInt(10),
		},
	}

	entries := []*Entry{
		{
			AccountID:   a.ID,
			EntrySide:   Debit,
			AccountSide: Debit,
			Amount:      big.NewInt(5),
			Currency:    "brl",
		},
		{
			AccountID:   a.ID,
			EntrySide:   Debit,
			AccountSide: Debit,
			Amount:      big.NewInt(5),
			Currency:    "usd",
		},
		{
			AccountID:   a.ID,
			EntrySide:   Credit,
			AccountSide: Debit,
			Amount:      big.NewInt(2),
			Currency:    "usd",
		},
	}

	err := a.ComputeFromEntries(entries, stubClock)
	if err != nil {
		t.Fatal("TestAccount_ComputeFromEntries: error found but none expected")
	}

	if len(a.Balance) != 2 {
		t.Fatalf("expect 2 currencies at balance. Found %v", len(a.Balance))
	}

	if a.Balance["brl"].Cmp(big.NewInt(5)) != 0 {
		t.Fatalf("expect brl to be 5. Found %v", a.Balance["brl"])
	}

	if a.Balance["usd"].Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("expect usd to be 3. Found %v", a.Balance["usd"])
	}
}

func TestAccount_ComputeFromEntries_Error(t *testing.T) {
	a := &Account{
		ID:         xuuid.MustParseString("0e77b405-4ba6-410a-aa85-2e4c15d791e5"),
		NormalSide: Debit,
		Balance: Balance{
			"brl": big.NewInt(10),
		},
	}

	entries := []*Entry{
		{
			AccountID:   xuuid.MustParseString("c5d1913c-959b-4502-b39b-9636a2c99559"),
			EntrySide:   Debit,
			AccountSide: Debit,
			Amount:      big.NewInt(5),
			Currency:    "brl",
		},
	}

	err := a.ComputeFromEntries(entries, stubClock)
	if err == nil {
		t.Fatal("TestAccount_ComputeFromEntries: error not found but one expected")
	}
}

func TestAccount_Delete_BalanceNotEmpty(t *testing.T) {
	a := &Account{
		Balance: Balance{
			"brl": big.NewInt(2),
		},
	}

	if err := a.Delete(stubClock); err != ErrBalanceNotEmpty {
		t.Fatal("expected error ErrBalanceNotEmpty")
	}
}

func TestAccount_Delete(t *testing.T) {
	a := &Account{}

	if err := a.Delete(stubClock); err != nil {
		t.Fatal("expected error nil")
	}
}
