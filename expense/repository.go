package expense

import (
	"database/sql"
	"fmt"
	"time"
)

type (
	Expense struct {
		ID          int
		Title       string
		Amount      int
		Type        string
		ExpenseDate time.Time
		Username    string
	}

	ExpenseRepo struct {
		DB *sql.DB
	}
)

func (e *ExpenseRepo) Insert(expense Expense) error {
	_, err := e.DB.Exec(
		"INSERT INTO expenses (title, amount, type, expense_date, username) VALUES (?, ?, ?, ?, ?)",
		expense.Title,
		expense.Amount,
		expense.Type,
		expense.ExpenseDate,
		expense.Username,
	)

	if err != nil {
		return fmt.Errorf("[Insert] %v", err)
	}

	return nil
}
