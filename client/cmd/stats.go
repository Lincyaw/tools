package cmd

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/lincyaw/tools/client/pkg/client"
	"github.com/spf13/cobra"
)

var (
	detailedStats bool
	statsHours    int
)

var statsCmd = &cobra.Command{
	Use:   "stats [short code]",
	Short: "Get short link statistics",
	Long:  `Get statistics for the specified short code, including click count, creation time, etc.`,
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		code := args[0]
		c := client.NewClient(baseURL)

		if detailedStats {
			// Get detailed statistics
			color.Cyan("Getting detailed statistics for short code '%s'...", code)
			if statsHours > 0 {
				color.Cyan("(Looking back %d hours)", statsHours)
			}

			stats, err := c.GetDetailedStats(code, statsHours)
			if err != nil {
				color.Red("✗ Failed to get detailed statistics: %v", err)
				return
			}

			color.Green("\n✓ Detailed statistics retrieved successfully!")
			fmt.Println()

			// Basic info
			color.Cyan("═══════════════════════════════════════════════")
			color.Cyan("Basic Information")
			color.Cyan("═══════════════════════════════════════════════")
			fmt.Printf("Short code:      %s\n", stats.Code)
			fmt.Printf("Original URL:    %s\n", stats.OriginalURL)
			fmt.Printf("Total clicks:    %d\n", stats.TotalClicks)
			fmt.Printf("Unique IPs:      %d\n", stats.UniqueIPs)
			fmt.Printf("Created at:      %s\n", stats.CreatedAt.Format(time.RFC3339))
			if stats.LastAccessedAt != nil {
				fmt.Printf("Last accessed:   %s\n", stats.LastAccessedAt.Format(time.RFC3339))
			} else {
				fmt.Printf("Last accessed:   Never accessed\n")
			}
			fmt.Println()

			// Hourly statistics
			if len(stats.HourlyStats) > 0 {
				color.Cyan("═══════════════════════════════════════════════")
				color.Cyan("Hourly Statistics (Top 10)")
				color.Cyan("═══════════════════════════════════════════════")
				fmt.Printf("%-20s %12s %12s\n", "Hour", "Accesses", "Unique IPs")
				fmt.Println("-------------------------------------------------------")
				limit := 10
				if len(stats.HourlyStats) < limit {
					limit = len(stats.HourlyStats)
				}
				for i := 0; i < limit; i++ {
					h := stats.HourlyStats[i]
					fmt.Printf("%-20s %12d %12d\n",
						h.HourBucket.Format("2006-01-02 15:04"),
						h.AccessCount,
						h.UniqueIPs)
				}
				fmt.Println()
			}

			// Location statistics
			if len(stats.LocationStats) > 0 {
				color.Cyan("═══════════════════════════════════════════════")
				color.Cyan("Geographic Distribution (Top 10)")
				color.Cyan("═══════════════════════════════════════════════")
				fmt.Printf("%-20s %-20s %-20s %10s\n", "Country", "Region", "City", "Accesses")
				fmt.Println("--------------------------------------------------------------------------------")
				limit := 10
				if len(stats.LocationStats) < limit {
					limit = len(stats.LocationStats)
				}
				for i := 0; i < limit; i++ {
					l := stats.LocationStats[i]
					country := l.Country
					if country == "" {
						country = "Unknown"
					}
					region := l.Region
					if region == "" {
						region = "Unknown"
					}
					city := l.City
					if city == "" {
						city = "Unknown"
					}
					fmt.Printf("%-20s %-20s %-20s %10d\n", country, region, city, l.AccessCount)
				}
				fmt.Println()
			}

			// Recent accesses
			if len(stats.RecentAccesses) > 0 {
				color.Cyan("═══════════════════════════════════════════════")
				color.Cyan("Recent Accesses (Latest 10)")
				color.Cyan("═══════════════════════════════════════════════")
				limit := 10
				if len(stats.RecentAccesses) < limit {
					limit = len(stats.RecentAccesses)
				}
				for i := 0; i < limit; i++ {
					r := stats.RecentAccesses[i]
					fmt.Printf("\n[%s]\n", r.AccessTime.Format("2006-01-02 15:04:05"))
					fmt.Printf("  IP:       %s\n", r.IPAddress)
					location := fmt.Sprintf("%s, %s, %s", r.Country, r.Region, r.City)
					if r.Country == "" {
						location = "Unknown"
					}
					fmt.Printf("  Location: %s\n", location)
					if r.UserAgent != "" {
						userAgent := r.UserAgent
						if len(userAgent) > 60 {
							userAgent = userAgent[:60] + "..."
						}
						fmt.Printf("  UA:       %s\n", userAgent)
					}
				}
				fmt.Println()
			}

		} else {
			// Get basic statistics
			color.Cyan("Getting statistics for short code '%s'...", code)

			stats, err := c.GetStats(code)
			if err != nil {
				color.Red("✗ Failed to get: %v", err)
				return
			}

			color.Green("\n✓ Statistics retrieved successfully!")
			fmt.Println()
			color.Cyan("Short code:      %s", stats.Code)
			color.Cyan("Original URL:    %s", stats.OriginalURL)
			color.Cyan("Click count:     %d", stats.ClickCount)
			color.Cyan("Created at:      %s", stats.CreatedAt.Format(time.RFC3339))
			if stats.LastAccessedAt != nil {
				color.Cyan("Last accessed:   %s", stats.LastAccessedAt.Format(time.RFC3339))
			} else {
				color.Cyan("Last accessed:   Never accessed")
			}
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
	statsCmd.Flags().BoolVarP(&detailedStats, "detailed", "d", false, "Show detailed statistics including hourly data and location info")
	statsCmd.Flags().IntVarP(&statsHours, "hours", "H", 0, "Number of hours to look back (0 = all time)")
}
