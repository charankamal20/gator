package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/charankamal20/gator/internal/config"
	"github.com/charankamal20/gator/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type state struct {
	db      *sql.DB
	queries *database.Queries
	conf    *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	dict map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if handler, exists := c.dict[cmd.name]; exists {
		return handler(s, cmd)
	}

	return fmt.Errorf("unknown command: %s", cmd.name)
}

func (c *commands) register(name string, handler func(*state, command) error) {
	if c.dict == nil {
		c.dict = make(map[string]func(*state, command) error)
	}

	c.dict[name] = handler
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("the login handler expects a single argument, the username.")
	}

	user, err := s.queries.GetUser(context.Background(), sql.NullString{ String: cmd.args[0], Valid: true })
	if err != nil {
		fmt.Println("user does not exist")
		return err
	}


	s.conf.SetUser(user.Name.String)
	fmt.Printf("User set to %s\n", user.Name.String)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("the login handler expects a single argument, the username.")
	}

	newuser := &database.CreateUserParams{
		ID: sql.NullString{
			Valid:  true,
			String: uuid.New().String(),
		},
		Name: sql.NullString{
			Valid:  true,
			String: cmd.args[0],
		},
		CreatedAt: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
		UpdatedAt: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}

	_, err := s.queries.GetUser(context.Background(), newuser.Name)
	if err == nil {
		fmt.Println("user already exists")
		return fmt.Errorf("user %s already exists", newuser.Name.String)
	}

	_, err = s.queries.CreateUser(
		context.Background(),
		*newuser,
	)

	if err != nil {
		fmt.Println("could not create user: ", err.Error())
		return err
	}

	s.conf.SetUser(newuser.Name.String)

	fmt.Println("User was created successfully.")
	fmt.Println("ID: ", newuser.ID.String)
	fmt.Println("Name: ", newuser.Name.String)
	fmt.Println("Created at: ", newuser.CreatedAt.Time.String())
	fmt.Println("Updated at: ", newuser.UpdatedAt.Time.String())

	return nil
}

func main() {
	conf := config.Read()
	db, err := sql.Open("postgres", conf.DBUrl)
	if err != nil {
		fmt.Println("could not connect to db: ", err.Error())
		os.Exit(1)
	}

	queries := database.New(db)

	currState := &state{
		conf:    &conf,
		db:      db,
		queries: queries,
	}

	allCommands := &commands{
		dict: make(map[string]func(*state, command) error, 0),
	}

	allCommands.register("login", handlerLogin)
	allCommands.register("register", handlerRegister)

	args := os.Args[1:]
	if len(args) < 2 {
		os.Exit(1)
	}

	cmd := &command{
		name: args[0],
		args: args[1:],
	}

	if err = allCommands.run(currState, *cmd); err != nil {
		os.Exit(1)
	}
}
