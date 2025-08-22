package pelucio

import (
	"math/big"
	"pelucio/x/xuuid"
	"reflect"
	"testing"
	"time"
)

func TestEntry_Apply_ErrNilBalance(t *testing.T) {
	entry := Entry{}

	err := entry.Apply(nil)
	if err != ErrNilBalance {
		t.Fatalf("expected error ErrNilBalance")
	}
}

func TestEntry_Apply_ErrOnAdd(t *testing.T) {
	entry := Entry{
		ID:          xuuid.MustParseString("6b3f5a20-4721-4cb5-8046-4af41e319736"),
		AccountID:   xuuid.MustParseString("cf1cdd05-3311-4e6b-8dbc-e0b0bcb59bde"),
		EntrySide:   Debit,
		AccountSide: Debit,
		Amount:      big.NewInt(-100),
		Currency:    "USD",
		CreatedAt:   time.Now(),
	}

	balance := Balance{}
	err := entry.Apply(balance)
	if err == nil {
		t.Fatalf("expected error on add operation. got nil")
	}
}

func TestEntry_Apply_DebitEntry_DebitAccount(t *testing.T) {
	balance := make(Balance)

	entry := Entry{
		ID:          xuuid.MustParseString("6b3f5a20-4721-4cb5-8046-4af41e319736"),
		AccountID:   xuuid.MustParseString("cf1cdd05-3311-4e6b-8dbc-e0b0bcb59bde"),
		EntrySide:   Debit,
		AccountSide: Debit,
		Amount:      big.NewInt(100),
		Currency:    "USD",
		CreatedAt:   time.Now(),
	}

	if err := entry.Apply(balance); err != nil {
		t.Errorf("Failed to apply entry: %v", err)
	}

	if balance["USD"] == nil {
		t.Error("Expected balance for USD to be set")
	}

	if balance["USD"].Cmp(big.NewInt(100)) != 0 {
		t.Errorf("Expected balance for USD to be 100, got %s", balance["USD"].String())
	}
}

func TestEntry_Apply_CreditEntry_DebitAccount(t *testing.T) {
	balance := make(Balance)

	balance["USD"] = big.NewInt(200)

	entry := Entry{
		ID:          xuuid.MustParseString("95b61c00-d121-4263-aeba-188ab345e493"),
		AccountID:   xuuid.MustParseString("5565cbb8-2104-4b51-a306-8304ff255b3e"),
		EntrySide:   Credit,
		AccountSide: Debit,
		Amount:      big.NewInt(100),
		Currency:    "USD",
		CreatedAt:   time.Now(),
	}

	if err := entry.Apply(balance); err != nil {
		t.Errorf("Failed to apply entry: %v", err)
	}

	if balance["USD"] == nil {
		t.Error("Expected balance for USD to be set")
	}

	if balance["USD"].Cmp(big.NewInt(100)) != 0 {
		t.Errorf("Expected balance for USD to be 100, got %s", balance["USD"].String())
	}
}

func TestEntry_Apply_CreditEntry_DebitAccount_NoBalance(t *testing.T) {
	balance := make(Balance)

	entry := Entry{
		ID:          xuuid.MustParseString("324b502b-371b-426c-b4d8-349c447f350a"),
		AccountID:   xuuid.MustParseString("1473232c-14b5-485e-80a6-2d9e7f74d3b1"),
		EntrySide:   Credit,
		AccountSide: Debit,
		Amount:      big.NewInt(100),
		Currency:    "USD",
		CreatedAt:   time.Now(),
	}

	if err := entry.Apply(balance); err == nil {
		t.Fatal("Expected error. Got nil")
	}
}

func TestEntry_Apply_CreditEntry_CreditAccount(t *testing.T) {
	balance := make(Balance)

	entry := Entry{
		ID:          xuuid.MustParseString("f2e1ee61-f6c1-495d-82a7-b95182169e97"),
		AccountID:   xuuid.MustParseString("fc947d4c-508c-4763-bb07-6f023f02ec08"),
		EntrySide:   Credit,
		AccountSide: Credit,
		Amount:      big.NewInt(100),
		Currency:    "USD",
		CreatedAt:   time.Now(),
	}

	if err := entry.Apply(balance); err != nil {
		t.Errorf("Failed to apply entry: %v", err)
	}

	if balance["USD"] == nil {
		t.Error("Expected balance for USD to be set")
	}

	if balance["USD"].Cmp(big.NewInt(100)) != 0 {
		t.Errorf("Expected balance for USD to be 100, got %s", balance["USD"].String())
	}
}

