package entity

import "time"

type Operation struct {
	Id              int       `db:"id"`
	UserId          int       `db:"user_id"`
	AccountId       int       `db:"account_id"`
	Amount          int       `db:"amount"`
	OperationTypeId int       `db:"operation_type_id"`
	Description     string    `db:"description"`
	CreatedAt       time.Time `db:"created_at"`
}

type OperationType struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}
