// cmd/opml.go
package cmd

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// OPML represents the structure of an OPML document.
type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Head    Head     `xml:"head"`
	Body    Body     `xml:"body"`
}

// Head contains metadata such as the title of the document.
type Head struct {
	Title       string `xml:"title"`
	DateCreated string `xml:"dateCreated"`
}

// Body contains the main content of the OPML document.
type Body struct {
	Outlines []Outline `xml:"outline"`
}

// Outline represents a single outline entry in an OPML document.
type Outline struct {
	Text    string `xml:"text,attr"`
	Type    string `xml:"type,attr,omitempty"`
	XMLURL  string `xml:"xmlUrl,attr,omitempty"`
	HTMLURL string `xml:"htmlUrl,attr,omitempty"`
}

var opmlCmd = &cobra.Command{
	Use:   "opml",
	Short: "Generate an OPML file from posts",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generating OPML file...")
		GenerateOPML()
	},
}

func GenerateOPML() {
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

	opml := OPML{
		Version: "2.0",
		Head: Head{
			Title:       "Blackmail.dev OPML Feed",
			DateCreated: time.Now().Format(time.RFC3339),
		},
		Body: Body{
			Outlines: make([]Outline, 0),
		},
	}

	for _, post := range posts {
		outline := Outline{
			Text:    post.Title,
			Type:    "rss",
			XMLURL:  fmt.Sprintf("https://blackmail.dev/%s", post.Path),
			HTMLURL: fmt.Sprintf("https://blackmail.dev/%s", post.Path),
		}
		opml.Body.Outlines = append(opml.Body.Outlines, outline)
	}

	output, err := xml.MarshalIndent(opml, "", "  ")
	if err != nil {
		fmt.Printf("Error generating OPML XML: %s\n", err)
		return
	}

	err = os.WriteFile("blog.opml", output, 0644)
	if err != nil {
		fmt.Printf("Error writing blog.opml: %s\n", err)
		return
	}

	fmt.Println("OPML file generated successfully!")
}
