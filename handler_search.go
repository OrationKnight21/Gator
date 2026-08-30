package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/OrationKnight21/gator/internal/database"
)

func handlerSearch(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: %s <search_query> [limit]", cmd.name)
	}

	query := cmd.args[0]
	limit := int32(10) // default limit

	if len(cmd.args) > 1 {
		parsedLimit, err := strconv.Atoi(cmd.args[1])
		if err == nil && parsedLimit > 0 {
			limit = int32(parsedLimit)
		}
	}

	posts, err := s.db.SearchPostsForUser(context.Background(), database.SearchPostsForUserParams{
		UserID: user.ID,
		Title:  query,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("could not search posts: %w", err)
	}

	if len(posts) == 0 {
		fmt.Printf("No posts found matching '%s'\n", query)
		return nil
	}

	fmt.Printf("Found %d posts matching '%s':\n\n", len(posts), query)
	for _, post := range posts {
		fmt.Printf("* %s\n", post.Title)
		fmt.Printf("  Feed: %s\n", post.FeedName)
		fmt.Printf("  Link: %s\n", post.Url)
		if post.Description.Valid {
			// Print snippet of description
			desc := strings.TrimSpace(post.Description.String)
			if len(desc) > 100 {
				desc = desc[:97] + "..."
			}
			fmt.Printf("  Summary: %s\n", desc)
		}
		fmt.Println()
	}

	return nil
}
