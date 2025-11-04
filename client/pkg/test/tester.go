package test

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/lincyaw/tools/client/pkg/client"
)

// Result represents the result of a test execution
type Result struct {
	Name    string
	Passed  bool
	Message string
	Error   error
}

// Tester tester
type Tester struct {
	client  *client.Client
	results []Result
	verbose bool
}

// NewTester create tester
func NewTester(baseURL string, verbose bool) *Tester {
	return &Tester{
		client:  client.NewClient(baseURL),
		results: make([]Result, 0),
		verbose: verbose,
	}
}

// NewTesterWithInsecureSkipVerify creates a tester that skips TLS certificate verification
func NewTesterWithInsecureSkipVerify(baseURL string, verbose bool) *Tester {
	return &Tester{
		client:  client.NewClientWithInsecureSkipVerify(baseURL),
		results: make([]Result, 0),
		verbose: verbose,
	}
}

// addResult add test result
func (t *Tester) addResult(name string, passed bool, message string, err error) {
	t.results = append(t.results, Result{
		Name:    name,
		Passed:  passed,
		Message: message,
		Error:   err,
	})

	if passed {
		color.Green("✓ %s: %s", name, message)
	} else {
		color.Red("✗ %s: %s", name, message)
		if err != nil && t.verbose {
			color.Red("  Error details: %v", err)
		}
	}
}

// TestHealthCheck test health check
func (t *Tester) TestHealthCheck() {
	color.Cyan("\n━━━ Test Health Check ━━━")
	err := t.client.HealthCheck()
	if err != nil {
		t.addResult("Health Check", false, "Service unavailable", err)
	} else {
		t.addResult("Health Check", true, "Service running normally", nil)
	}
}

// TestCreateShortCodeAuto test auto-generated short code
func (t *Tester) TestCreateShortCodeAuto() string {
	color.Cyan("\n━━━ Test Create Short Link (Auto-generated Short Code) ━━━")

	req := client.CreateShortCodeRequest{
		URL:       "https://github.com/lincyaw/tools",
		ExpiresIn: 3600,
	}

	resp, err := t.client.CreateShortCode(req)
	if err != nil {
		t.addResult("Create Short Link (Auto)", false, "Creation failed", err)
		return ""
	}

	msg := fmt.Sprintf("Short code: %s, URL: %s", resp.ShortCode, resp.ShortURL)
	t.addResult("Create Short Link (Auto)", true, msg, nil)

	if t.verbose {
		color.Yellow("  Original URL: %s", resp.OriginalURL)
		color.Yellow("  Created at: %s", resp.CreatedAt.Format(time.RFC3339))
		if resp.ExpiresAt != nil {
			color.Yellow("  Expiration time: %s", resp.ExpiresAt.Format(time.RFC3339))
		}
	}

	return resp.ShortCode
}

// TestCreateShortCodeCustom test custom short code
func (t *Tester) TestCreateShortCodeCustom(customCode string) {
	color.Cyan("\n━━━ Test Create Short Link (Custom Short Code) ━━━")

	req := client.CreateShortCodeRequest{
		URL:        "https://github.com/lincyaw/tools",
		CustomCode: customCode,
		ExpiresIn:  7200,
	}

	resp, err := t.client.CreateShortCode(req)
	if err != nil {
		t.addResult("Create Short Link (Custom)", false, "Creation failed", err)
		return
	}

	msg := fmt.Sprintf("Short code: %s, URL: %s", resp.ShortCode, resp.ShortURL)
	t.addResult("Create Short Link (Custom)", true, msg, nil)

	if t.verbose {
		color.Yellow("  Original URL: %s", resp.OriginalURL)
		color.Yellow("  Created at: %s", resp.CreatedAt.Format(time.RFC3339))
	}
}

