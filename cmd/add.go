/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

// adding a new comment to show the commit and checkout
import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [file_name]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fileName := args[0]

		// check if file exists
		_, err := os.Stat(fileName)

		if os.IsNotExist(err) {
			fmt.Printf("File %s does not exist\n", fileName)
		} else if err != nil {
			fmt.Printf("Error Checking file: %s\n", err)
		}
		// write go code to add file to staging area
		// write go code to create blob object
		// write go code to write blob object to .jit/objects directory
		// write go code to write blob object to staging area
		// write go code to update index file

	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
