package pelucio

import (
	"math/big"

	"github.com/gofrs/uuid/v5"
)

// TransferBetweenCreditAccounts creates a transaction for transferring an amount from one credit account to another.
func TransferBetweenCreditAccounts(externalID string, fromAccountID, toAccountID uuid.UUID, amount *big.Int, currency Currency) *Transaction {
	return NewTransaction(nil).
		WithExternalID(externalID).
		AddEntry(fromAccountID, Debit, Credit, amount, currency).
		AddEntry(toAccountID, Credit, Credit, amount, currency).
		MustBuild()
}

// TransferBetweenDebitAccounts creates a transaction for transferring an amount from one debit account to another.
func TransferBetweenDebitAccounts(externalID string, fromAccountID, toAccountID uuid.UUID, amount *big.Int, currency Currency) *Transaction {
	return NewTransaction(nil).
		WithExternalID(externalID).
		AddEntry(fromAccountID, Credit, Debit, amount, currency).
		AddEntry(toAccountID, Debit, Credit, amount, currency).
		MustBuild()
}

// Deposit creates a transaction for depositing an amount into a credit account from a cash account with debit normal side.
func Deposit(externalID string, cashAccountID, toAccountID uuid.UUID, amount *big.Int, currency Currency) *Transaction {
	return NewTransaction(nil).
		WithExternalID(externalID).
		AddEntry(cashAccountID, Debit, Debit, amount, currency).
		AddEntry(toAccountID, Credit, Credit, amount, currency).
		MustBuild()
}

// DepositWithFee its like Deposit but takes part of the amount as fee which is transferred to a fee account.
func DepositWithFee(externalID string, cashAccountID, feeAccount, toAccountID uuid.UUID, feeTax int8, amount *big.Int, currency Currency) *Transaction {
	depositFee := TakePercent(amount, feeTax)

	return NewTransaction(nil).
		WithExternalID(externalID).
		AddEntry(cashAccountID, Debit, Debit, amount, currency).
		AddEntry(feeAccount, Credit, Credit, depositFee, currency).
		AddEntry(toAccountID, Credit, Credit, new(big.Int).Sub(amount, depositFee), currency).
		MustBuild()
}

// Withdraw creates a transaction for withdrawing an amount from a credit account and a cash account with debit normal side.
func Withdraw(externalID string, cashAccountID, fromAccount uuid.UUID, amount *big.Int, currency Currency) *Transaction {
	return NewTransaction(nil).
		WithExternalID(externalID).
		AddEntry(cashAccountID, Credit, Debit, amount, currency).
		AddEntry(fromAccount, Debit, Credit, amount, currency).
		MustBuild()
}

// WithdrawWithFee is the same as Withdraw but takes part of the amount as fee which is transferred to a fee account.
func WithdrawWithFee(externalID string, cashAccountID, feeAccount, fromAccount uuid.UUID, feeTax int8, amount *big.Int, currency Currency) *Transaction {
	withdrawFee := TakePercent(amount, feeTax)

	return NewTransaction(nil).
		WithExternalID(externalID).
		AddEntry(cashAccountID, Credit, Debit, new(big.Int).Sub(amount, withdrawFee), currency).
		AddEntry(feeAccount, Credit, Credit, withdrawFee, currency).
		AddEntry(fromAccount, Debit, Credit, amount, currency).
		MustBuild()
}

type TradeTransaction struct {
	ExternalID string

	// Taker account ID
	TakerAccountID uuid.UUID
	// Taker asset to be removed from taker account and added to maker account
	TakerAsset Currency
	// total amount to be debited from taker account
	TakerAmount *big.Int
	// Fee percentage to be removed from taker amount before adding to maker account
	TakerFeePercentage int8

	// Maker account ID
	MakerAccountID uuid.UUID
	// Maker asset to be removed from Maker account and added to taker account
	MakerAsset Currency
	// total amount to be debited from maker account
	MakerAmount *big.Int
	// Fee percentage to be removed from maker amount before adding to taker account
	MakerFeePercentage int8

	FeeAccountID uuid.UUID
}

func Trade(externalID string, trade TradeTransaction) *Transaction {
	takerFee := TakePercent(trade.TakerAmount, trade.TakerFeePercentage)
	makerFee := TakePercent(trade.MakerAmount, trade.MakerFeePercentage)

	// create transaction builder
	t := NewTransaction(nil).
		WithExternalID(externalID)

	// debit from taker account and add to maker account and fee
	t.
		AddEntry(trade.TakerAccountID, Debit, Credit, trade.TakerAmount, trade.TakerAsset).
		AddEntry(trade.MakerAccountID, Credit, Credit, new(big.Int).Sub(trade.TakerAmount, takerFee), trade.TakerAsset)
	if takerFee.Cmp(big.NewInt(0)) > 0 {
		t.AddEntry(trade.FeeAccountID, Credit, Credit, takerFee, trade.TakerAsset)
	}

	// debit from maker account and add to taker account and fee
	t.
		AddEntry(trade.MakerAccountID, Debit, Credit, trade.MakerAmount, trade.MakerAsset).
		AddEntry(trade.TakerAccountID, Credit, Credit, new(big.Int).Sub(trade.MakerAmount, makerFee), trade.MakerAsset)
	if makerFee.Cmp(big.NewInt(0)) > 0 {
		t.AddEntry(trade.FeeAccountID, Credit, Credit, makerFee, trade.MakerAsset)
	}

	return t.MustBuild()
}

func TakePercent(value *big.Int, percentage int8) *big.Int {
	if percentage == 0 {
		return big.NewInt(0)
	}

	if percentage < 0 {
		panic("percentage must be positive")
	}

	bPercentage := big.NewInt(int64(percentage))
	hundred := big.NewInt(100)

	valueWithPercentage := new(big.Int).Mul(value, bPercentage)
	return new(big.Int).Div(valueWithPercentage, hundred)
}
