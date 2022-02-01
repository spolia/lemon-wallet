package movement

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func TestSaveMovement_ok(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		require.NoError(t, err)
	}
	repository := New(db)
	defer db.Close()

	movement := Movement{
		Type:         DepositMov,
		Amount:       100.2,
		CurrencyName: USDT,
		UserID:       1,
	}
	// When
	query := "INSERT INTO movements_usdt(mov_type,currency_name,tx_amount,user_id)VALUES (?,?,?,?);"
	mock.ExpectExec(query).
		WithArgs(movement.Type, movement.CurrencyName, movement.Amount, movement.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// then
	movementID, err := repository.Save(context.Background(), movement)
	require.NoError(t, err)
	require.Equal(t, int64(1), movementID)
}

func TestSaveMovement_Error(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		require.NoError(t, err)
	}
	repository := New(db)
	defer db.Close()

	movement := Movement{
		Type:         DepositMov,
		Amount:       100.2,
		CurrencyName: USDT,
		UserID:       1,
	}

	// When
	query := "INSERT INTO movements_usdt(mov_type,currency_name,tx_amount,user_id)VALUES (?,?,?,?);"
	mock.ExpectExec(query).
		WithArgs(movement.Type, movement.CurrencyName, movement.Amount, movement.UserID).WillReturnError(&mysql.MySQLError{
		Number: 1264,
	})

	// then
	movementID, err := repository.Save(context.Background(), movement)
	require.Error(t, err)
	require.EqualError(t, ErrorInsufficientBalance, err.Error())
	require.Equal(t, int64(0), movementID)
}

func TestSaveMovement_ErrorWrongCurrency(t *testing.T) {
	// Given
	repository := New(nil)

	// When
	movementID, err := repository.Save(context.Background(), Movement{
		Type:         DepositMov,
		Amount:       100.2,
		CurrencyName: "wrong",
		UserID:       1,
	})

	// Then
	require.Error(t, err)
	require.EqualError(t, ErrorWrongCurrency, err.Error())
	require.Equal(t, int64(0), movementID)
}

func TestInitSave_ok(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		require.NoError(t, err)
	}
	repository := New(db)
	defer db.Close()

	movement := Movement{
		Type:         DepositMov,
		UserID:       1,
	}
	// When
	mock.ExpectBegin()

	mock.ExpectExec("INSERT INTO movements_usdt(mov_type,tx_amount,total_amount,user_id)VALUES (?,?,?,?);").
		WithArgs(movement.Type, movement.Amount, movement.TotalAmount, movement.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO movements_ars(mov_type,tx_amount,total_amount,user_id)VALUES (?,?,?,?);").
		WithArgs(movement.Type, movement.Amount, movement.TotalAmount, movement.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO movements_btc(mov_type,tx_amount,total_amount,user_id)VALUES (?,?,?,?);").
		WithArgs(movement.Type, movement.Amount, movement.TotalAmount, movement.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	// then
	err = repository.InitSave(context.Background(), movement)
	require.NoError(t, err)
}