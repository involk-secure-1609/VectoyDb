package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	interactive bool
	rootCmd     = &cobra.Command{
		Use:   "mycli",
		Short: "My CLI application with interactive mode",
		Long:  `A longer description of your application`,
		Run: func(cmd *cobra.Command, args []string) {
			// This is executed when the root command is called
			fmt.Printf("Config File: %s\n", viper.GetString("config"))
			fmt.Printf("Debug Mode: %v\n", viper.GetBool("debug"))

			// Check if interactive mode is enabled
			if interactive {
				runInteractiveMode()
			}
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mycli.yaml)")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "enable debug mode")
	rootCmd.PersistentFlags().BoolVarP(&interactive, "interactive", "i", false, "run in interactive mode")

	// Bind flags to viper
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	// Add any additional commands
	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".mycli" (without extension)
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".mycli")
	}

	// Read environment variables
	viper.AutomaticEnv()

	// Read in config
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("MyApp v1.0")
	},
}

func runInteractiveMode() {
	fmt.Println("Welcome to Interactive Mode")
	fmt.Println("Type 'help' for commands or 'exit' to quit")
	fmt.Println("----------------------------------------")

	scanner := bufio.NewScanner(os.Stdin)

	// Main loop
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		switch input {
		case "exit", "quit":
			fmt.Println("Goodbye!")
			return
		case "help":
			showHelp()
		case "":
			// Skip empty input
		default:
			processCommand(input)
		}
	}
}

func showHelp() {
	fmt.Println("Available interactive commands:")
	fmt.Println("  help  - Show this help message")
	fmt.Println("  exit  - Exit interactive mode")
	fmt.Println("  debug - Show debug status")
	// Add more commands
}

func processCommand(cmd string) {
	parts := strings.Fields(cmd)
	command := parts[0]
	args := parts[1:]

	switch command {
	case "debug":
		fmt.Printf("Debug mode: %v\n", viper.GetBool("debug"))
	case "echo":
		fmt.Println(strings.Join(args, " "))
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}

func main() {
	// interactive = true
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
