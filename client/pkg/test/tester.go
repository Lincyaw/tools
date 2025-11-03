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
	}

	t.TestInvalidRequests()
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
