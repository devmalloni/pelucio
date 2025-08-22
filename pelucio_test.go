package pelucio

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/devmalloni/pelucio/x/xtime"
	"github.com/devmalloni/pelucio/x/xuuid"

	"github.com/stretchr/testify/mock"
)

func TestPelucio_CreateAccount_Success(t *testing.T) {
	readWriter := new(ReadWriterMock)
	pelucio := NewPelucio(WithReadWriter(readWriter), WithClock(xtime.DefaultClock))

	externalID := "external_id"

	readWriter.
		On("ReadAccountByExternalID", externalID).
		Return(nil, ErrNotFound)

	readWriter.
		On("WriteAccount", mock.AnythingOfType("*pelucio.Account"), false).
		Return(nil)

	err := pelucio.CreateAccount(context.Background(), externalID, "name", Debit, nil)

	if err != nil {
		t.Fatalf("expected no error. got %v", err)
	}

	if !readWriter.AssertExpectations(t) {
		t.Errorf("mock expected calls not occurred")
	}
}

func TestPelucio_CreateAccount_ExternalIDError(t *testing.T) {
	readWriter := new(ReadWriterMock)
	pelucio := NewPelucio(WithReadWriter(readWriter), WithClock(xtime.DefaultClock))

	externalID := "external_id"

	expectedAccount := &Account{Name: "teste"}
	readWriter.
		On("ReadAccountByExternalID", externalID).
		Return(expectedAccount, nil)

	err := pelucio.CreateAccount(context.Background(), externalID, "name", Debit, nil)

	if err != ErrExternalIDAlreadyInUse {
		t.Fatalf("expected error ErrExternalIDAlreadyInUse. got %v", err)
	}

	if !readWriter.AssertExpectations(t) {
		t.Errorf("mock expected calls not occurred")
	}
}

func TestPelucio_CreateAccount_AnyReadError(t *testing.T) {
	readWriter := new(ReadWriterMock)
	pelucio := NewPelucio(WithReadWriter(readWriter), WithClock(xtime.DefaultClock))

	externalID := "external_id"
	expectedErr := errors.New("any error")
	readWriter.
		On("ReadAccountByExternalID", externalID).
		Return(nil, expectedErr)

	err := pelucio.CreateAccount(context.Background(), externalID, "name", Debit, nil)

	if err != expectedErr {
		t.Fatalf("expected error. got %v", err)
	}

	if !readWriter.AssertExpectations(t) {
		t.Errorf("mock expected calls not occurred")
	}
}

func TestPelucio_CreateAccount_AnyWriteErr(t *testing.T) {
	readWriter := new(ReadWriterMock)
	pelucio := NewPelucio(WithReadWriter(readWriter), WithClock(xtime.DefaultClock))

	externalID := "external_id"

	readWriter.
		On("ReadAccountByExternalID", externalID).
		Return(nil, ErrNotFound)

	expectedErr := errors.New("any error")
	readWriter.
		On("WriteAccount", mock.AnythingOfType("*pelucio.Account"), false).
		Return(expectedErr)

	err := pelucio.CreateAccount(context.Background(), externalID, "name", Debit, nil)

	if err != expectedErr {
		t.Fatalf("expected error. got %v", err)
	}

	if !readWriter.AssertExpectations(t) {
		t.Errorf("mock expected calls not occurred")
	}
}

func TestPelucio_UpdateAccount_Success(t *testing.T) {
	readWriter := new(ReadWriterMock)
	pelucio := NewPelucio(WithReadWriter(readWriter), WithClock(xtime.DefaultClock))

	account := NewAccount("external_id", "name", Debit, nil, xtime.DefaultClock)

	readWriter.
		On("ReadAccount", account.ID).
		Return(account, nil)

	readWriter.
		On("WriteAccount", account, true).
		Return(nil)

	err := pelucio.UpdateAccount(context.Background(), account.ID, "name", nil)

	if err != nil {
		t.Fatalf("expected no error. got %v", err)
	}

	if !readWriter.AssertExpectations(t) {
		t.Errorf("mock expected calls not occurred")
	}
}

