package movement

import (
	"context"
	"errors"
	"time"
)

const (
	DepositMov = "deposit"
	ExtractMov = "extract"
	BTC        = "BTC"
	ARS        = "ARS"
	USDT       = "USDT"
)

var currencyTable = map[string]string{
	USDT: "movements_usdt",
	ARS:  "movements_ars",
	BTC:  "movements_btc",
}

var (
	ErrorInsufficientBalance = errors.New("movement: insufficient balance")
	ErrorWrongUser           = errors.New("movement: wrong user")
	ErrorWrongCurrency       = errors.New("movement: wrong currency")
	ErrorNoMovements         = errors.New("movement: there no movements")
)

type AccountExtract map[string]float64

type Repository interface {
	Save(ctx context.Context, movement Movement) (int64, error)
	ListAll(ctx context.Context, id int64) ([]Movement, error)
	InitInsert(ctx context.Context, movement Movement) error
	GetAccountExtract(ctx context.Context, id int64) (AccountExtract, error)
	Search(ctx context.Context, limit, offset uint64, movType, currencyName string, userID int64) ([]Row, error)
}

type Movement struct {
	ID           int64   `json: "id"`
	Type         string  `json: "type" binding:"required"`
	Amount       float64 `json: "amount" binding:"required"`
	CurrencyName string  `json: "currencyname" binding:"required"`
	UserID       int64   `json: "userid" binding:"required"`
	TotalAmount  float64 `json: "totalamount"`
}

type Currency struct {
	ID     int64
	Name   string
	Digits float64
}

type Row struct {
	CurrencyName string
	Type         string
	DateCreated  time.Time
	Amount       float64
	TotalAmount  float64
}

func getCurrencyTable(currency string) string {
	return currencyTable[currency]
}
