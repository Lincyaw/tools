package cmd

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/lincyaw/tools/client/pkg/client"
	"github.com/spf13/cobra"
)

var (
	url        string
	customCode string
	expiresIn  int
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create short link",
	Long:  `Create a new short link, can automatically generate short code or use custom short code.`,
	Run: func(_ *cobra.Command, _ []string) {
		c := client.NewClient(baseURL)

		req := client.CreateShortCodeRequest{
			URL:        url,
			CustomCode: customCode,
			ExpiresIn:  expiresIn,
		}

		color.Cyan("Creating short link...")
		if verbose {
			color.Yellow("URL: %s", url)
			if customCode != "" {
				color.Yellow("Custom short code: %s", customCode)
			}
			if expiresIn > 0 {
				color.Yellow("Expiration time: %d hours", expiresIn)
			}
		}

		resp, err := c.CreateShortCode(req)
		if err != nil {
			color.Red("✗ Creation failed: %v", err)
			return
		}

		color.Green("\n✓ Short link created successfully!")
		fmt.Println()
		color.Cyan("Short code:      %s", resp.ShortCode)
		color.Cyan("Short link:    %s", resp.ShortURL)
		color.Cyan("Original URL:   %s", resp.OriginalURL)
		color.Cyan("Created at:  %s", resp.CreatedAt.Format(time.RFC3339))
		if resp.ExpiresAt != nil {
			color.Cyan("Expires at:  %s", resp.ExpiresAt.Format(time.RFC3339))
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&url, "long-url", "l", "", "The long URL to shorten (required)")
	createCmd.Flags().StringVarP(&customCode, "code", "c", "", "Custom short code (optional, auto-generated if not provided)")
	createCmd.Flags().IntVarP(&expiresIn, "expires", "e", 0, "Expiration time (hours, optional)")

	if err := createCmd.MarkFlagRequired("long-url"); err != nil {
		panic(fmt.Sprintf("failed to mark flag as required: %v", err))
	}
}
