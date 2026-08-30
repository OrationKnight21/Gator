package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/OrationKnight21/gator/internal/config"
	"github.com/OrationKnight21/gator/internal/database"
	"github.com/google/uuid"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}
type command struct {
	name string
	args []string
}
type commands struct {
	comNames map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.comNames[cmd.name]
	if !ok {
		return errors.New("the command does not exist")
	}
	return handler(s, cmd)
}
func (c *commands) register(name string, f func(*state, command) error) {
	c.comNames[name] = f
}
func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUser(context.Background())
	if err != nil {
		return err
	}
	fmt.Print("reset successful\n")
	return nil
}
func users(s *state, cmd command) error {
	usr, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, val := range usr {
		if s.cfg.CurrentUserName == val.Name {
			fmt.Printf("* %v (current)\n", val.Name)
		} else {
			fmt.Printf("* %v\n", val.Name)
		}
	}
	return nil
}
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("insufficient arguments")
	}
	usr, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	err = s.cfg.SetUser(usr.Name)
	if err != nil {
		return err
	}
	fmt.Print("The user has been set\n")
	return nil
}
func registerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("insufficient arguments")
	}
	usr, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.args[0]})
	if err != nil {
		return err
	}
	err = s.cfg.SetUser(usr.Name)
	if err != nil {
		return err
	}
	fmt.Println("the user was created")
	log.Printf("%+v\n", usr)
	return nil
}
func addfeed(s *state, cmd command, usr database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("insufficient arguments")
	}
	name := cmd.args[0]
	url := cmd.args[1]
	feedusr, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: name, Url: url, UserID: usr.ID})
	if err != nil {
		return err
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: feedusr.CreatedAt, UpdatedAt: feedusr.UpdatedAt, UserID: feedusr.UserID, FeedID: feedusr.ID})
	if err != nil {
		return err
	}
	fmt.Printf("%+v \n", feedusr)
	return nil
}
func handlerFeeds(s *state, cmd command) error {
	feed, err := s.db.GetFeedName(context.Background())
	if err != nil {
		return err
	}
	for i := range feed {
		fmt.Printf("%+v \n %+v \n %+v \n", feed[i].FeedName, feed[i].Url, feed[i].UserName)
	}
	return nil
}
func follow(s *state, cmd command, usr database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("insufficient arguments")
	}
	url := cmd.args[0]
	feedname, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return err
	}
	feedfollow, err := s.db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    usr.ID, FeedID: feedname.ID})
	if err != nil {
		return err
	}
	fmt.Printf("%v \n", feedfollow.FeedName)
	fmt.Printf("%v \n", feedfollow.UserName)
	return nil
}
func following(s *state, cmd command, usr database.User) error {
	followingFeed, err := s.db.GetFeedFollowsForUser(context.Background(), usr.ID)
	if err != nil {
		return err
	}
	for i := range followingFeed {
		fmt.Printf("%v \n", followingFeed[i].FeedName)
	}
	return nil
}
func handlerUnfollow(s *state, cmd command, usr database.User) error {
	url := cmd.args[0]
	feedId, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return err
	}
	return s.db.DeleteFeedFollows(context.Background(),
		database.DeleteFeedFollowsParams{UserID: usr.ID, FeedID: feedId.ID})
}
func browse(s *state, cmd command, usr database.User) error {
	limit := 2
	if len(cmd.args) < 1 {
		if specifiedLimit, err := strconv.Atoi(cmd.args[0]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("invalid limit: %v", err)
		}
	}
	if len(cmd.args) > 1 {
		return errors.New("wrong number of arguments entered")
	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: usr.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), usr.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}
	return nil
}