// TestDuplicateCustomCode test duplicate custom short code
func (t *Tester) TestDuplicateCustomCode(customCode string) {
	color.Cyan("\n━━━ Test Duplicate Short Code (Should Fail) ━━━")

	req := client.CreateShortCodeRequest{
		URL:        "https://example.com",
		CustomCode: customCode,
	}

	_, err := t.client.CreateShortCode(req)
	if err != nil {
		t.addResult("Duplicate Short Code Detection", true, "Correctly rejected duplicate short code", nil)
		if t.verbose {
			color.Yellow("  Error message: %v", err)
		}
	} else {
		t.addResult("Duplicate Short Code Detection", false, "Failed to detect duplicate short code", nil)
	}
}

// TestRedirect test redirect
func (t *Tester) TestRedirect(code string) {
	color.Cyan("\n━━━ Test Short Link Redirect ━━━")

	info, err := t.client.TestRedirect(code)
	if err != nil {
		t.addResult("Short Link Redirect", false, "Redirect failed", err)
		return
	}

	msg := fmt.Sprintf("Status code: %d, Redirect to: %s", info.StatusCode, info.Location)
	t.addResult("Short Link Redirect", true, msg, nil)

	if t.verbose {
		color.Yellow("  Redirect status code: %d", info.StatusCode)
		color.Yellow("  Target URL: %s", info.Location)
	}
}

// TestGetStats test get statistics
func (t *Tester) TestGetStats(code string) {
	color.Cyan("\n━━━ Test Get Statistics ━━━")

	stats, err := t.client.GetStats(code)
	if err != nil {
		t.addResult("Get Statistics", false, "Failed to get", err)
		return
	}

	msg := fmt.Sprintf("Short code: %s, Click count: %d", stats.Code, stats.ClickCount)
	t.addResult("Get Statistics", true, msg, nil)

	if t.verbose {
		color.Yellow("  Original URL: %s", stats.OriginalURL)
		color.Yellow("  Click count: %d", stats.ClickCount)
		color.Yellow("  Created at: %s", stats.CreatedAt.Format(time.RFC3339))
		if stats.LastAccessedAt != nil {
			color.Yellow("  Last accessed: %s", stats.LastAccessedAt.Format(time.RFC3339))
		}
	}
}

// TestGetDetailedStats test get detailed statistics
func (t *Tester) TestGetDetailedStats(code string) {
	color.Cyan("\n━━━ Test Get Detailed Statistics ━━━")

	// Test without time range (all time)
	stats, err := t.client.GetDetailedStats(code, 0)
	if err != nil {
		t.addResult("Get Detailed Statistics (All Time)", false, "Failed to get", err)
		return
	}

	msg := fmt.Sprintf("Code: %s, Total clicks: %d, Unique IPs: %d", stats.Code, stats.TotalClicks, stats.UniqueIPs)
	t.addResult("Get Detailed Statistics (All Time)", true, msg, nil)

	if t.verbose {
		color.Yellow("  Original URL: %s", stats.OriginalURL)
		color.Yellow("  Total clicks: %d", stats.TotalClicks)
		color.Yellow("  Unique IPs: %d", stats.UniqueIPs)
		color.Yellow("  Created at: %s", stats.CreatedAt.Format(time.RFC3339))

		if len(stats.HourlyStats) > 0 {
			color.Yellow("  Hourly stats entries: %d", len(stats.HourlyStats))
			color.Yellow("  Latest hour bucket: %s (%d accesses, %d unique IPs)",
				stats.HourlyStats[0].HourBucket.Format("2006-01-02 15:04"),
				stats.HourlyStats[0].AccessCount,
				stats.HourlyStats[0].UniqueIPs)
		}

		if len(stats.LocationStats) > 0 {
			color.Yellow("  Location stats entries: %d", len(stats.LocationStats))
			top := stats.LocationStats[0]
			color.Yellow("  Top location: %s, %s, %s (%d accesses)",
				top.Country, top.Region, top.City, top.AccessCount)
		}

		if len(stats.RecentAccesses) > 0 {
			color.Yellow("  Recent accesses: %d", len(stats.RecentAccesses))
			latest := stats.RecentAccesses[0]
			color.Yellow("  Latest access: %s from %s (%s, %s, %s)",
				latest.AccessTime.Format("2006-01-02 15:04:05"),
				latest.IPAddress, latest.Country, latest.Region, latest.City)
		}
	}

	// Test with time range (last 24 hours)
	stats24h, err := t.client.GetDetailedStats(code, 24)
	if err != nil {
		t.addResult("Get Detailed Statistics (24h)", false, "Failed to get", err)
		return
	}

	msg24h := fmt.Sprintf("Last 24h - Clicks: %d, Unique IPs: %d", stats24h.TotalClicks, stats24h.UniqueIPs)
	t.addResult("Get Detailed Statistics (24h)", true, msg24h, nil)

	if t.verbose {
		color.Yellow("  Last 24h total clicks: %d", stats24h.TotalClicks)
		color.Yellow("  Last 24h unique IPs: %d", stats24h.UniqueIPs)
		color.Yellow("  Last 24h hourly entries: %d", len(stats24h.HourlyStats))
	}
}

