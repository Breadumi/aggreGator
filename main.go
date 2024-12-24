package main

import (
	"encoding/json"
	"fmt"

	"github.com/Breadumi/aggreGator/internal/config"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	cfg.SetUser("Bread")
	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonBytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonBytes))

}
