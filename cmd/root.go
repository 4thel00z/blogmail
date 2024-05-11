// cmd/root.go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "os"
)

var rootCmd = &cobra.Command{
    Use:   "blackmail",
    Short: "Blackmail.dev CLI for blog management",
    Long:  `Blackmail.dev CLI to manage blog posts and feeds.`,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func init() {
    rootCmd.AddCommand(rssCmd)
    rootCmd.AddCommand(opmlCmd)
    rootCmd.AddCommand(createCmd)
}