func TestPelucio_UpdateAccount_AccountNotFound(t *testing.T) {
	readWriter := new(ReadWriterMock)
	pelucio := NewPelucio(WithReadWriter(readWriter), WithClock(xtime.DefaultClock))

	account := NewAccount("external_id", "name", Debit, nil, xtime.DefaultClock)
	readWriter.
		On("ReadAccount", account.ID).
		Return(nil, ErrNotFound)

	err := pelucio.UpdateAccount(context.Background(), account.ID, "name", nil)

	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound. got %v", err)
	}

	if !readWriter.AssertExpectations(t) {
		t.Errorf("mock expected calls not occurred")
	}
}

func TestPelucio_UpdateAccount_WriteError(t *testing.T) {
	readWriter := new(ReadWriterMock)
	pelucio := NewPelucio(WithReadWriter(readWriter), WithClock(xtime.DefaultClock))

	account := NewAccount("external_id", "name", Debit, nil, xtime.DefaultClock)

	readWriter.
		On("ReadAccount", account.ID).
		Return(account, nil)

	readWriter.
		On("WriteAccount", account, true).
		Return(ErrNotFound)

	err := pelucio.UpdateAccount(context.Background(), account.ID, "name", nil)

	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound. got %v", err)
	}

	if !readWriter.AssertExpectations(t) {
		t.Errorf("mock expected calls not occurred")
	}
}

func TestPelucio_DeleteAccount_Success(t *testing.T) {
	readWriter := new(ReadWriterMock)
	pelucio := NewPelucio(WithReadWriter(readWriter), WithClock(xtime.DefaultClock))

	account := NewAccount("external_id", "name", Debit, nil, xtime.DefaultClock)

	readWriter.
		On("ReadAccount", account.ID).
		Return(account, nil)

	readWriter.
		On("WriteAccount", account, true).
		Return(nil)

	err := pelucio.DeleteAccount(context.Background(), account.ID)

	if err != nil {
		t.Fatalf("expected no error. got %v", err)
	}

	if !readWriter.AssertExpectations(t) {
		t.Errorf("mock expected calls not occurred")
	}
}

func TestPelucio_DeleteAccount_AccountNotFound(t *testing.T) {
	readWriter := new(ReadWriterMock)
	pelucio := NewPelucio(WithReadWriter(readWriter), WithClock(xtime.DefaultClock))

	account := NewAccount("external_id", "name", Debit, nil, xtime.DefaultClock)
	readWriter.
		On("ReadAccount", account.ID).
		Return(nil, ErrNotFound)

	err := pelucio.DeleteAccount(context.Background(), account.ID)

	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound. got %v", err)
	}

	if !readWriter.AssertExpectations(t) {
		t.Errorf("mock expected calls not occurred")
	}
}

func TestPelucio_DeleteAccount_AccountWithBalance(t *testing.T) {
	readWriter := new(ReadWriterMock)
	pelucio := NewPelucio(WithReadWriter(readWriter), WithClock(xtime.DefaultClock))

	account := NewAccount("external_id", "name", Debit, nil, xtime.DefaultClock)
	account.Balance = Balance{
		Currency("BRL"): big.NewInt(1),
	}

	readWriter.
		On("ReadAccount", account.ID).
		Return(account, nil)

	err := pelucio.DeleteAccount(context.Background(), account.ID)

	if err == nil {
		t.Fatalf("expected error. got %v", err)
	}

	if !readWriter.AssertExpectations(t) {
		t.Errorf("mock expected calls not occurred")
	}
}

func TestPelucio_DeleteAccount_WriteError(t *testing.T) {
	readWriter := new(ReadWriterMock)
	pelucio := NewPelucio(WithReadWriter(readWriter), WithClock(xtime.DefaultClock))

	account := NewAccount("external_id", "name", Debit, nil, xtime.DefaultClock)

	readWriter.
		On("ReadAccount", account.ID).
		Return(account, nil)

	readWriter.
		On("WriteAccount", account, true).
		Return(ErrNotFound)

	err := pelucio.DeleteAccount(context.Background(), account.ID)

	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound. got %v", err)
	}

	if !readWriter.AssertExpectations(t) {
		t.Errorf("mock expected calls not occurred")
	}
}

