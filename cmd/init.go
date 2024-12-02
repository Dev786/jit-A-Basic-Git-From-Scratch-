/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"jit/internals"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// write go code get root directory
		dir, err := os.Getwd()

		if err != nil {
			fmt.Printf("Error Fetching Root: %s\n", err)
		}

		// write go code to create .jit directory
		jitDir := filepath.Join(dir, internals.ROOT)
		err = os.Mkdir(jitDir, internals.FILE_MODE)

		if err != nil {
			fmt.Printf("Error Creating .jit Directory: %s\n", err)
		}

		// create objects directory
		err = os.Mkdir(filepath.Join(jitDir, internals.DB_PATH), internals.FILE_MODE)

		if err != nil {
			fmt.Printf("Error Creating Objects Directory: %s\n", err)
		}

		// create refs directory
		err = os.Mkdir(filepath.Join(jitDir, internals.REF_PATH), internals.FILE_MODE)

		if err != nil {
			fmt.Printf("Error Creating Refs Directory: %s\n", err)
		}

		// create refs folder
		err = os.Mkdir(filepath.Join(jitDir, internals.REF_PATH), internals.FILE_MODE)

		if err != nil {
			fmt.Printf("Error Creating Refs Folder: %s\n", err)
		}

		// create refs/heads folder
		err = os.Mkdir(filepath.Join(jitDir, internals.REF_PATH, internals.REFS_HEADS), internals.FILE_MODE)

		if err != nil {
			fmt.Printf("Error Creating refs/head Folder: %s\n", err)
		}

		// create an empty main file in refs/heads
		_, err = os.Create(filepath.Join(jitDir, internals.REF_PATH, internals.REFS_HEADS, internals.MAIN_BRANCH))

		if err != nil {
			fmt.Printf("Error Creating Main File: %s\n", err)
		}

		// create refs/heads/main file
		_, err = os.Create(filepath.Join(jitDir, internals.REF_PATH, "heads", "main"))

		if err != nil {
			fmt.Printf("Error Creating Main File: %s\n", err)
		}

		// add the reference to the main branch
		err = os.WriteFile(filepath.Join(jitDir, internals.HEAD), []byte("ref: refs/heads/main"), 0644)

		if err != nil {
			fmt.Printf("Error Writing HEAD File: %s\n", err)
		}
		fmt.Println("Initialized empty Jit repository in", jitDir)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
