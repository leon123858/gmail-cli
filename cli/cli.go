package cli

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	numEmails int
	rootCmd   = &cobra.Command{
		Use:   "gmail-cli",
		Short: "A CLI tool for managing and reading emails from multiple Gmail accounts",
		Long:  `This CLI tool allows you to manage multiple Gmail accounts and read recent emails using OAuth2 authentication.`,
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(runCmd)

	configCmd.AddCommand(configAddCmd)
	configCmd.AddCommand(configDeleteCmd)
	configCmd.AddCommand(configSetCmd)

	runCmd.AddCommand(runReadCmd)
	runReadCmd.Flags().IntVarP(&numEmails, "count", "n", 20, "Number of emails to read from each account")
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `Add or remove email accounts, and set the path to credentials.json file.`,
}

var configAddCmd = &cobra.Command{
	Use:   "add <email>",
	Short: "Add a new email account",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		email := args[0]
		accounts := viper.GetStringSlice("accounts")
		if contains(accounts, email) {
			fmt.Printf("Account %s already exists.\n", email)
			return
		}

		accounts = append(accounts, email)
		viper.Set("accounts", accounts)
		if err := viper.WriteConfig(); err != nil {
			fmt.Printf("Error writing config: %v\n", err)
			return
		}

		fmt.Printf("Added account: %s\n", email)
	},
}

var configDeleteCmd = &cobra.Command{
	Use:   "delete <email>",
	Short: "Delete an email account",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		email := args[0]
		accounts := viper.GetStringSlice("accounts")
		if !contains(accounts, email) {
			fmt.Printf("Account %s does not exist.\n", email)
			return
		}

		accounts = remove(accounts, email)
		viper.Set("accounts", accounts)
		if err := viper.WriteConfig(); err != nil {
			fmt.Printf("Error writing config: %v\n", err)
			return
		}

		fmt.Printf("Deleted account: %s\n", email)
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <id> <secret>",
	Short: "Set the client ID and secret on GCP OAuth2 credentials",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		secret := args[1]

		viper.Set("id", id)
		viper.Set("secret", secret)
		if err := viper.WriteConfig(); err != nil {
			fmt.Printf("Error writing config: %v\n", err)
			return
		}

		fmt.Println("Client ID and secret set.")
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Gmail operations",
	Long:  `Execute Gmail operations such as reading emails.`,
}

var runReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Read emails from configured accounts",
	Run: func(cmd *cobra.Command, args []string) {
		accounts := viper.GetStringSlice("accounts")
		if len(accounts) == 0 {
			fmt.Println("No accounts configured. Use 'gmail-cli config add <email>' to add an account.")
			return
		}

		var wg sync.WaitGroup
		for _, account := range accounts {
			wg.Add(1)
			go readEmails(account, numEmails, &wg)
		}
		wg.Wait()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
