package db

import (
	"context"
	"testing"
	"time"

	"github.com/andrew/bankoferico/util"
	"github.com/stretchr/testify/require"
)

//create a random account to be used in actual tests
//because it requires a database connection, we need to set it up in the queries object first (main_test)
func createRandomAccount(t *testing.T) Account {
	//create mock data that would get sent to create account function
	arg := CreateAccountParams{
		Owner:    util.RandomName(),
		Balance:  util.RandomAmount(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	//use testify assertions to have no error and not an empty account
	require.NoError(t, err)
	require.NotEmpty(t, account)

	//check if created account matches the input data
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

//Test the Create Account function
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

//Test the Get Account function
func TestGetAccount(t *testing.T) {
	//create random account first
	account1 := createRandomAccount(t)
	//get that account with an ID
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	//check for errors
	require.NoError(t, err)
	//fetched account fields not empty
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

//test the Update Account function
func TestUpdateAccount(t *testing.T) {
	//create an account
	account1 := createRandomAccount(t)

	//update our mock account
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: 69,
	}
	_, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	//get that account from db
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	//compare the two
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account2.Balance, int64(69))
}

//Test the Delete Account function
func TestDeleteAccount(t *testing.T) {
	//create an account
	account1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.Empty(t, account2)
}

//Test GetAllAccounts function
func TestListAccounts(t *testing.T) {
	//create 10 accounts
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	//require the returned length to be 5
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
