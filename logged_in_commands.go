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

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return errors.New("expected two arguments <feed name> and <url>")
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
	err = handlerFollow(s, cmd, user)
	if err != nil {
		return err
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {

	if len(cmd.args) != 1 {
		return errors.New("expected 1 argument <url>")
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

	row, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("User %s is now following %s\n", row.UserName, row.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	if len(cmd.args) != 0 {
		return errors.New("expected no arguments")
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
