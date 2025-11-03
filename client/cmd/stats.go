package cmd

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/lincyaw/tools/client/pkg/client"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats [short code]",
	Short: "Get short link statistics",
	Long:  `Get statistics for the specified short code, including click count, creation time, etc.`,
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		code := args[0]
		c := client.NewClient(baseURL)

		color.Cyan("Getting statistics for short code '%s'...", code)

		stats, err := c.GetStats(code)
		if err != nil {
			color.Red("✗ Failed to get: %v", err)
			return
		}

		color.Green("\n✓ Statistics retrieved successfully!")
		fmt.Println()
		color.Cyan("Short code:      %s", stats.Code)
		color.Cyan("Original URL:   %s", stats.OriginalURL)
		color.Cyan("Click count:  %d", stats.ClickCount)
		color.Cyan("Created at:  %s", stats.CreatedAt.Format(time.RFC3339))
		if stats.LastAccessedAt != nil {
			color.Cyan("Last accessed:  %s", stats.LastAccessedAt.Format(time.RFC3339))
		} else {
			color.Cyan("Last accessed:  Never accessed")
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