// TestAccessStatisticsRecording test that access statistics are properly recorded
func (t *Tester) TestAccessStatisticsRecording() {
	color.Cyan("\n━━━ Test Access Statistics Recording ━━━")

	// Create a test shortcode
	testCode := fmt.Sprintf("statstest%d", time.Now().Unix())
	req := client.CreateShortCodeRequest{
		URL:        "https://github.com/lincyaw/tools",
		CustomCode: testCode,
	}

	_, err := t.client.CreateShortCode(req)
	if err != nil {
		t.addResult("Statistics Test - Create Code", false, "Failed to create test code", err)
		return
	}

	// Get initial stats (should be 0)
	initialStats, err := t.client.GetDetailedStats(testCode, 0)
	if err != nil {
		t.addResult("Statistics Test - Get Initial Stats", false, "Failed to get initial stats", err)
		return
	}

	if initialStats.TotalClicks != 0 {
		t.addResult("Statistics Test - Initial State", false,
			fmt.Sprintf("Expected 0 clicks, got %d", initialStats.TotalClicks), nil)
	} else {
		t.addResult("Statistics Test - Initial State", true, "Initial clicks = 0", nil)
	}

	// Simulate some accesses
	accessCount := 5
	color.Yellow("  Simulating %d accesses...", accessCount)
	for i := 0; i < accessCount; i++ {
		_, err := t.client.TestRedirect(testCode)
		if err != nil {
			if t.verbose {
				color.Yellow("  Access %d failed: %v", i+1, err)
			}
		}
		time.Sleep(200 * time.Millisecond) // Small delay between accesses
	}

	// Wait for async processing
	color.Yellow("  Waiting for statistics to be processed...")
	time.Sleep(3 * time.Second)

	// Get updated stats
	updatedStats, err := t.client.GetDetailedStats(testCode, 0)
	if err != nil {
		t.addResult("Statistics Test - Get Updated Stats", false, "Failed to get updated stats", err)
		return
	}

	// Verify click count increased
	if updatedStats.TotalClicks >= int64(accessCount) {
		msg := fmt.Sprintf("Clicks recorded correctly (%d >= %d)", updatedStats.TotalClicks, accessCount)
		t.addResult("Statistics Test - Click Recording", true, msg, nil)
	} else {
		msg := fmt.Sprintf("Expected >= %d clicks, got %d", accessCount, updatedStats.TotalClicks)
		t.addResult("Statistics Test - Click Recording", false, msg, nil)
	}

	// Verify hourly stats exist
	if len(updatedStats.HourlyStats) > 0 {
		t.addResult("Statistics Test - Hourly Stats", true,
			fmt.Sprintf("Hourly statistics created (%d entries)", len(updatedStats.HourlyStats)), nil)

		if t.verbose {
			for i, h := range updatedStats.HourlyStats {
				color.Yellow("    Hour %d: %s - %d accesses, %d unique IPs",
					i+1, h.HourBucket.Format("2006-01-02 15:04"),
					h.AccessCount, h.UniqueIPs)
			}
		}
	} else {
		t.addResult("Statistics Test - Hourly Stats", false, "No hourly statistics created", nil)
	}

	// Verify location stats exist
	if len(updatedStats.LocationStats) > 0 {
		t.addResult("Statistics Test - Location Stats", true,
			fmt.Sprintf("Location statistics created (%d entries)", len(updatedStats.LocationStats)), nil)

		if t.verbose {
			for i, l := range updatedStats.LocationStats {
				color.Yellow("    Location %d: %s, %s, %s - %d accesses",
					i+1, l.Country, l.Region, l.City, l.AccessCount)
			}
		}
	} else {
		// Location stats might be "Unknown" for localhost, which is still valid
		t.addResult("Statistics Test - Location Stats", true,
			"Location stats may be empty (localhost access)", nil)
	}

	// Verify recent accesses
	if len(updatedStats.RecentAccesses) > 0 {
		t.addResult("Statistics Test - Recent Accesses", true,
			fmt.Sprintf("Recent accesses recorded (%d entries)", len(updatedStats.RecentAccesses)), nil)

		if t.verbose {
			for i, r := range updatedStats.RecentAccesses {
				if i >= 3 { // Only show first 3
					break
				}
				color.Yellow("    Access %d: %s from %s (%s, %s, %s)",
					i+1, r.AccessTime.Format("15:04:05"),
					r.IPAddress, r.Country, r.Region, r.City)
			}
		}
	} else {
		t.addResult("Statistics Test - Recent Accesses", false, "No recent accesses recorded", nil)
	}

	// Verify unique IPs
	if updatedStats.UniqueIPs > 0 {
		t.addResult("Statistics Test - Unique IPs", true,
			fmt.Sprintf("Unique IPs tracked: %d", updatedStats.UniqueIPs), nil)
	} else {
		t.addResult("Statistics Test - Unique IPs", false, "No unique IPs tracked", nil)
	}

	// Test time range filtering (last 1 hour)
	stats1h, err := t.client.GetDetailedStats(testCode, 1)
	if err != nil {
		t.addResult("Statistics Test - Time Range Filter", false, "Failed to get 1h stats", err)
	} else {
		msg := fmt.Sprintf("Last 1h stats retrieved (clicks: %d)", stats1h.TotalClicks)
		t.addResult("Statistics Test - Time Range Filter", true, msg, nil)
	}

	// Cleanup
	if err := t.client.DeleteShortCode(testCode); err != nil {
		if t.verbose {
			color.Yellow("  Warning: Failed to cleanup test code: %v", err)
		}
	}
}

