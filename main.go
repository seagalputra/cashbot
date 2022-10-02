package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/seagalputra/cashbot/command_history"
	"github.com/seagalputra/cashbot/expense"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	_ "modernc.org/sqlite"
)

var (
	db *sql.DB

	selector = &telebot.ReplyMarkup{}
	btnYes   = selector.Data("Yes", "yes")
	btnNo    = selector.Data("No", "no")

	userState = map[string]interface{}{}

	telegramApiToken string

	helpMsg = `
I can help you to record your money expenses.

You can control me by sending the following commands:

/addexpense - add new expense
/editexpense - edit an expense
/viewexpense - view an available expense
/configure - config your Google Sheets account
/help - show available commands
`
)

func main() {
	// setup sqlite db
	fn := filepath.Join("db", "development.sqlite")
	db, err := sql.Open("sqlite", fn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	telegramApiToken = os.Getenv("TELEGRAM_BOT_API_KEY")

	commandHistoryRepo := &command_history.CommandHistoryRepo{DB: db}
	expenseRepo := &expense.ExpenseRepo{DB: db}

	pref := telebot.Settings{
		Token: telegramApiToken,
		Poller: &telebot.LongPoller{
			Timeout: 10 * time.Second,
		},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Use(middleware.Logger())

	b.Handle(telebot.OnText, func(c telebot.Context) error {
		// For handling add expense
		username := c.Message().Chat.Username
		history, err := commandHistoryRepo.FindByUsername(username)
		if err != nil {
			log.Println(err)
		}

		step := history.Step
		userKey := username + "_" + history.CommandName
		log.Println(userState)
		log.Println(userKey)
		var msg string
		switch step {
		case 1:
			msg = "How many amount do you want to record?"
			exp := userState[userKey].(*expense.Expense)
			exp.Title = c.Text()
			userState[userKey] = exp
			commandHistoryRepo.IncrementStepByUsername(username)
		case 2:
			msg = "When you do the expense?"
			exp := userState[userKey].(*expense.Expense)
			amount, err := strconv.Atoi(c.Text())
			if err != nil {
				log.Println(err)
			}
			exp.Amount = amount
			userState[userKey] = exp
			commandHistoryRepo.IncrementStepByUsername(username)
		case 3:
			msg = "Here's your expense. Please re-check if it wrongs"
			exp := userState[userKey].(*expense.Expense)
			// TODO: change based on user input
			exp.ExpenseDate = time.Now()
			userState[userKey] = exp
			commandHistoryRepo.IncrementStepByUsername(username)

			selector.Inline(
				selector.Row(btnYes, btnNo),
			)
		}

		return c.Send(msg, selector)
	})

	b.Handle("/start", func(c telebot.Context) error {
		username := c.Chat().Username
		err := commandHistoryRepo.DeleteByUsername(username)
		if err != nil {
			log.Println(err)
		}

		return c.Send("Welcome to Cash Bot")
	})

	b.Handle("/addexpense", func(c telebot.Context) error {
		msg := "Please provide the title of the expense"
		username := c.Message().Chat.Username

		history := command_history.CommandHistory{}
		history.Username = username
		history.CommandName = "addexpense"
		history.Step = 1

		commandHistoryRepo.InsertHistory(history)

		userKey := username + "_" + "addexpense"
		userState[userKey] = &expense.Expense{
			Username: username,
		}

		return c.Send(msg)
	})

	b.Handle("/help", func(c telebot.Context) error {
		return c.Send(helpMsg)
	})

	b.Handle("/cancel", func(c telebot.Context) error {
		return fmt.Errorf("Not yet implemented!")
	})

	b.Handle(&btnYes, func(c telebot.Context) error {
		username := c.Message().Chat.Username
		history, err := commandHistoryRepo.FindByUsername(username)
		if err != nil {
			log.Println(err)
		}

		userKey := username + "_" + history.CommandName
		exp := userState[userKey].(*expense.Expense)
		err = expenseRepo.Insert(*exp)
		if err != nil {
			log.Println(err)
		}

		commandHistoryRepo.DeleteByUsername(username)
		return c.Send("Your expense was successfully added!")
	})

	b.Handle(&btnNo, func(c telebot.Context) error {
		return c.Respond()
	})

	b.Start()
}
