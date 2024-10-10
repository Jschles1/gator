package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Jschles1/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, _ command) error {
	fmt.Println("Collecting feeds every " + s.timeBetweenRequests.String())
	ticker := time.NewTicker(s.timeBetweenRequests)
	for range ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("error: feed name and/or URL not provided")
	}
	feedName := cmd.arguments[0]
	url := cmd.arguments[1]
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		Name:      feedName,
		Url:       url,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	newFeed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return fmt.Errorf("error: feed with URL already exists")
		}
		return err
	}

	// Create feed follow
	newFeedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    newFeed.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_, err = s.db.CreateFeedFollow(context.Background(), newFeedFollow)
	if err != nil {
		return err
	}

	fmt.Println("Feed successfully created:")
	fmt.Println("Name: ", newFeed.Name)
	fmt.Println("URL: ", newFeed.Url)
	return nil
}

func handlerFeeds(s *state, _ command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Feeds:")
	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			continue
		}
		fmt.Println("***************")
		fmt.Println("Feed Name: ", feed.Name)
		fmt.Println("URL: ", feed.Url)
		fmt.Println("By: ", user.Name)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("error: URL to follow not provided")
	}
	urlToFollow := cmd.arguments[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), urlToFollow)
	if err != nil {
		return err
	}
	newFeedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	feedFollowResult, err := s.db.CreateFeedFollow(context.Background(), newFeedFollow)
	if err != nil {
		return err
	}
	fmt.Println(feedFollowResult.UserName + " is now following " + feedFollowResult.FeedName)
	return nil
}

func handlerFollowing(s *state, _ command, user database.User) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	fmt.Println(user.Name + " is currently following:")
	for _, feed := range feedFollows {
		fmt.Println("- " + feed.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("error: feed URL not provided")
	}
	url := cmd.arguments[0]
	params := database.UnfollowParams{
		UserID: user.ID,
		Url:    url,
	}
	err := s.db.Unfollow(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Println(user.Name+" has successfully unfollowed feed at ", url)
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	}
	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Println("Your Posts:")
	for _, post := range posts {
		fmt.Println()
		fmt.Println(post.Title)
		fmt.Println()
		fmt.Println("Published at: ", post.PublishedAt)
		fmt.Println()
		fmt.Println(outputDescription(post.Description))
		fmt.Println()
		fmt.Println("****************************")
	}
	return nil
}