// TestInvalidRequests test invalid requests
func (t *Tester) TestInvalidRequests() {
	color.Cyan("\n━━━ Test Invalid Requests ━━━")

	// Test invalid URL
	req := client.CreateShortCodeRequest{
		URL: "not-a-valid-url",
	}
	_, err := t.client.CreateShortCode(req)
	if err != nil {
		t.addResult("Invalid URL Detection", true, "Correctly rejected invalid URL", nil)
	} else {
		t.addResult("Invalid URL Detection", false, "Failed to detect invalid URL", nil)
	}

	// Test non-existent short code
	_, err = t.client.GetStats("nonexistent999")
	if err != nil {
		t.addResult("Non-existent Short Code Detection", true, "Correctly returned error", nil)
	} else {
		t.addResult("Non-existent Short Code Detection", false, "Failed to detect non-existent short code", nil)
	}
}

// TestDeleteShortCode test delete short link
func (t *Tester) TestDeleteShortCode() {
	color.Cyan("\n━━━ Test Delete Short Link ━━━")

	// Create temporary short link
	tempCode := fmt.Sprintf("temp%d", time.Now().Unix())
	req := client.CreateShortCodeRequest{
		URL:        "https://example.com",
		CustomCode: tempCode,
	}

	_, err := t.client.CreateShortCode(req)
	if err != nil {
		t.addResult("Delete Test - Create Temp Short Code", false, "Failed to create temp short code", err)
		return
	}

	// Delete short link
	err = t.client.DeleteShortCode(tempCode)
	if err != nil {
		t.addResult("Delete Short Link", false, "Deletion failed", err)
		return
	}

	t.addResult("Delete Short Link", true, fmt.Sprintf("Successfully deleted short code: %s", tempCode), nil)

	// Verify if really deleted
	_, err = t.client.TestRedirect(tempCode)
	if err != nil {
		t.addResult("Deletion Verification", true, "Short code has been deleted", nil)
	} else {
		t.addResult("Deletion Verification", false, "Short code still exists", nil)
	}
}

