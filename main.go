package main

import (
	"fmt"
	"os"

	"github.com/0xF000D/goenv/pkg/services"
	"github.com/0xF000D/goenv/pkg/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   fmt.Sprintf("%s [command]", utils.APP_NAME),
	Short: fmt.Sprintf("%s - inject env variables at runtime", utils.APP_NAME),
	Long:  fmt.Sprintf("%s is a next-generation tool for managing environment variables", utils.APP_NAME),
	Args:  cobra.ArbitraryArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// fmt.Println("In presistentPreRun func")
	},
	Run: func(cmd *cobra.Command, args []string) {
		// pass
	},
}

func init() {
	// Sub commands
	rootCmd.AddCommand(runCmd())
	rootCmd.AddCommand(encryptCmd())
	rootCmd.AddCommand(decryptCmd())
}

func main() {
	// p := tea.NewProgram(helpers.NewSpinner("Encrypting"))
	// if _, err := p.Run(); err != nil {
	// 	fmt.Println(err)
	// }
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runCmd() *cobra.Command {
	var envFiles []string
	var envStrings []string
	// var overload bool

	cmd := &cobra.Command{
		Use:   "run -- <command you want to exceute>",
		Short: "Inject env variables on runtime",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(args)
			fmt.Println(envFiles)
			fmt.Println("Running command with injected env variables...")
		},
	}

	cmd.Flags().StringSliceVarP(&envStrings, "env", "e", []string{}, "environment variable(s) set as string")
	cmd.Flags().StringSliceVarP(&envFiles, "env-file", "f", []string{}, "path(s) to your env file(s)")

	return cmd
}

func encryptCmd() *cobra.Command {
	var envFilesPaths []string
	var envKeysFilePath string

	cmd := &cobra.Command{
		Use:   "encrypt",
		Short: "Convert plain .env file(s) to encrypted .env file(s)",
		Run: func(cmd *cobra.Command, args []string) {
			services.Encrypt(envFilesPaths, envKeysFilePath)
		},
	}

	cmd.Flags().StringSliceVarP(&envFilesPaths, "env-file", "f", []string{".env"}, "path(s) to your plain env file(s)")
	cmd.Flags().StringVarP(&envKeysFilePath, "env-keys-file", "fk", ".env.keys", "paths to you .env.keys file")

	return cmd
}

func decryptCmd() *cobra.Command {
	var envFilesPaths []string
	var envKeysFilePath string

	cmd := &cobra.Command{
		Use:   "decrypt",
		Short: "Convert encrypted .env file(s) to plain .env file(s)",
		Run: func(cmd *cobra.Command, args []string) {
			services.Decrypt(envFilesPaths, envKeysFilePath)
		},
	}

	cmd.Flags().StringSliceVarP(&envFilesPaths, "env-file", "f", []string{".env"}, "path(s) to your encrypted env file(s)")
	cmd.Flags().StringVarP(&envKeysFilePath, "env-keys-file", "fk", ".env.keys", "paths to you .env.keys file")

	return cmd
}