func TestEntry_Apply_DebitEntry_CreditAccount(t *testing.T) {
	balance := make(Balance)

	balance["USD"] = big.NewInt(200)

	entry := Entry{
		ID:          xuuid.MustParseString("cc2524ff-cf98-44d5-830c-e26ccccf02cb"),
		AccountID:   xuuid.MustParseString("dcbb8ac1-9206-454f-b93b-c655a6557bee"),
		EntrySide:   Debit,
		AccountSide: Credit,
		Amount:      big.NewInt(100),
		Currency:    "USD",
		CreatedAt:   time.Now(),
	}

	if err := entry.Apply(balance); err != nil {
		t.Errorf("Failed to apply entry: %v", err)
	}

	if balance["USD"] == nil {
		t.Error("Expected balance for USD to be set")
	}

	if balance["USD"].Cmp(big.NewInt(100)) != 0 {
		t.Errorf("Expected balance for USD to be 100, got %s", balance["USD"].String())
	}
}

func TestEntry_Apply_DebitEntry_CreditAccount_NoBalance(t *testing.T) {
	balance := make(Balance)

	entry := Entry{
		ID:          xuuid.MustParseString("da841ebd-0d99-4751-9121-568ad6b7d444"),
		AccountID:   xuuid.MustParseString("e99ecee2-4c35-44be-9724-8a1291af09ca"),
		EntrySide:   Debit,
		AccountSide: Credit,
		Amount:      big.NewInt(100),
		Currency:    "USD",
		CreatedAt:   time.Now(),
	}

	if err := entry.Apply(balance); err == nil {
		t.Fatal("Expected error. Got nil")
	}
}

func TestBalance_IsBalanced_False_Less_Currencies(t *testing.T) {
	positive := Balance{
		Currency("BRL"):  big.NewInt(100),
		Currency("USD"):  big.NewInt(200),
		Currency("EUR"):  big.NewInt(300),
		Currency("USDC"): big.NewInt(0),
	}

	negative := Balance{
		Currency("BRL"): big.NewInt(100),
		Currency("USD"): big.NewInt(200),
	}

	if positive.IsBalanced(negative) {
		t.Error("Expected positive balance to be balanced")
	}
}

func TestBalance_IsBalanced_False_DifferentAmount_Currencies(t *testing.T) {
	positive := Balance{
		Currency("BRL"):  big.NewInt(100),
		Currency("USD"):  big.NewInt(200),
		Currency("EUR"):  big.NewInt(300),
		Currency("USDC"): big.NewInt(0),
	}

	negative := Balance{
		Currency("BRL"): big.NewInt(200),
		Currency("USD"): big.NewInt(5),
		Currency("EUR"): big.NewInt(300),
	}

	if positive.IsBalanced(negative) {
		t.Error("Expected positive balance to be balanced")
	}
}

func TestBalance_IsBalanced_True(t *testing.T) {
	positive := Balance{
		Currency("BRL"):  big.NewInt(100),
		Currency("USD"):  big.NewInt(200),
		Currency("EUR"):  big.NewInt(300),
		Currency("USDC"): big.NewInt(0),
	}

	negative := Balance{
		Currency("BRL"):  big.NewInt(100),
		Currency("USD"):  big.NewInt(200),
		Currency("EUR"):  big.NewInt(300),
		Currency("USDC"): big.NewInt(0),
	}

	if !positive.IsBalanced(negative) {
		t.Error("Expected positive balance to be balanced")
	}
}

func TestEntry_OperationOnBalance(t *testing.T) {
	debitOnCreditAccountEntry := Entry{
		EntrySide:   Debit,
		AccountSide: Credit,
	}
	if debitOnCreditAccountEntry.OperationOnBalance() != OperationKindSub {
		t.Fatalf("debit on credit account entry expected sub")
	}

	creditOnCreditAccountEntry := Entry{
		EntrySide:   Credit,
		AccountSide: Credit,
	}
	if creditOnCreditAccountEntry.OperationOnBalance() != OperationKindAdd {
		t.Fatalf("credit on credit account entry expected add")
	}

	debitOnDebitAccountEntry := Entry{
		EntrySide:   Debit,
		AccountSide: Debit,
	}
	if debitOnDebitAccountEntry.OperationOnBalance() != OperationKindAdd {
		t.Fatalf("debit on credit account entry expected sub")
	}

	creditOnDebitAccountEntry := Entry{
		EntrySide:   Credit,
		AccountSide: Debit,
	}
	if creditOnDebitAccountEntry.OperationOnBalance() != OperationKindSub {
		t.Fatalf("credit on credit account entry expected sub")
	}
}

func TestEntry_OperationOnBalance_Invalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic as expected")
		} else if r != "invalid entry side or account side" {
			t.Errorf("Unexpected panic value: %v", r)
		}
	}()

	debitOnCreditAccountEntry := Entry{
		EntrySide:   "foo",
		AccountSide: "bar",
	}

	debitOnCreditAccountEntry.OperationOnBalance()
}

