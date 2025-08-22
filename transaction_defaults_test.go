package pelucio

import (
	"math/big"
	"testing"

	"github.com/gofrs/uuid/v5"
)

func TestDeposit(t *testing.T) {
	d := Deposit("ext-123", uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), big.NewInt(1000), "USD")

	if !d.IsBalanced() {
		t.Errorf("Deposit transaction is not balanced: %v", d)
	}
}

func TestDepositWithFee(t *testing.T) {
	d := DepositWithFee("ext-123", uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), 10, big.NewInt(1000), "USD")

	if !d.IsBalanced() {
		t.Errorf("DepositWithFee transaction is not balanced: %v", d)
	}

	if len(d.Entries) != 3 {
		t.Errorf("DepositWithFee should have 3 entries, got %d", len(d.Entries))
	}

	if d.Entries[0].Amount.Cmp(big.NewInt(1000)) != 0 {
		t.Errorf("DepositWithFee cash account entry amount should be 1000, got %v", d.Entries[0].Amount)
	}
	if d.Entries[0].EntrySide != Debit {
		t.Errorf("DepositWithFee cash account entry side should be Debit, got %v", d.Entries[0].EntrySide)
	}

	if d.Entries[1].Amount.Cmp(big.NewInt(100)) != 0 {
		t.Errorf("DepositWithFee cash account entry amount should be 1000, got %v", d.Entries[0].Amount)
	}
	if d.Entries[1].EntrySide != Credit {
		t.Errorf("DepositWithFee cash account entry side should be Debit, got %v", d.Entries[0].EntrySide)
	}

	if d.Entries[2].Amount.Cmp(big.NewInt(900)) != 0 {
		t.Errorf("DepositWithFee cash account entry amount should be 1000, got %v", d.Entries[0].Amount)
	}
	if d.Entries[2].EntrySide != Credit {
		t.Errorf("DepositWithFee cash account entry side should be Debit, got %v", d.Entries[0].EntrySide)
	}
}

func TestWithdraw(t *testing.T) {
	d := Withdraw("ext-123", uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), big.NewInt(1000), "USD")

	if !d.IsBalanced() {
		t.Errorf("Deposit transaction is not balanced: %v", d)
	}
}

func TestWithdrawWithFee(t *testing.T) {
	d := WithdrawWithFee("ext-123", uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), 10, big.NewInt(1000), "USD")
	if !d.IsBalanced() {
		t.Errorf("WithdrawWithFee transaction is not balanced: %v", d)
	}

	if len(d.Entries) != 3 {
		t.Errorf("WithdrawWithFee should have 3 entries, got %d", len(d.Entries))
	}

	if d.Entries[0].Amount.Cmp(big.NewInt(900)) != 0 {
		t.Errorf("WithdrawWithFee cash account entry amount should be 1000, got %v", d.Entries[0].Amount)
	}
	if d.Entries[0].EntrySide != Credit {
		t.Errorf("WithdrawWithFee cash account entry side should be Debit, got %v", d.Entries[0].EntrySide)
	}

	if d.Entries[1].Amount.Cmp(big.NewInt(100)) != 0 {
		t.Errorf("WithdrawWithFee cash account entry amount should be 1000, got %v", d.Entries[0].Amount)
	}
	if d.Entries[1].EntrySide != Credit {
		t.Errorf("WithdrawWithFee cash account entry side should be Debit, got %v", d.Entries[0].EntrySide)
	}

	if d.Entries[2].Amount.Cmp(big.NewInt(1000)) != 0 {
		t.Errorf("WithdrawWithFee cash account entry amount should be 1000, got %v", d.Entries[0].Amount)
	}
	if d.Entries[2].EntrySide != Debit {
		t.Errorf("WithdrawWithFee cash account entry side should be Debit, got %v", d.Entries[0].EntrySide)
	}
}

func TestTrade(t *testing.T) {
	d := Trade("ext-123", TradeTransaction{
		TakerAccountID:     uuid.Must(uuid.NewV4()),
		TakerAsset:         "USD",
		TakerAmount:        big.NewInt(1000),
		TakerFeePercentage: 10,
		MakerAccountID:     uuid.Must(uuid.NewV4()),
		MakerAsset:         "EUR",
		MakerAmount:        big.NewInt(900),
		MakerFeePercentage: 5,
		FeeAccountID:       uuid.Must(uuid.NewV4()),
	})

	if !d.IsBalanced() {
		t.Errorf("Trade transaction is not balanced: %v", d)
	}
}

func TestTrade_WithoutFee(t *testing.T) {
	d := Trade("ext-123", TradeTransaction{
		TakerAccountID: uuid.Must(uuid.NewV4()),
		TakerAsset:     "USD",
		TakerAmount:    big.NewInt(1000),
		MakerAccountID: uuid.Must(uuid.NewV4()),
		MakerAsset:     "EUR",
		MakerAmount:    big.NewInt(900),
	})

	if !d.IsBalanced() {
		t.Errorf("Trade transaction is not balanced: %v", d)
	}
}
