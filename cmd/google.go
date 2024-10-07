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

// fetchGoogleEmails fetches emails from Google search results for the given domain
func fetchGoogleEmails(domain string) []string {
	searchQuery := fmt.Sprintf("email \"%s\"", domain)

	// Encode the search query
	encodedQuery := url.QueryEscape(searchQuery)

	// Construct the Google search URL
	searchURL := fmt.Sprintf("https://www.google.com/search?num=100&pws=0&q=%s", encodedQuery)

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
	emails := extractEmails(doc)

	// Process emails: split, lowercase, and filter
	return processEmails(emails, domain)
}

// extractEmails traverses the HTML nodes to find and collect email addresses
func extractEmails(n *html.Node) []string {
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

// processEmails processes the extracted emails according to the specified steps
func processEmails(emails []string, domain string) []string {
	var processed []string
	domainLower := strings.ToLower(domain) // To handle case insensitivity

	for _, email := range emails {
		emailLower := strings.ToLower(email) // Convert email to lowercase for comparison

		// Check if email ends with the specified domain
		if exactgoogleMatch {
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

// googleCmd represents the google command
var googleCmd = &cobra.Command{
	Use:   "google",
	Short: "Fetch emails from Google search results for a given domain.",
	Long:  `This command allows you to fetch email addresses associated with a specified domain from Google search results.

Examples:
echo "domain.com" | emailfinder google
cat domains.txt | emailfinder google
cat domains.txt | emailfinder google -e`,
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
			emails := fetchGoogleEmails(domain)

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

var exactgoogleMatch bool // Declare exactgoogleMatch as a global variable

func init() {
	rootCmd.AddCommand(googleCmd)

	googleCmd.Flags().BoolVarP(&exactgoogleMatch, "exact-match", "e", false, "Match emails exactly with the domain (e.g., @domain)")
}
