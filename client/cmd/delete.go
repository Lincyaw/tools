package cmd

import (
	"github.com/fatih/color"
	"github.com/lincyaw/tools/client/pkg/client"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [short code]",
	Short: "Delete short link",
	Long:  `Delete the specified short link.`,
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		code := args[0]
		c := client.NewClient(baseURL)

		color.Cyan("Deleting short code '%s'...", code)

		err := c.DeleteShortCode(code)
		if err != nil {
			color.Red("✗ Deletion failed: %v", err)
			return
		}

		color.Green("✓ Short code '%s' deleted successfully!", code)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
