package command_history

import (
	"database/sql"
	"fmt"
)

type (
	CommandHistory struct {
		ID          int
		CommandName string
		Username    string
		Step        int
	}

	CommandHistoryRepo struct {
		DB *sql.DB
	}
)

func (ch *CommandHistoryRepo) InsertHistory(history CommandHistory) error {
	_, err := ch.DB.Exec(
		"INSERT INTO command_histories (username, step, command_name) VALUES (?, ?, ?)",
		history.Username,
		history.Step,
		history.CommandName,
	)

	if err != nil {
		return fmt.Errorf("[InsertHistory] %v", err)
	}

	return nil
}

func (ch *CommandHistoryRepo) IncrementStepByUsername(username string) error {
	_, err := ch.DB.Exec("UPDATE command_histories SET step = (SELECT step FROM command_histories WHERE username = ?) + 1 WHERE username = ?", username, username)

	if err != nil {
		return fmt.Errorf("[IncrementStepByUsername] %v", err)
	}

	return nil
}

func (ch *CommandHistoryRepo) FindByUsername(username string) (*CommandHistory, error) {
	var history CommandHistory
	row := ch.DB.QueryRow("SELECT id, username, step, command_name FROM command_histories WHERE username = ?", username)

	if err := row.Scan(&history.ID, &history.Username, &history.Step, &history.CommandName); err != nil {
		return nil, fmt.Errorf("[FindByUsername] %v", err)
	}

	return &history, nil
}

func (ch *CommandHistoryRepo) DeleteByUsername(username string) error {
	_, err := ch.DB.Exec("DELETE FROM command_histories WHERE username = ?", username)
	if err != nil {
		return fmt.Errorf("[DeleteByUsername] %v", err)
	}

	return nil
}
