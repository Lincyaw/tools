package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/lincyaw/tools/client/pkg/client"
	"github.com/spf13/cobra"
)

var redirectCmd = &cobra.Command{
	Use:   "redirect [short code]",
	Short: "Test short link redirect",
	Long:  `Test the redirect function of the specified short code, display the redirect status code and target URL.`,
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		code := args[0]
		var c *client.Client
		if insecureSkipVerify {
			c = client.NewClientWithInsecureSkipVerify(baseURL)
		} else {
			c = client.NewClient(baseURL)
		}

		color.Cyan("Testing redirect for short code '%s'...", code)

		info, err := c.TestRedirect(code)
		if err != nil {
			color.Red("✗ Redirect test failed: %v", err)
			return
		}

		color.Green("\n✓ Redirect test successful!")
		fmt.Println()
		color.Cyan("Short code:        %s", code)
		color.Cyan("Status code:      %d", info.StatusCode)
		color.Cyan("Redirect to:    %s", info.Location)
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(redirectCmd)
}
