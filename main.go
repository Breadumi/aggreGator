package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Breadumi/aggreGator/internal/config"
	"github.com/Breadumi/aggreGator/internal/database"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = goose.Up(db, "sql/schema")
	if err != nil {
		fmt.Println(err)
		return
	}
	dbQueries := database.New(db)

	s := state{
		db:  dbQueries,
		cfg: cfg,
	}

	cmds := commands{
		commandMap: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmd := command{}

	if len(os.Args) < 2 {
		fmt.Println("expected function name")
		os.Exit(1)
	}

	if len(os.Args) >= 2 {
		cmd.name = os.Args[1]
	}
	if len(os.Args) >= 3 {
		cmd.args = os.Args[2:]
	}

	err = cmds.run(&s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	os.Exit(0)

}
