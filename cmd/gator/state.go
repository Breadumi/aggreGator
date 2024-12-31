package main

import (
	"github.com/Breadumi/aggreGator/internal/config"
	"github.com/Breadumi/aggreGator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}
