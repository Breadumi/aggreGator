package main

import (
	"fmt"
	"os"

	"github.com/Breadumi/aggreGator/internal/config"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	s := state{
		cfg: &cfg,
	}

	cmds := commands{
		commandMap: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
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

}
