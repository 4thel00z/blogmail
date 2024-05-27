// cmd/create.go
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type Post struct {
	Title   string `json:"title"`
	Date    string `json:"date"`
	Path    string `json:"path"`
	Summary string `json:"summary"`
}

func readPosts(filename string) ([]Post, error) {
	var posts []Post
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []Post{}, nil
		}
		return nil, err
	}
	err = json.Unmarshal(data, &posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func writePosts(filename string, posts []Post) error {
	data, err := json.Marshal(posts)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new blog post",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating a new blog post...")
		createPost()
	},
}

func createPost() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the post title: ")
	postTitle, _ := reader.ReadString('\n')
	postTitle = strings.TrimSpace(postTitle)

	fmt.Print("Enter the post summary: ")
	postSummary, _ := reader.ReadString('\n')
	postSummary = strings.TrimSpace(postSummary)

	fmt.Print("Enter the post date (YYYY-MM-DD), or press enter to use today's date: ")
	dateInput, _ := reader.ReadString('\n')
	dateInput = strings.TrimSpace(dateInput)

	var postDate time.Time
	if dateInput == "" {
		postDate = time.Now()
	} else {
		var err error
		postDate, err = time.Parse("2006-01-02", dateInput)
		if err != nil {
			fmt.Println("Invalid date format. Using today's date instead.")
			postDate = time.Now()
		}
	}

	postDirName := fmt.Sprintf("%s_%s", postDate.Format("20060102"), strings.ReplaceAll(postTitle, " ", "_"))
	postPath := fmt.Sprintf("posts/%s", postDirName)
	if err := os.MkdirAll(postPath, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory: %s\n", err)
		return
	}

	filePath := fmt.Sprintf("%s/index.html", postPath)
	content := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <link rel="stylesheet" href="../../css/style.css">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/styles/github-dark-dimmed.min.css">
    <script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/highlight.min.js"></script>
    <!-- and it's easy to individually load additional languages -->
    <script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/go.min.js"></script>
    <script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/python.min.js"></script>
    <script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/bash.min.js"></script>
    <script src="https://cdn.jsdelivr.net/gh/highlightjs/cdn-release@11.9.0/build/languages/makefile.min.js"></script>

    <script>
        document.addEventListener('DOMContentLoaded', (event) => {
            hljs.highlightAll();
        });
    </script>
</head>
<body>
    <button class="back-button" onclick="window.history.back();"><i class="fa fa-arrow-left"></i> Back</button>
    <h1>%s</h1>
    <p>Date: %s</p>
    <p>%s</p>
    <!-- Example placeholder for code -->
    <pre><code class="language-go">// Your Go code here</code></pre>
</body>
</html>`, postTitle, postTitle, postDate.Format("January 2, 2006"), postSummary)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		fmt.Printf("Error writing to file: %s\n", err)
		return
	}

	posts, err := readPosts("posts.json")
	if err != nil {
		fmt.Printf("Error reading posts: %s\n", err)
		return
	}

	newPost := Post{
		Title:   postTitle,
		Date:    postDate.Format("2006-01-02"),
		Path:    postPath,
		Summary: postSummary,
	}
	posts = append(posts, newPost)

	slices.SortFunc(posts, func(a, b Post) int {
		return strings.Compare(b.Date, a.Date)
	})

	if err := writePosts("posts.json", posts); err != nil {
		fmt.Printf("Error writing posts: %s\n", err)
		return
	}

	fmt.Printf("New post created successfully: %s\n", filePath)
}
