package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/charankamal20/gator/internal/config"
	"github.com/charankamal20/gator/internal/database"
	"github.com/charankamal20/gator/internal/pkg/rss"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type State struct {
	db      *sql.DB
	queries *database.Queries
	conf    *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	dict map[string]func(*State, command) error
}

func (c *commands) run(s *State, cmd command) error {
	if handler, exists := c.dict[cmd.name]; exists {
		return handler(s, cmd)
	}

	return fmt.Errorf("unknown command: %s", cmd.name)
}

func (c *commands) register(name string, handler func(*State, command) error) {
	if c.dict == nil {
		c.dict = make(map[string]func(*State, command) error)
	}

	c.dict[name] = handler
}

func handlerLogin(s *State, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("the login handler expects a single argument, the username.")
	}

	user, err := s.queries.GetUser(context.Background(), sql.NullString{String: cmd.args[0], Valid: true})
	if err != nil {
		fmt.Println("user does not exist")
		return err
	}

	s.conf.SetUser(user.Name.String)
	fmt.Printf("User set to %s\n", user.Name.String)

	return nil
}

func handlerRegister(s *State, cmd command) error {
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

func resetHandler(s *State, cmd command) error {
	err := s.queries.DeleteAllUsers(context.Background())
	if err != nil {
		fmt.Println("could not reset users: ", err.Error())
		return err
	}

	fmt.Println("All users have been deleted successfully.")
	return nil
}

func handleUsers(s *State, cmd command) error {
	users, err := s.queries.GetAllUsers(context.Background())
	if err != nil {
		fmt.Println("could not fetch users: ", err.Error())
		return err
	}

	var currUser = s.conf.CurrentUsername
	for _, user := range users {
		if (user.Name.String == currUser) {
			fmt.Printf(" * %s (current)\n", user.Name.String)
			continue
		}

		fmt.Printf(" * %s\n", user.Name.String)
	}

	return nil
}

func handleAgg(s *State, cmd command) error {

	data, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Title: ", data.Channel.Title)
	fmt.Println("Description: ", data.Channel.Description)
	for _, item := range data.Channel.Item {
		fmt.Println("Item Title: ", item.Title)
		fmt.Println("Item Link: ", item.Link)
		fmt.Println("Item Description: ", item.Description)
		fmt.Println("Item PubDate: ", item.PubDate)
		fmt.Println("--------------------------------------------------")
	}

	return nil
}

func handleAddFeed(s *State, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("the addfeed handler expects two arguments, the feed URL and name.")
	}

	user, err := s.queries.GetUser(context.Background(), sql.NullString{String: s.conf.CurrentUsername, Valid: true})
	if err != nil {
		return err
	}

	feedURL := cmd.args[1]
	feedName := cmd.args[0]

	feed := &database.CreateFeedParams{
		ID:     uuid.New().String(),
		Name:   feedName,
		Url:    feedURL,
		UserID: user.ID.String,
	}

	newFeed, err := s.queries.CreateFeed(context.Background(), *feed)
	if err != nil {
		fmt.Println("could not add feed: ", err.Error())
		return err
	}

	fmt.Println("Feed added successfully:")
	fmt.Printf("ID: %s\n", newFeed.ID)
	fmt.Printf("Name: %s\n", newFeed.Name)
	fmt.Printf("URL: %s\n", newFeed.Url)
	fmt.Println("User ID: ", newFeed.UserID)
	fmt.Println("CreatedAt: ", newFeed.CreatedAt.String())
	fmt.Println("UpdatedAt: ", newFeed.UpdatedAt.String())

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

	currState := &State{
		conf:    &conf,
		db:      db,
		queries: queries,
	}

	allCommands := &commands{
		dict: make(map[string]func(*State, command) error, 0),
	}

	allCommands.register("login", handlerLogin)
	allCommands.register("register", handlerRegister)
	allCommands.register("reset", resetHandler)
	allCommands.register("users", handleUsers)
	allCommands.register("agg", handleAgg)
	allCommands.register("addfeed", handleAddFeed)

	args := os.Args[1:]
	cmd := &command{
		name: args[0],
		args: args[1:],
	}

	if err = allCommands.run(currState, *cmd); err != nil {
		os.Exit(1)
	}
}