func TestEntry_Reverse_CreditAccount_DebitSide(t *testing.T) {
	originalEntry := Entry{
		ID:            xuuid.MustParseString("63e11bba-3daf-44a7-a263-399a665ab699"),
		TransactionID: xuuid.MustParseString("0fe5ba57-74db-4d80-9f00-18c037c98196"),
		AccountID:     xuuid.MustParseString("d2b88dbf-925a-489c-ac67-2fa9c4c4a63f"),
		EntrySide:     Debit,
		AccountSide:   Credit,
		Amount:        big.NewInt(1),
		Currency:      "BRL",
	}

	reversedEntry := originalEntry.Reverse(xuuid.MustParseString("e4cb4476-39e2-402f-886d-df60767a5924"), stubClock)

	if reversedEntry.EntrySide != Credit {
		t.Fatalf("expected reversed entry side to be credit")
	}

	if reversedEntry.Amount.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("expected reversed entry amount to be 1")
	}

	if reversedEntry.AccountID != originalEntry.AccountID {
		t.Fatalf("expected reversed entry account to be the same as original")
	}

	if reversedEntry.AccountSide != originalEntry.AccountSide {
		t.Fatalf("expected reversed entry account side to be the same as original")
	}

	if reversedEntry.Currency != originalEntry.Currency {
		t.Fatalf("expected reversed entry currency to be the same as original")
	}

	if !xuuid.Equal(reversedEntry.TransactionID, xuuid.MustParseString("e4cb4476-39e2-402f-886d-df60767a5924")) {
		t.Fatalf("expected reversed entry transaction id to be foo")
	}
}

func TestEntry_Reverse_CreditAccount_CreditSide(t *testing.T) {
	originalEntry := Entry{
		ID:            xuuid.MustParseString("bf8dcd72-8743-4f4f-8f68-7af664719bd8"),
		TransactionID: xuuid.MustParseString("ad68c516-738a-4b74-9712-397971923f08"),
		AccountID:     xuuid.MustParseString("67553216-35bf-4c1a-bd71-12a2ddf80801"),
		EntrySide:     Credit,
		AccountSide:   Credit,
		Amount:        big.NewInt(1),
		Currency:      "BRL",
	}

	reversedEntry := originalEntry.Reverse(xuuid.MustParseString("da494476-041b-4a7d-a41c-2b9129cec160"), stubClock)

	if reversedEntry.EntrySide != Debit {
		t.Fatalf("expected reversed entry side to be debit")
	}

	if reversedEntry.Amount.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("expected reversed entry amount to be 1")
	}

	if reversedEntry.AccountID != originalEntry.AccountID {
		t.Fatalf("expected reversed entry account to be the same as original")
	}

	if reversedEntry.AccountSide != originalEntry.AccountSide {
		t.Fatalf("expected reversed entry account side to be the same as original")
	}

	if reversedEntry.Currency != originalEntry.Currency {
		t.Fatalf("expected reversed entry currency to be the same as original")
	}

	if !xuuid.Equal(reversedEntry.TransactionID, xuuid.MustParseString("da494476-041b-4a7d-a41c-2b9129cec160")) {
		t.Fatalf("expected reversed entry transaction id to be foo")
	}
}

func TestEntry_Reverse_DebitAccount_CreditSide(t *testing.T) {
	originalEntry := Entry{
		ID:            xuuid.MustParseString("ccc57fd5-d24d-43c6-8d2e-55981653e0f7"),
		TransactionID: xuuid.MustParseString("490908ee-4cfc-4680-aa20-1499b1580603"),
		AccountID:     xuuid.MustParseString("7a8ef1c7-9f1e-42a9-972a-2cd2bd4c3c1e"),
		EntrySide:     Credit,
		AccountSide:   Debit,
		Amount:        big.NewInt(1),
		Currency:      "BRL",
	}

	reversedEntry := originalEntry.Reverse(xuuid.MustParseString("e759ec41-4836-4201-ab3c-6ef13c75b703"), stubClock)

	if reversedEntry.EntrySide != Debit {
		t.Fatalf("expected reversed entry side to be debit")
	}

	if reversedEntry.Amount.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("expected reversed entry amount to be 1")
	}

	if reversedEntry.AccountID != originalEntry.AccountID {
		t.Fatalf("expected reversed entry account to be the same as original")
	}

	if reversedEntry.AccountSide != originalEntry.AccountSide {
		t.Fatalf("expected reversed entry account side to be the same as original")
	}

	if reversedEntry.Currency != originalEntry.Currency {
		t.Fatalf("expected reversed entry currency to be the same as original")
	}

	if !xuuid.Equal(reversedEntry.TransactionID, xuuid.MustParseString("e759ec41-4836-4201-ab3c-6ef13c75b703")) {
		t.Fatalf("expected reversed entry transaction id to be foo")
	}
}

