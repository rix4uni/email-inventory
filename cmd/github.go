package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// githubCmd represents the github command
var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "A brief description of your command",
	Long: `A longer descriptionn of your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("github called")
	},
}

func init() {
	rootCmd.AddCommand(githubCmd)

	githubCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// https://github.com/search?q=%40dell.com&type=code