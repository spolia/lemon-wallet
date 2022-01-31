package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r repository) Save(ctx context.Context, firstName, lastName, alias, email string) (int64, error) {
	result, err := r.db.Exec("INSERT INTO users(first_name,last_name,alias,email)VALUES (?,?,?,?);",
		firstName, lastName, alias, email)
	if err != nil {
		if err.(*mysql.MySQLError).Number==1062{
			return 0 , ErrorAlreadyExist
		}
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r repository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.Exec("DELETE FROM users Where id = ?;", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected < 1 {
		return errors.New("any row affected")
	}

	return nil
}

func (r repository) Get(ctx context.Context, id int64) (User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM users Where id = ?;",id)
	if err != nil {
		return User{}, err
	}

	var user User
	if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Alias, &user.Email); err != nil {
		if err == sql.ErrNoRows {
			return User{}, ErrorUserNotFound
		}
		return User{}, err
	}

	return user, nil
}