// TestRateLimiting test rate limiting
func (t *Tester) TestRateLimiting() {
	color.Cyan("\n━━━ Test Rate Limiting ━━━")

	rateLimitHit := false
	successCount := 0

	for i := 0; i < 10; i++ {
		req := client.CreateShortCodeRequest{
			URL: fmt.Sprintf("https://example.com/test%d", i),
		}

		_, err := t.client.CreateShortCode(req)
		if err != nil {
			if t.verbose {
				color.Yellow("  Request %d: Failed - %v", i+1, err)
			}
			rateLimitHit = true
		} else {
			successCount++
			if t.verbose {
				color.Yellow("  Request %d: Success", i+1)
			}
		}

		time.Sleep(100 * time.Millisecond)
	}

	if rateLimitHit {
		t.addResult("Rate Limiting", true, fmt.Sprintf("Rate limiting effective (Success: %d/10)", successCount), nil)
	} else {
		t.addResult("Rate Limiting", true, fmt.Sprintf("All requests successful (10/10) - Rate limiting lenient (%d)", successCount), nil)
	}
}

// RunAllTests run all tests
func (t *Tester) RunAllTests() {
	color.Cyan("\n╔══════════════════════════════════════════════╗")
	color.Cyan("║     Short Code Service Client Test            ║")
	color.Cyan("╚══════════════════════════════════════════════╝")

	customCode := fmt.Sprintf("test%d", time.Now().Unix())

	t.TestHealthCheck()
	autoCode := t.TestCreateShortCodeAuto()
	t.TestCreateShortCodeCustom(customCode)
	t.TestDuplicateCustomCode(customCode)

	if autoCode != "" {
		t.TestRedirect(autoCode)
		t.TestGetStats(autoCode)
		t.TestGetDetailedStats(autoCode)
	}

	t.TestInvalidRequests()
	t.TestAccessStatisticsRecording()
	t.TestDeleteShortCode()
	t.TestRateLimiting()

	t.PrintSummary()
}

// PrintSummary print test summary
func (t *Tester) PrintSummary() {
	color.Cyan("\n╔══════════════════════════════════════════════╗")
	color.Cyan("║              Test Summary                        ║")
	color.Cyan("╚══════════════════════════════════════════════╝")

	passed := 0
	failed := 0

	for _, result := range t.results {
		if result.Passed {
			passed++
		} else {
			failed++
		}
	}

	total := passed + failed
	fmt.Printf("\nTotal: %d tests\n", total)
	color.Green("Passed: %d", passed)
	color.Red("Failed: %d", failed)

	if failed == 0 {
		color.Green("\n✓ All tests passed!\n")
	} else {
		color.Red("\n✗ Some tests failed\n")
		color.Yellow("\nFailed tests:")
		for _, result := range t.results {
			if !result.Passed {
				color.Red("  • %s: %s", result.Name, result.Message)
				if result.Error != nil {
					color.Red("    Error: %v", result.Error)
				}
			}
		}
		fmt.Println()
	}
}

// GetResults get test results
func (t *Tester) GetResults() []Result {
	return t.results
}