func TestPelucio_FindAccounts(t *testing.T) {
	readWriter := new(ReadWriterMock)
	pelucio := NewPelucio(WithReadWriter(readWriter), WithClock(xtime.DefaultClock))

	query := ReadAccountFilter{
		FromDate: xtime.DefaultClock.NilNow(),
	}
	account := NewAccount("external_id", "name", Debit, nil, xtime.DefaultClock)

	readWriter.
		On("ReadAccounts", query).
		Return([]*Account{account}, nil)

	accounts, err := pelucio.FindAccounts(context.Background(), query)

	if err != nil {
		t.Errorf("expectd error nil. got %v", err)
	}

	if len(accounts) != 1 {
		t.Errorf("expected len of 1. got %v", len(accounts))
	}

	if !readWriter.AssertExpectations(t) {
		t.Errorf("mock expected calls not occurred")
	}
}

func TestPelucio_FindAccountByID(t *testing.T) {
	readWriter := new(ReadWriterMock)
	pelucio := NewPelucio(WithReadWriter(readWriter), WithClock(xtime.DefaultClock))

	account := NewAccount("external_id", "name", Debit, nil, xtime.DefaultClock)

	readWriter.
		On("ReadAccount", account.ID).
		Return(account, nil)

	res, err := pelucio.FindAccountByID(context.Background(), account.ID)

	if err != nil {
		t.Errorf("expectd error nil. got %v", err)
	}

	if res.ID != account.ID {
		t.Errorf("expected account id to be %v. got %v", res.ID, account.ID)
	}

	if !readWriter.AssertExpectations(t) {
		t.Errorf("mock expected calls not occurred")
	}
}

func TestPelucio_FindAccountByExternalID(t *testing.T) {
	readWriter := new(ReadWriterMock)
	pelucio := NewPelucio(WithReadWriter(readWriter), WithClock(xtime.DefaultClock))

	account := NewAccount("external_id", "name", Debit, nil, xtime.DefaultClock)

	readWriter.
		On("ReadAccountByExternalID", account.ExternalID).
		Return(account, nil)

	res, err := pelucio.FindAccountByExternalID(context.Background(), account.ExternalID)

	if err != nil {
		t.Errorf("expectd error nil. got %v", err)
	}

	if res.ID != account.ID {
		t.Errorf("expected account id to be %v. got %v", res.ID, account.ID)
	}

	if !readWriter.AssertExpectations(t) {
		t.Errorf("mock expected calls not occurred")
	}
}

func TestPelucio_BalanceOf(t *testing.T) {
	readWriter := new(ReadWriterMock)
	pelucio := NewPelucio(WithReadWriter(readWriter), WithClock(xtime.DefaultClock))

	account := NewAccount("external_id", "name", Debit, nil, xtime.DefaultClock)
	account.Balance = Balance{
		Currency("BRL"): big.NewInt(1),
	}

	readWriter.
		On("ReadAccount", account.ID).
		Return(account, nil)

	balance, err := pelucio.BalanceOf(context.Background(), account.ID)

	if err != nil {
		t.Fatalf("expected no error. got %v", err)
	}

	if len(balance) != 1 {
		t.Errorf("expected balance to have 1. Got %v", len(balance))
	}

	if balance["BRL"].Cmp(account.Balance["BRL"]) != 0 {
		t.Errorf("expected currency brl to have 1. got %v", balance["BRL"])
	}

	if !readWriter.AssertExpectations(t) {
		t.Errorf("mock expected calls not occurred")
	}
}

func TestPelucio_BalanceOf_AccountNotFound(t *testing.T) {
	readWriter := new(ReadWriterMock)
	pelucio := NewPelucio(WithReadWriter(readWriter), WithClock(xtime.DefaultClock))

	account := NewAccount("external_id", "name", Debit, nil, xtime.DefaultClock)
	readWriter.
		On("ReadAccount", account.ID).
		Return(nil, ErrNotFound)

	_, err := pelucio.BalanceOf(context.Background(), account.ID)

	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound. got %v", err)
	}

	if !readWriter.AssertExpectations(t) {
		t.Errorf("mock expected calls not occurred")
	}
}

