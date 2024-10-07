package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// fetchGithubSavedEmails fetches emails for the given domain and prints the results
func fetchGithubSavedEmails(domain string) {
	// Construct the URL with the domain
	url := fmt.Sprintf("https://raw.githubusercontent.com/rix4uni/EmailFinder/main/Emails/%s.txt", domain)
	fmt.Printf("Fetching emails for domain: %s\n", domain)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch email list: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Customize the error message for a 404 status code
	if resp.StatusCode == http.StatusNotFound {
		fmt.Fprintf(os.Stderr, "Error: %s.txt file does not exist\n", domain)
		return
	} else if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Error: Received status code %d\n", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response body: %v\n", err)
		return
	}

	// Process the email list
	emails := strings.Split(string(body), "\n")
	var foundEmails []string
	for _, email := range emails {
		email = strings.TrimSpace(email) // Trim whitespace
		if email != "" { // Ignore empty lines
			foundEmails = append(foundEmails, email)
		}
	}

	// Print found emails
	if len(foundEmails) > 0 {
		for _, email := range foundEmails {
			fmt.Println(email)
		}
		fmt.Printf("Found %d emails for: %s\n", len(foundEmails), domain)
	} else {
		fmt.Printf("No emails found for domain: %s\n", domain)
	}

	fmt.Println() // Add a blank line between results for each domain
}

// savedCmd represents the saved command
var savedCmd = &cobra.Command{
	Use:   "saved",
	Short: "Fetches emails for given domains using https://github.com/rix4uni/EmailFinder/tree/main/Emails",
	Long:  `This command fetches and displays emails for the specified domains from a remote text file.
In github.com/rix4uni/EmailFinder/tree/main/Emails.

Examples:
echo "domain.com" | emailfinder saved
cat domains.txt | emailfinder saved`,
	Run: func(cmd *cobra.Command, args []string) {
		scanner := bufio.NewScanner(os.Stdin)

		// Read domains from stdin
		for scanner.Scan() {
			domain := strings.TrimSpace(scanner.Text())
			fetchGithubSavedEmails(domain)
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(savedCmd)
}