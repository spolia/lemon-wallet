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

func (r repository) Save(ctx context.Context, movement Movement) (int64, error) {
	var table string
	fmt.Print(movement.CurrencyName)
	if table = getCurrencyTable(movement.CurrencyName); table == "" {
		return 0, ErrorWrongCurrency
	}

	query := fmt.Sprintf("INSERT INTO %s(mov_type,currency_name,tx_amount,user_id)VALUES (?,?,?,?);", table)
	result, err := r.db.Exec(query, movement.Type, movement.CurrencyName, movement.Amount, movement.UserID)
	if err != nil {
		if err.(*mysql.MySQLError).Number == 1264 {
			return 0, ErrorInsufficientBalance
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

// this insert is just for first insert
func (r repository) InitInsert(ctx context.Context, movement Movement) error {
	// todo: solve what happen if some error occurs, roll back the rows created
	for _, v := range currencyTable {
		query := fmt.Sprintf("INSERT INTO %s(mov_type,tx_amount,total_amount,user_id)VALUES (?,?,?,?);", v)
		_, err := r.db.Exec(query, movement.Type, movement.Amount, movement.TotalAmount, movement.UserID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r repository) ListAll(ctx context.Context, id int64) ([]Movement, error) {
	var movements []Movement
	rows, err := r.db.Query("SELECT * FROM movements where id = ?;", id)
	if err != nil {
		return []Movement{}, err
	}

	err = rows.Scan(&movements)
	if err != nil {
		return []Movement{}, err
	}

	return movements, nil
}

func (r repository) GetAccountExtract(ctx context.Context, id int64) (AccountExtract, error) {
	var accountExtract = make(AccountExtract, 0)
	for k, v := range currencyTable {
		var queryResult struct {
			totalAmount float64
		}

		row := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT total_amount FROM %s WHERE date_created = (SELECT MAX(date_created)"+
			"FROM %s WHERE user_id = ?)", v, v), id)
		if err := row.Scan(&queryResult.totalAmount); err != nil {
			return AccountExtract{}, err
		}

		accountExtract[k] = queryResult.totalAmount
	}

	return accountExtract, nil
}

func (r repository) Search(ctx context.Context, limit, offset uint64, movType, currencyName string, userID int64) ([]Row, error) {
	var currencyTables = make([]string, 0)
	if currencyName == "" {
		for _, v := range currencyTable {
			currencyTables = append(currencyTables, v)
		}
	} else {
		currencyTables = append(currencyTables, currencyTable[currencyName])
	}

	var movements []Row
	for _, v := range currencyTables {

		sqlQuery := fmt.Sprintf("SELECT mov_type, currency_name, date_created, tx_amount, total_amount "+
			"FROM %s WHERE user_id = ?", v)
		if movType != "" {
			sqlQuery = fmt.Sprintf("%s AND mov_type = '%s'", sqlQuery, movType)
		}

		if limit > 0 {
			sqlQuery = fmt.Sprintf("%s LIMIT %v OFFSET %v;", sqlQuery, limit, offset)
		}
		fmt.Print(sqlQuery)
		rows, err := r.db.QueryContext(ctx, sqlQuery, userID)
		if err != nil {
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
