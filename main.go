package main

import (
	"database/sql"
	_ "embed"
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/seagalputra/cashbot/expense"
	"github.com/seagalputra/cashbot/history"
	"github.com/seagalputra/cashbot/tgbot"
	_ "modernc.org/sqlite"
)

var (
	db *sql.DB

	//go:embed assets/help.txt
	helpMsg string
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

	commandHistoryRepo := &history.CommandHistoryRepo{DB: db}
	expenseRepo := &expense.ExpenseRepo{DB: db}

	tgbot.HelpMsg = helpMsg
	telegramBot := &tgbot.TelegramBot{
		CommandHistoryRepo: *commandHistoryRepo,
		ExpenseRepo:        *expenseRepo,
	}

	telegramBot.SetupTelegramBot()
}
