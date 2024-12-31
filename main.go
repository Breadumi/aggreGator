package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"

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

	// migrate databases, discard goose migration output but capture errors for later logging
	old := os.Stdout
	r, w, _ := os.Pipe()
	var buf bytes.Buffer
	log.SetOutput(&buf)

	err = goose.Up(db, "sql/schema") // successful output is not printed to console

	// close new output containing unwanted logs
	w.Close()
	r.Close()
	os.Stdout = old

	// log migration errors in console if necessary
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
		commandMap:     make(map[string]func(*state, command) error),
		commandDescMap: make(map[string]string),
	}

	cmds.register("login", handlerLogin,
		"gator login [user]\n\tLogin to a user named [user].")
	cmds.register("register", handlerRegister,
		"gator register [user]\n\tRegister a user named [user].")
	cmds.register("reset", handlerReset,
		"gator reset\n\tDelete all info stored in database.")
	cmds.register("users", handlerUsers,
		"gator users\n\tList all registered users.")
	cmds.register("agg", handlerAgg,
		"gator agg\n\tContinuously save new posts from the current user's followed feeds.")
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed),
		"gator addfeed [url]\n\tAdd a new feed at [url] and sets current user to follow [url].")
	cmds.register("feeds", handlerFeeds,
		"gator feeds\n\tList all tracked feeds.")
	cmds.register("follow", middlewareLoggedIn(handlerFollow),
		"gator follow [url]\n\tHave current user follow the feed at [url].")
	cmds.register("following", middlewareLoggedIn(handlerFollowing),
		"gator following\n\tList all feeds followed by current user.")
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow),
		"gator unfollow [url]\n\tUnfollow feed at [url] for current user.")
	cmds.register("browse", middlewareLoggedIn(handlerBrowse),
		"gator browse [limit]\n\tBrowse posts from current user's followed feeds. Default [limit] is 2.")
	cmd := command{}

	if len(os.Args) < 2 {
		fmt.Println("expected function name - type <gator help> for more info")
		os.Exit(1)
	}

	if len(os.Args) >= 1 && os.Args[1] == "help" {
		cmds.printHelp()
		return
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

func (c *commands) printHelp() {
	fmt.Printf("Commands:\n")
	keys := make([]string, 0, len(c.commandDescMap))

	for k := range c.commandDescMap {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, key := range keys {
		fmt.Printf("%s\n", c.commandDescMap[key])
	}
}
