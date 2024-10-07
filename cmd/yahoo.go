package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"os"

	"golang.org/x/net/html"
	"github.com/spf13/cobra"
)

// fetchYahooEmails fetches emails from Yahoo search results for the given domain
func fetchYahooEmails(domain string) []string {
	searchQuery := fmt.Sprintf("email \"@%s\"", domain)

	// Encode the search query
	encodedQuery := url.QueryEscape(searchQuery)

	// Construct the Yahoo search URL
	searchURL := fmt.Sprintf("https://search.yahoo.com/search?p=%s&pz=100", encodedQuery)

	// Create an HTTP client and set the User-Agent header
	client := &http.Client{}
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		fmt.Println("Error creating the request:", err)
		return nil
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36")

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching the URL:", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: status code", resp.StatusCode)
		return nil
	}

	// Parse the HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return nil
	}

	// Extract emails
	emails := extractyahooEmails(doc)

	// Process emails: split, lowercase, and filter
	return processyahooEmails(emails, domain)
}

// extractyahooEmails traverses the HTML nodes to find and collect email addresses
func extractyahooEmails(n *html.Node) []string {
	var emails []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			// Use regex to find email addresses in the text
			re := regexp.MustCompile(`([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`)
			matches := re.FindAllString(n.Data, -1)
			emails = append(emails, matches...)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)

	// Remove duplicates
	emailSet := make(map[string]struct{})
	for _, email := range emails {
		emailSet[email] = struct{}{}
	}

	var uniqueEmails []string
	for email := range emailSet {
		uniqueEmails = append(uniqueEmails, email)
	}

	return uniqueEmails
}

// processyahooEmails processes the extracted emails according to the specified steps
func processyahooEmails(emails []string, domain string) []string {
	var processed []string
	domainLower := strings.ToLower(domain) // To handle case insensitivity

	for _, email := range emails {
		emailLower := strings.ToLower(email) // Convert email to lowercase for comparison

		// Check if email ends with the specified domain
		if exactyahooMatch {
			// Match only if the email ends exactly with @domain
			if strings.HasSuffix(emailLower, "@"+domainLower) {
				processed = append(processed, emailLower)
			}
		} else {
			// Match if the email ends with the specified domain
			if strings.HasSuffix(emailLower, domainLower) {
				processed = append(processed, emailLower)
			}
		}
	}

	return processed
}

// yahooCmd represents the Yahoo command
var yahooCmd = &cobra.Command{
	Use:   "yahoo",
	Short: "Fetch emails from Yahoo search results for a given domain.",
	Long:  `This command allows you to fetch email addresses associated with a specified domain from Yahoo search results.

Examples:
echo "domain.com" | emailfinder yahoo
cat domains.txt | emailfinder yahoo
cat domains.txt | emailfinder yahoo -e`,
	Run: func(cmd *cobra.Command, args []string) {
		// Use a scanner to read domains from standard input
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			domain := strings.TrimSpace(scanner.Text())
			if domain == "" {
				continue
			}

			fmt.Printf("Fetching emails for domain: %s\n", domain)

			// Perform the email extraction for the domain
			emails := fetchYahooEmails(domain)

			if len(emails) == 0 {
				fmt.Printf("No emails found for domain: %s\n", domain)
			} else {
				for _, email := range emails {
					fmt.Println(email)
				}
				fmt.Printf("Found %d emails for: %s\n", len(emails), domain)
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
		}
	},
}

var exactyahooMatch bool // Declare exactyahooMatch as a global variable

func init() {
	rootCmd.AddCommand(yahooCmd)

	yahooCmd.Flags().BoolVarP(&exactyahooMatch, "exact-match", "e", false, "Match emails exactly with the domain (e.g., @domain)")
}
