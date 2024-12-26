package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Breadumi/aggreGator/internal/database"
	"github.com/google/uuid"
)

type command struct {
	name string
	args []string
}

type commands struct {
	commandMap map[string]func(*state, command) error
}

func (cmds *commands) register(name string, f func(*state, command) error) {
	cmds.commandMap[name] = f
}

func (cmds *commands) run(s *state, cmd command) error {
	if f, ok := cmds.commandMap[cmd.name]; ok {
		err := f(s, cmd)
		if err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("no function %s exists", cmd.name)
	}
}

func handlerLogin(s *state, cmd command) error {
	if cmd.args == nil || len(cmd.args) != 1 {
		return errors.New("username required")
	}

	_, err := s.db.GetUserByName(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("user %s does not exist", cmd.args[0])
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User has been set to %s\n", cmd.args[0])

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if cmd.args == nil || len(cmd.args) != 1 {
		return errors.New("username required")
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}

	s.cfg.SetUser(user.Name)
	fmt.Printf("User %s has been created:\n", user.Name)
	userJSON, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("User: %s\n", userJSON)

	return nil
}

func handlerReset(s *state, cmd command) error {

	if cmd.args != nil || len(cmd.args) > 0 {
		return errors.New("too many arguments: expected none")
	}

	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return errors.New("error deleting all users")
	}
	fmt.Println("All users successfully deleted!")
	return nil
}

func handlerUsers(s *state, cmd command) error {

	if cmd.args != nil || len(cmd.args) > 0 {
		return errors.New("too many arguments: expected none")
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}
