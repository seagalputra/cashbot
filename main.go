package main

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/seagalputra/cashbot/command_history"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	_ "modernc.org/sqlite"
)

var (
	db *sql.DB

	selector = &telebot.ReplyMarkup{}
	btnYes   = selector.Data("Yes", "yes")
	btnNo    = selector.Data("No", "no")

	// TODO: save to OS env
	API_TOKEN string = "1706311893:AAHogb1MjYlTL1bK6un-tY5pMhhnYx0_K7I"

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

	commandHistoryRepo := &command_history.CommandHistoryRepo{DB: db}

	pref := telebot.Settings{
		Token: API_TOKEN,
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
		var msg string

		switch step {
		case 1:
			msg = "How many amount do you want to record?"
			commandHistoryRepo.IncrementStepByUsername(username)
		case 2:
			msg = "When you do the expense?"
			commandHistoryRepo.IncrementStepByUsername(username)
		case 3:
			msg = "Here's your expense. Please re-check if it wrongs"
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

		history := command_history.CommandHistory{}
		history.Username = c.Message().Chat.Username
		history.CommandName = "addexpense"
		history.Step = 1

		commandHistoryRepo.InsertHistory(history)

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
		commandHistoryRepo.DeleteByUsername(username)

		// TODO: add expense to database
		return c.Send("Your expense was successfully added!")
	})

	b.Handle(&btnNo, func(c telebot.Context) error {
		return c.Respond()
	})

	b.Start()
}
