package db

import (
	"context"
	"testing"
	"time"

	"github.com/andrew/bankoferico/util"
	"github.com/stretchr/testify/require"
)

//create a random transaction to be used in actual tests
//because it requires a database connection, we need to set it up in the queries object first (main_test)
func createRandomTransaction(t *testing.T, a Account) Transaction {
	//create mock data that would get sent to create account function
	arg := CreateTransactionParams{
		AccountID: a.ID,
		Amount:    util.RandomAmount(),
	}

	transaction, err := testQueries.CreateTransaction(context.Background(), arg)
	//use testify assertions to have no error and not an empty transaction
	require.NoError(t, err)
	require.NotEmpty(t, transaction)

	//check if created account matches the input data
	require.Equal(t, arg.AccountID, transaction.AccountID)
	require.Equal(t, arg.Amount, transaction.Amount)
	require.NotZero(t, transaction.CreatedAt)
	return transaction
}

func TestCreateTransaction(t *testing.T) {
	account := createRandomAccount(t)
	createRandomTransaction(t, account)
}

func TestGetTransaction(t *testing.T) {
	account := createRandomAccount(t)
	trans1 := createRandomTransaction(t, account)

	trans2, err := testQueries.GetTransaction(context.Background(), trans1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, trans2)

	require.Equal(t, trans1.ID, trans2.ID)
	require.Equal(t, trans1.AccountID, trans2.AccountID)
	require.Equal(t, trans1.Amount, trans2.Amount)
	require.WithinDuration(t, trans1.CreatedAt, trans2.CreatedAt, time.Second)
}

func TestListTransactions(t *testing.T) {
	//create a bunch of transactions
	a := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomTransaction(t, a)
	}

	arg := ListTransactionsParams{
		AccountID: a.ID,
		Limit:     5,
		Offset:    5,
	}

	strans, err := testQueries.ListTransactions(context.Background(), arg)
	require.NoError(t, err)

	//require the returned length to be 5
	require.Len(t, strans, 5)

	for _, trans := range strans {
		require.NotEmpty(t, trans)
	}
}
