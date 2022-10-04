package tgbot

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/seagalputra/cashbot/expense"
	"github.com/seagalputra/cashbot/history"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type TelegramBot struct {
	CommandHistoryRepo history.CommandHistoryRepo
	ExpenseRepo        expense.ExpenseRepo
}

var (
	selector = &telebot.ReplyMarkup{}
	btnYes   = selector.Data("Yes", "yes")
	btnNo    = selector.Data("No", "no")

	expenseTypeSelector = &telebot.ReplyMarkup{}
	btnIncome           = expenseTypeSelector.Data("Income", "income", "income")
	btnOutcome          = expenseTypeSelector.Data("Outcome", "outcome", "outcome")

	userState = map[string]interface{}{}

	telegramApiToken string

	HelpMsg string
)

func (t *TelegramBot) HandleCommandStart(c telebot.Context) error {
	username := c.Chat().Username
	err := t.CommandHistoryRepo.DeleteByUsername(username)
	if err != nil {
		log.Println(err)
	}

	return c.Send("Welcome to Cash Bot")
}

func (t *TelegramBot) HandleText(c telebot.Context) error {
	// For handling add expense
	username := c.Message().Chat.Username
	history, err := t.CommandHistoryRepo.FindByUsername(username)
	if err != nil {
		log.Println(err)
	}

	step := history.Step
	userKey := username + "_" + history.CommandName
	var msg string
	switch step {
	case 1:
		msg = "How many amount do you want to record?"
		exp := userState[userKey].(*expense.Expense)
		exp.Title = c.Text()
		userState[userKey] = exp
		t.CommandHistoryRepo.IncrementStepByUsername(username)

		return c.Send(msg)
	case 2:
		msg = "When you do the expense?"
		exp := userState[userKey].(*expense.Expense)
		amount, err := strconv.Atoi(c.Text())
		if err != nil {
			log.Println(err)
		}
		exp.Amount = amount
		userState[userKey] = exp
		t.CommandHistoryRepo.IncrementStepByUsername(username)

		return c.Send(msg)
	case 3:
		msg = "Here's your expense. Please re-check if it wrongs"
		exp := userState[userKey].(*expense.Expense)
		expenseDate, err := expense.ParseExpenseDate(c.Text())
		if err != nil {
			msg = "Oops, i'm failed to recognized your time format, use day-month-year (DD-MM-YYYY) pattern."
			log.Println(err)

			return c.Send(msg)
		}

		exp.ExpenseDate = *expenseDate
		userState[userKey] = exp
		t.CommandHistoryRepo.IncrementStepByUsername(username)

		selector.Inline(
			selector.Row(btnYes, btnNo),
		)

		return c.Send(msg, selector)
	}

	return nil
}

func (t *TelegramBot) HandleAddExpense(c telebot.Context) error {
	msg := "What type of your expense is?"
	username := c.Message().Chat.Username

	h := history.CommandHistory{}
	h.Username = username
	h.CommandName = "addexpense"

	t.CommandHistoryRepo.InsertHistory(h)

	userKey := username + "_" + "addexpense"
	userState[userKey] = &expense.Expense{
		Username: username,
	}

	expenseTypeSelector.Inline(
		expenseTypeSelector.Row(btnIncome, btnOutcome),
	)

	return c.Send(msg, expenseTypeSelector)
}

func (t *TelegramBot) HandleHelp(c telebot.Context) error {
	return c.Send(HelpMsg)
}

func (t *TelegramBot) HandleCancel(c telebot.Context) error {
	msg := "The current operation has been cancelled. Please let me know when you have any request to me"
	username := c.Message().Chat.Username
	t.CommandHistoryRepo.DeleteByUsername(username)

	return c.Send(msg)
}

func (t *TelegramBot) HandleBtnYes(c telebot.Context) error {
	username := c.Message().Chat.Username
	history, err := t.CommandHistoryRepo.FindByUsername(username)
	if err != nil {
		log.Println(err)
	}

	userKey := username + "_" + history.CommandName
	exp := userState[userKey].(*expense.Expense)
	err = t.ExpenseRepo.Insert(*exp)
	if err != nil {
		log.Println(err)
	}

	t.CommandHistoryRepo.DeleteByUsername(username)
	return c.Send("Your expense was successfully added!")
}

func (t *TelegramBot) HandleBtnNo(c telebot.Context) error {
	return c.Respond()
}

func (t *TelegramBot) HandleBtnIncome(c telebot.Context) error {
	username := c.Message().Chat.Username
	history, err := t.CommandHistoryRepo.FindByUsername(username)
	if err != nil {
		log.Println(err)
	}

	userKey := username + "_" + history.CommandName
	exp := userState[userKey].(*expense.Expense)
	exp.Type = c.Data()
	userState[userKey] = exp

	t.CommandHistoryRepo.IncrementStepByUsername(username)

	msg := "Please provide the title of the expense"

	return c.Send(msg)
}

func (t *TelegramBot) HandleBtnOutcome(c telebot.Context) error {
	return c.Respond()
}

func (t *TelegramBot) SetupTelegramBot() {
	telegramApiToken = os.Getenv("TELEGRAM_BOT_API_KEY")

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

	b.Handle("/addexpense", t.HandleAddExpense)
	b.Handle("/cancel", t.HandleCancel)
	b.Handle("/help", t.HandleHelp)
	b.Handle("/start", t.HandleCommandStart)
	b.Handle(&btnIncome, t.HandleBtnIncome)
	b.Handle(&btnNo, t.HandleBtnNo)
	b.Handle(&btnOutcome, t.HandleBtnOutcome)
	b.Handle(&btnYes, t.HandleBtnYes)
	b.Handle(telebot.OnText, t.HandleText)

	b.Start()
}
