package cmd

import (
	"github.com/fatih/color"
	"github.com/lincyaw/tools/client/pkg/test"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run complete test suite",
	Long:  `Run the complete test suite, testing all API endpoint functions.`,
	Run: func(_ *cobra.Command, _ []string) {
		var tester *test.Tester
		if insecureSkipVerify {
			tester = test.NewTesterWithInsecureSkipVerify(baseURL, verbose)
			if verbose {
				color.Yellow("⚠ Warning: Skipping TLS certificate verification")
			}
		} else {
			tester = test.NewTester(baseURL, verbose)
		}
		tester.RunAllTests()
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
