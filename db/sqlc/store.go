package db

import (
	"context"
	"database/sql"
	"fmt"
)

//Store provides all functions to execute db queries and transactions
//Extending queries functionality since it can only support one database function at a time (composition)
type Store struct {
	*Queries
	db *sql.DB
}

//NewStore creates a store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// Executes a generic database transaction
// Takes in context and a callback fn
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// Start a new db transaction with context and default transaction options (nil)
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Create a new queries object with that transaction
	q := New(tx)

	// Call the callback fn with the created queries
	err = fn(q)
	// Commit or rollback the transaction based on the error returned
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// Include the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// Result of the transfer transaction
type TransferTxResult struct {
	Transfer        Transfer    `json:"transfer"`
	FromAccount     Account     `json:"from_account"`
	ToAccount       Account     `json:"to_account"`
	FromTransaction Transaction `json:"from_transaction"`
	ToTransaction   Transaction `json:"to_transaction"`
}

// TransferTx performs a money transfer from one account to the other
//It creates a transfer record, adds account entry, and updates account balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	//create an empty result
	var result TransferTxResult

	//call the transaction function and inside the single database functions
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		//create transfer
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		//create from transaction
		result.FromTransaction, err = q.CreateTransaction(ctx, CreateTransactionParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		//create to transaction
		result.ToTransaction, err = q.CreateTransaction(ctx, CreateTransactionParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		//This if statement makes sure the updating happens in the same order with concurrent transactions preventing a deadlock
		if arg.FromAccountID < arg.ToAccountID {
			// update from and to accounts with a helper func
			result.FromAccount, result.ToAccount, err = updateMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = updateMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return nil
	})

	return result, err
}

func updateMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1, account2 Account, err error) {
	account1, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:      accountID1,
		Balance: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:      accountID2,
		Balance: amount2,
	})
	if err != nil {
		return
	}

	return
}
