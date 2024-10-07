package cmd

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

// fetchEmails fetches and parses the HTML for a given domain, extracting emails
func fetchEmails(domain string) ([]string, error) {
	url := fmt.Sprintf("https://www.skymem.info/srch?q=%s", domain)

	// Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch url: %v", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Extract emails by finding the <a> tags
	var emails []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		email := s.Text()
		if strings.Contains(email, "@") {
			emails = append(emails, email)
		}
	})

	return emails, nil
}

// saveEmailsToFile saves unique emails to a file using the 'unew' command
func saveEmailsToFile(emails []string, domain string) error {
	filePath := fmt.Sprintf("Emails/%s.txt", domain)

	// Ensure the Emails directory exists
	err := os.MkdirAll("Emails", 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Prepare to run the unew command
	cmd := exec.Command("unew", "-q", filePath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %v", err)
	}

	// Start the command
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start unew command: %v", err)
	}

	// Write emails to the stdin of the unew command
	for _, email := range emails {
		_, err := stdin.Write([]byte(email + "\n"))
		if err != nil {
			return fmt.Errorf("failed to write email to stdin: %v", err)
		}
	}
	stdin.Close()

	// Wait for the unew command to finish
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("unew command failed: %v", err)
	}

	fmt.Printf("Emails saved to %s\n", filePath)
	return nil
}

// skymemCmd represents the skymem command
var skymemCmd = &cobra.Command{
	Use:   "skymem",
	Short: "Fetch emails from a domain using skymem",
	Long: `This command fetches emails associated with a given domain from skymem.info.
You can provide one or more domain names via standard input and optionally use the -save flag to save unique emails to a file.

Examples:
echo "domain.com" | emailfinder skymem
cat domains.txt | emailfinder skymem
cat domains.txt | emailfinder skymem -s`,
	Run: func(cmd *cobra.Command, args []string) {
		scanner := bufio.NewScanner(os.Stdin)
		firstDomain := true // Flag to check if it's the first domain

		for scanner.Scan() {
			domain := scanner.Text()
			if domain == "" {
				continue
			}

			if !firstDomain {
				fmt.Println() // Print a blank line between domains
			}
			firstDomain = false

			fmt.Printf("Fetching emails for domain: %s\n", domain)

			emails, err := fetchEmails(domain)
			if err != nil {
				log.Printf("Error fetching emails for %s: %v\n", domain, err)
				continue
			}

			emailCount := len(emails)
			if emailCount > 0 {
				fmt.Printf("Found %d emails for: %s\n", emailCount, domain)
				for _, email := range emails {
					fmt.Println(email)
				}

				// If the -save flag is provided, save emails to file
				if saveFlag {
					err := saveEmailsToFile(emails, domain)
					if err != nil {
						log.Printf("Error saving emails for %s: %v\n", domain, err)
					}
				}
			} else {
				fmt.Printf("No emails found for %s\n", domain)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Error reading input: %v\n", err)
		}
	},
}

var saveFlag bool

func init() {
	rootCmd.AddCommand(skymemCmd)

	skymemCmd.Flags().BoolVarP(&saveFlag, "save", "s", false, "Save unique emails to Emails/domain.txt")
}
