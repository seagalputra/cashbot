package main

import (
	"database/sql"
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/seagalputra/cashbot/command_history"
	"github.com/seagalputra/cashbot/expense"
	"github.com/seagalputra/cashbot/handler"
	_ "modernc.org/sqlite"
)

var (
	db *sql.DB
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// setup sqlite db
	fn := filepath.Join("db", "development.sqlite")
	db, err := sql.Open("sqlite", fn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	LoadEnv()

	commandHistoryRepo := &command_history.CommandHistoryRepo{DB: db}
	expenseRepo := &expense.ExpenseRepo{DB: db}

	telegramBot := &handler.TelegramBot{
		CommandHistoryRepo: *commandHistoryRepo,
		ExpenseRepo:        *expenseRepo,
	}

	telegramBot.SetupTelegramBot()
}