func TestEntry_Reverse_DebitAccount_DebitSide(t *testing.T) {
	originalEntry := Entry{
		ID:            xuuid.MustParseString("fe6ac534-ec89-4ae3-ab50-e4a87d378683"),
		TransactionID: xuuid.MustParseString("78dc8432-8246-4cf3-9970-27b582f97357"),
		AccountID:     xuuid.MustParseString("ee5c723c-cdd5-4503-94d6-252da84aaca0"),
		EntrySide:     Debit,
		AccountSide:   Debit,
		Amount:        big.NewInt(1),
		Currency:      "BRL",
	}

	reversedEntry := originalEntry.Reverse(xuuid.MustParseString("20bc3b14-339a-42f0-a210-a5a737de6d52"), stubClock)

	if reversedEntry.EntrySide != Credit {
		t.Fatalf("expected reversed entry side to be debit")
	}

	if reversedEntry.Amount.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("expected reversed entry amount to be 1")
	}

	if reversedEntry.AccountID != originalEntry.AccountID {
		t.Fatalf("expected reversed entry account to be the same as original")
	}

	if reversedEntry.AccountSide != originalEntry.AccountSide {
		t.Fatalf("expected reversed entry account side to be the same as original")
	}

	if reversedEntry.Currency != originalEntry.Currency {
		t.Fatalf("expected reversed entry currency to be the same as original")
	}

	if !xuuid.Equal(reversedEntry.TransactionID, xuuid.MustParseString("20bc3b14-339a-42f0-a210-a5a737de6d52")) {
		t.Fatalf("expected reversed entry transaction id to be foo")
	}
}

func TestBalance_Add(t *testing.T) {
	balance := make(Balance)

	err := balance.Add("brl", big.NewInt(0))
	if err != ErrNotPositiveAmount {
		t.Fatalf("error is expected")
	}

	err = balance.Add("brl", big.NewInt(-1))
	if err != ErrNotPositiveAmount {
		t.Fatalf("error is expected")
	}

	err = balance.Add("brl", big.NewInt(1))
	if err != nil {
		t.Fatalf("expected to be added 1, got %v", err)
	}
}

func TestBalance_Sub(t *testing.T) {
	balance := make(Balance)

	err := balance.Sub("brl", big.NewInt(0))
	if err != ErrNotPositiveAmount {
		t.Fatalf("error is expected")
	}

	err = balance.Sub("brl", big.NewInt(-1))
	if err != ErrNotPositiveAmount {
		t.Fatalf("error is expected")
	}

	err = balance.Sub("brl", big.NewInt(1))
	if err != ErrInsufficientBalance {
		t.Fatalf("expected insufficient balance")
	}

	balance.Add("brl", big.NewInt(1))

	err = balance.Sub("brl", big.NewInt(2))
	if err != ErrInsufficientBalance {
		t.Fatalf("expected insufficient balance")
	}

	err = balance.Sub("brl", big.NewInt(1))
	if err != nil {
		t.Fatalf("expected to be added 1, got %v", err)
	}
}

func TestBalance_HasBalance(t *testing.T) {
	b := Balance{
		"brl": big.NewInt(0),
		"usd": big.NewInt(0),
	}

	if b.HasBalance() {
		t.Fatalf("expected has balance to be false")
	}

	b.Add("brl", big.NewInt(1))

	if !b.HasBalance() {
		t.Fatalf("expected has balance to be true")
	}
}

func TestBalance_HasBalance_EmptyBalance(t *testing.T) {
	b := Balance{}

	if b.HasBalance() {
		t.Fatalf("expected has balance to be false")
	}
}

func TestBalance_Decimal(t *testing.T) {
	b := Balance{
		"brl": big.NewInt(12345), // 123.45
		"usd": big.NewInt(6789),  // 67.89
	}

	decimals := b.Decimal(4)

	expected := map[Currency]string{
		"brl": "1.2345",
		"usd": "0.6789",
	}

	if !reflect.DeepEqual(decimals, expected) {
		t.Fatalf("expected decimals to be %v, got %v", expected, decimals)
	}
}

func TestBalance_DecimalFromMap(t *testing.T) {
	b := Balance{
		"brl": big.NewInt(12345), // 123.45
		"usd": big.NewInt(6789),  // 67.89
		"btc": big.NewInt(1),     // 67.89
	}

	decimals := b.DecimalFromMap(map[Currency]int{
		"brl": 4,
		"usd": 2,
	}, 8)

	expected := map[Currency]string{
		"brl": "1.2345",
		"usd": "67.89",
		"btc": "0.00000001",
	}

	if !reflect.DeepEqual(decimals, expected) {
		t.Fatalf("expected decimals to be %v, got %v", expected, decimals)
	}
}
