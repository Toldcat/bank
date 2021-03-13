package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	//make a new store because TransferTx function is a method on it
	store := NewStore(testDBConn)
	//Create random accounts
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	//Create arguments to pass to the test fn
	arg := TransferTxParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        10,
	}

	result, err := store.TransferTx(context.Background(), arg)
	//check transfers
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, result.Transfer.FromAccountID, acc1.ID)
	require.Equal(t, result.Transfer.ToAccountID, acc2.ID)
	require.Equal(t, result.Transfer.Amount, arg.Amount)
	require.NotZero(t, result.Transfer.ID)
	require.NotZero(t, result.Transfer.CreatedAt)

	_, err = store.GetTransfer(context.Background(), result.Transfer.ID)
	require.NoError(t, err)

	//check from transactions
	ft := result.FromTransaction
	require.NotEmpty(t, ft)
	require.Equal(t, ft.AccountID, acc1.ID)
	require.Equal(t, ft.Amount, -arg.Amount)
	require.NotZero(t, ft.ID)
	require.NotZero(t, ft.CreatedAt)
	_, err = store.GetTransaction(context.Background(), ft.ID)
	require.NoError(t, err)

	//check to transactions
	tt := result.ToTransaction
	require.NotEmpty(t, tt)
	require.Equal(t, tt.AccountID, acc2.ID)
	require.Equal(t, tt.Amount, arg.Amount)
	require.NotZero(t, tt.ID)
	require.NotZero(t, tt.CreatedAt)
	_, err = store.GetTransaction(context.Background(), tt.ID)
	require.NoError(t, err)

	//check accounts owners
	fa := result.FromAccount
	require.NotEmpty(t, fa)
	require.Equal(t, fa.ID, acc1.ID)

	ta := result.ToAccount
	require.NotEmpty(t, ta)
	require.Equal(t, ta.ID, acc2.ID)

	//check accounts balance
	diff1 := acc1.Balance - fa.Balance
	require.Equal(t, diff1, arg.Amount)

	diff2 := ta.Balance - acc2.Balance
	require.Equal(t, diff2, arg.Amount)
	require.True(t, diff1 > 0)
	require.True(t, diff1 == diff2)

}
