package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/OrationKnight21/gator/internal/config"
	"github.com/OrationKnight21/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	dbURL := "postgres://postgres@localhost:5432/gator?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	read, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	s := state{db: dbQueries, cfg: &read}
	nmap := commands{make(map[string]func(*state, command) error)}
	nmap.register("login", handlerLogin)
	nmap.register("register", registerLogin)
	nmap.register("reset", handlerReset)
	nmap.register("users", users)
	nmap.register("agg", handlerAgg)
	nmap.register("addfeed", middlewareLoggedIn(addfeed))
	nmap.register("feeds", handlerFeeds)
	nmap.register("follow", middlewareLoggedIn(follow))
	nmap.register("following", middlewareLoggedIn(following))
	nmap.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	nmap.register("browse", middlewareLoggedIn(browse))
	cArgs := os.Args
	if len(cArgs) < 2 {
		log.Fatal("insufficient arguments")
	}
	v := command{cArgs[1], cArgs[2:]}
	err = nmap.run(&s, v)
	if err != nil {
		log.Fatal(err)
	}
}
