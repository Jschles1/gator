package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/Jschles1/gator/internal/config"
	"github.com/Jschles1/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db                  *database.Queries
	c                   *config.Config
	timeBetweenRequests time.Duration
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("error: Not enough arguments provided")
		os.Exit(1)
	}
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(fmt.Errorf("error reading file: %w", err))
		os.Exit(1)
	}

	dbURL := cfg.DbURL

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(fmt.Errorf("error connecting to database: %w", err))
		os.Exit(1)
	}
	dbQueries := database.New(db)

	timeBetweenRequests, err := time.ParseDuration("15s")
	if err != nil {
		fmt.Println(fmt.Errorf("error parsing time duration: %w", err))
		os.Exit(1)
	}

	appState := &state{
		c:                   cfg,
		db:                  dbQueries,
		timeBetweenRequests: timeBetweenRequests,
	}

	cmds := commands{
		entries: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("feeds", handlerFeeds)
	cmds.register("help", handlerHelp)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	cmd := command{
		name: args[1],
		arguments: func() []string {
			if len(args) > 2 {
				return args[2:]
			}
			return []string{}
		}(),
	}

	err = cmds.run(appState, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
