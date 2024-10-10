package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Jschles1/gator/internal/database"
	"github.com/google/uuid"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		// Get the user here
		user, err := s.db.GetUser(context.Background(), s.c.CurrentUserName)
		if err != nil {
			return err
		}

		// Call the original handler, passing in the retrieved user
		return handler(s, cmd, user)
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("error: username not provided")
	}
	username := cmd.arguments[0]
	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return fmt.Errorf("error: user doesn't exist")
		}
		return err
	}
	s.c.SetUser(user.Name)
	fmt.Println(user.Name + " has been logged in")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("error: username not provided")
	}
	username := cmd.arguments[0]
	params := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      username,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	newUser, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return fmt.Errorf("error: user already exists")
		}
		return err
	}
	fmt.Println(newUser.Name + " has been successfully registered.")
	fmt.Println(newUser)
	s.c.SetUser(newUser.Name)
	return nil
}

func handlerReset(s *state, _ command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Users table successfully reset.")
	return nil
}

func handlerUsers(s *state, _ command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	currentUser := s.c.CurrentUserName
	for _, user := range users {
		if user.Name == currentUser {
			fmt.Println("* " + user.Name + " (current)")
		} else {
			fmt.Println("* " + user.Name)
		}
	}
	return nil
}
