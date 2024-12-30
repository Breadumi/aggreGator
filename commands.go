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
		return err
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

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	rss, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}

	fmt.Printf("RSS data:\n %+v\n", *rss)

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return errors.New("expected two arguments <feed name> and <url>")
	}

	user, err := s.db.GetUserByName(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	feedName := cmd.args[0]
	url := cmd.args[1]

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       url,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}

	fmt.Printf("Feed %s has been created:\n", feed.Name)
	feedJSON, err := json.MarshalIndent(feed, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("User: %s\n", feedJSON)

	// Add a feed follow line for current user to this url
	cmd.args = cmd.args[1:]
	cmd.name = "follow"
	err = handlerFollow(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {

	if cmd.args != nil || len(cmd.args) > 0 {
		return errors.New("too many arguments: expected none")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		userID, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("Name: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("User: %s\n", userID.Name)
	}

	return nil
}

func handlerFollow(s *state, cmd command) error {

	if len(cmd.args) != 1 {
		return errors.New("expected 1 argument <url>")
	}

	user, err := s.db.GetUserByName(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	rows, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}
	for i, row := range rows {
		rowJSON, err := json.MarshalIndent(row, "", "  ")
		if err != nil {
			return err
		}
		fmt.Printf("Row %v: %s\n", i, rowJSON)
	}

	return nil
}

func handlerFollowing(s *state, cmd command) error {

	if len(cmd.args) != 0 {
		return errors.New("expected no arguments")
	}

	user, err := s.db.GetUserByName(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	userFeedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, row := range userFeedFollows {
		fmt.Printf("%s\n", row.FeedName)
	}

	return nil
}
