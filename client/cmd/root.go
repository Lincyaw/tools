package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	baseURL string
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "shortcode-client",
	Short: "Short link service client tool",
	Long: `Short link service client tool for testing and interacting with the short link service.

Supported operations:
  - Create short links (auto-generated or custom short codes)
  - Get short link statistics
  - Delete short links
  - Run complete test suite`,
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&baseURL, "url", "u", "http://localhost", "Service base URL")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Show verbose output")
}
