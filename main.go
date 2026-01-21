package main

import (
	"database/sql"
	"log"
	"os"

	cfg "github.com/hreshchyshynt/gator/internal/config"
	"github.com/hreshchyshynt/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	config, err := cfg.Read()
	if err != nil {
		log.Fatalf("error read config: %v", err)
	}

	db, err := sql.Open("postgres", config.DbUrl)
	if err != nil {
		log.Fatalln("Can't open database connection")
	}

	state := NewState(config, database.New(db))
	commands := NewCommands()

	commands.Register("login", handleLogin)
	commands.Register("register", handleRegister)
	commands.Register("reset", handleReset)
	commands.Register("users", handleGetAllUsers)
	commands.Register("agg", handleAggregate)
	commands.Register("addfeed", middlewareLoggedIn(handleAddFeed))
	commands.Register("feeds", handleListFeeds)
	commands.Register("follow", middlewareLoggedIn(handleFeedFollow))
	commands.Register("following", middlewareLoggedIn(handleFollowing))
	commands.Register("unfollow", middlewareLoggedIn(handleUnfollowing))
	commands.Register("browse", middlewareLoggedIn(handleBrowse))
	commands.Register("fetch", handleFetch)

	args := os.Args

	if len(args) < 2 {
		log.Fatalln("No arguments provided")
	}

	command := NewCommand(
		args[1],
		args[2:],
	)

	err = commands.Run(&state, command)
	if err != nil {
		log.Fatal(err)
	}

	config, err = cfg.Read()
	if err != nil {
		log.Fatalf("error read new config: %v", err)
	}
}
