package cmd

import (
	"github.com/lincyaw/tools/client/pkg/test"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run complete test suite",
	Long:  `Run the complete test suite, testing all API endpoint functions.`,
	Run: func(_ *cobra.Command, _ []string) {
		tester := test.NewTester(baseURL, verbose)
		tester.RunAllTests()
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
