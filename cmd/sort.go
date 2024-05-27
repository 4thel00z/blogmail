package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"slices"
	"strings"
)

var sortCmd = &cobra.Command{
	Use:   "sort",
	Short: "Sort the blog posts by date",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Sorting blog posts by date...")
		sortPosts()
	},
}

func sortPosts() {
	posts, err := readPosts("posts.json")
	if err != nil {
		fmt.Printf("Error reading posts: %s\n", err)
		return
	}

	slices.SortFunc(posts, func(a, b Post) int {
		return strings.Compare(b.Date, a.Date)
	})

	if err := writePosts("posts.json", posts); err != nil {
		fmt.Printf("Error writing posts: %s\n", err)
		return
	}

	fmt.Println("Posts sorted successfully")
}

func init() {
	rootCmd.AddCommand(sortCmd)
}
