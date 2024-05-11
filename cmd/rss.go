// cmd/rss.go
package cmd

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// RSS is the root type for an RSS document.
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

// Channel represents an RSS channel.
type Channel struct {
	Title         string `xml:"title"`
	Link          string `xml:"link"`
	Description   string `xml:"description"`
	LastBuildDate string `xml:"lastBuildDate"`
	Items         []Item `xml:"item"`
}

// Item represents a single entry in an RSS feed.
type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Guid        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
}

var rssCmd = &cobra.Command{
	Use:   "rss",
	Short: "Generate an RSS feed from posts",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generating RSS feed...")
		GenerateRSS()
	},
}

func GenerateRSS() {
	var posts []Post
	data, err := os.ReadFile("posts.json")
	if err != nil {
		fmt.Printf("Error reading posts.json: %s\n", err)
		return
	}
	err = json.Unmarshal(data, &posts)
	if err != nil {
		fmt.Printf("Error unmarshalling posts: %s\n", err)
		return
	}

	rss := RSS{
		Version: "2.0",
		Channel: Channel{
			Title:         "Blackmail.dev Blog",
			Link:          "https://blackmail.dev",
			Description:   "Technical insights and software development stories.",
			LastBuildDate: time.Now().Format(time.RFC1123Z),
			Items:         make([]Item, 0),
		},
	}

	for _, post := range posts {
		item := Item{
			Title:       post.Title,
			Link:        fmt.Sprintf("https://blackmail.dev/%s", post.Path),
			Description: post.Summary,
			Guid:        fmt.Sprintf("https://blackmail.dev/%s", post.Path),
			PubDate:     formatDate(post.Date),
		}
		rss.Channel.Items = append(rss.Channel.Items, item)
	}

	output, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		fmt.Printf("Error generating RSS XML: %s\n", err)
		return
	}

	err = os.WriteFile("feed.xml", output, 0644)
	if err != nil {
		fmt.Printf("Error writing feed.xml: %s\n", err)
		return
	}

	fmt.Println("RSS feed generated successfully!")
}

func formatDate(dateStr string) string {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return ""
	}
	return t.Format(time.RFC1123Z)
}