func TestPelucio_BalanceOfAccountFromLedger_Success(t *testing.T) {
	readWriter := new(ReadWriterMock)
	pelucio := NewPelucio(WithReadWriter(readWriter), WithClock(xtime.DefaultClock))

	account := NewAccount("external_id", "name", Debit, nil, xtime.DefaultClock)
	account.Balance = Balance{
		Currency("BRL"): big.NewInt(1),
	}

	entries := []*Entry{
		{
			ID:          xuuid.New(),
			AccountID:   account.ID,
			EntrySide:   Debit,
			AccountSide: Debit,
			Amount:      big.NewInt(2),
			Currency:    "BRL",
		},
		{
			ID:          xuuid.New(),
			AccountID:   account.ID,
			EntrySide:   Credit,
			AccountSide: Debit,
			Amount:      big.NewInt(1),
			Currency:    "BRL",
		},
	}

	readWriter.
		On("ReadAccount", account.ID).
		Return(account, nil)

	readWriter.
		On("ReadEntriesOfAccount", account.ID).
		Return(entries, nil)

	balance, err := pelucio.BalanceOfAccountFromLedger(context.Background(), account.ID)

	if err != nil {
		t.Fatalf("expected no error. got %v", err)
	}

	if len(balance) != 1 {
		t.Errorf("expected balance to have 1. Got %v", len(balance))
	}

	if balance["BRL"].Cmp(big.NewInt(1)) != 0 {
		t.Errorf("expected currency brl to have 1. got %v", balance["BRL"])
	}

	if !readWriter.AssertExpectations(t) {
		t.Errorf("mock expected calls not occurred")
	}
}

func TestPelucio_BalanceOfAccountFromLedger_AccountIDNil(t *testing.T) {
	readWriter := new(ReadWriterMock)
	pelucio := NewPelucio(WithReadWriter(readWriter), WithClock(xtime.DefaultClock))

	_, err := pelucio.BalanceOfAccountFromLedger(context.Background(), xuuid.Empty)

	if err == nil {
		t.Fatalf("expected no error. got %v", err)
	}
}

func TestPelucio_BalanceOfAccountFromLedger_AccountNotFound(t *testing.T) {
	readWriter := new(ReadWriterMock)
	pelucio := NewPelucio(WithReadWriter(readWriter), WithClock(xtime.DefaultClock))

	account := NewAccount("external_id", "name", Debit, nil, xtime.DefaultClock)

	readWriter.
		On("ReadAccount", account.ID).
		Return(nil, ErrNotFound)

	_, err := pelucio.BalanceOfAccountFromLedger(context.Background(), account.ID)

	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound. got %v", err)
	}

	if !readWriter.AssertExpectations(t) {
		t.Errorf("mock expected calls not occurred")
	}
}

func TestPelucio_BalanceOfAccountFromLedger_ReadEntriesError(t *testing.T) {
	readWriter := new(ReadWriterMock)
	pelucio := NewPelucio(WithReadWriter(readWriter), WithClock(xtime.DefaultClock))

	account := NewAccount("external_id", "name", Debit, nil, xtime.DefaultClock)
	account.Balance = Balance{
		Currency("BRL"): big.NewInt(1),
	}

	readWriter.
		On("ReadAccount", account.ID).
		Return(account, nil)

	readWriter.
		On("ReadEntriesOfAccount", account.ID).
		Return(nil, ErrNotFound)

	_, err := pelucio.BalanceOfAccountFromLedger(context.Background(), account.ID)

	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound. got %v", err)
	}

	if !readWriter.AssertExpectations(t) {
		t.Errorf("mock expected calls not occurred")
	}
}

func TestPelucio_BalanceOfAccountFromLedger_ComputeEntryError(t *testing.T) {
	readWriter := new(ReadWriterMock)
	pelucio := NewPelucio(WithReadWriter(readWriter), WithClock(xtime.DefaultClock))

	account := NewAccount("external_id", "name", Debit, nil, xtime.DefaultClock)
	account.Balance = Balance{
		Currency("BRL"): big.NewInt(1),
	}

	entries := []*Entry{
		{
			ID:          xuuid.New(),
			AccountID:   account.ID,
			EntrySide:   Debit,
			AccountSide: Credit,
			Amount:      big.NewInt(2),
			Currency:    "BRL",
		},
	}

	readWriter.
		On("ReadAccount", account.ID).
		Return(account, nil)

	readWriter.
		On("ReadEntriesOfAccount", account.ID).
		Return(entries, nil)

	_, err := pelucio.BalanceOfAccountFromLedger(context.Background(), account.ID)

	if err != ErrAccountSideMismatch {
		t.Fatalf("expected ErrAccountSideMismatch. got %v", err)
	}

	if !readWriter.AssertExpectations(t) {
		t.Errorf("mock expected calls not occurred")
	}
}
