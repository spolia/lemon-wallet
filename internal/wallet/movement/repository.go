package movement

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) *repository {
	return &repository{db: db}
}

// Save inserts a new movement in the database
func (r repository) Save(ctx context.Context, movement Movement) (int64, error) {
	var table string
	if table = getCurrencyTable(movement.CurrencyName); table == "" {
		return 0, ErrorWrongCurrency
	}

	query := fmt.Sprintf("INSERT INTO %s(mov_type,currency_name,tx_amount,user_id)VALUES (?,?,?,?);", table)
	result, err := r.db.ExecContext(ctx, query, movement.Type, movement.CurrencyName, movement.Amount, movement.UserID)
	if err != nil {
		// when tx_amount - total_amount is less than 0
		if err.(*mysql.MySQLError).Number == 1264 {
			return 0, ErrorInsufficientBalance
		}
		// wrong type
		if err.(*mysql.MySQLError).Number == 1265 {
			return 0, ErrorWrongOperation
		}

		if err.(*mysql.MySQLError).Number == 1048 {
			return 0, ErrorWrongUser
		}

		return 0, err
	}

	movID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return movID, nil
}

// InitSave saves initials movements for a new user
func (r repository) InitSave(ctx context.Context, movement Movement) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, v := range movementTables {
		query := fmt.Sprintf("INSERT INTO %s(mov_type,tx_amount,total_amount,user_id)VALUES (?,?,?,?);", v)

		if _, err = tx.ExecContext(ctx, query, movement.Type, movement.Amount, movement.TotalAmount, movement.UserID);
			err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetAccountExtract given an id returns the last movements for each currency
func (r repository) GetAccountExtract(ctx context.Context, id int64) (AccountExtract, error) {
	var accountExtract = make(AccountExtract, 0)
	for k, v := range movementTables {
		var queryResult struct {
			totalAmount float64
		}

		row := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT total_amount FROM %s WHERE date_created = (SELECT MAX(date_created) "+
			"FROM %s WHERE user_id = ?)", v, v), id)
		if err := row.Scan(&queryResult.totalAmount); err != nil {
			return AccountExtract{}, err
		}

		accountExtract[k] = queryResult.totalAmount
	}

	return accountExtract, nil
}

// Search searches the movements for an user applying different filters
func (r repository) Search(ctx context.Context, userID int64, limit, offset uint64, movType, currencyName string) ([]Row, error) {
	var tables = getCurrenciesTables(currencyName)
	var movements []Row
	for _, v := range tables {
		sqlQuery := fmt.Sprintf("SELECT mov_type, currency_name, date_created, tx_amount, total_amount "+
			"FROM %s WHERE user_id = ?", v)
		if movType != "" {
			sqlQuery = fmt.Sprintf("%s AND mov_type = '%s'", sqlQuery, movType)
		}

		if limit > 0 {
			sqlQuery = fmt.Sprintf("%s LIMIT %v OFFSET %v;", sqlQuery, limit, offset)
		}

		rows, err := r.db.QueryContext(ctx, sqlQuery, userID)
		if err != nil || rows.Err()!= nil{
			return []Row{}, err
		}

		for rows.Next() {
			var result Row
			err = rows.Scan(&result.Type, &result.CurrencyName, &result.DateCreated, &result.Amount, &result.TotalAmount)
			if err != nil {
				return []Row{}, err
			}
			movements = append(movements, result)
		}
	}

	if len(movements) == 0 {
		return movements, ErrorNoMovements
	}

	return movements, nil
}
