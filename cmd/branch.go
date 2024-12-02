/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"jit/internals"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func walkForBranch(path string, branches *[]string) error {
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			*branches = append(*branches, filePath)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := filepath.Join(internals.ROOT, internals.REF_PATH, internals.REFS_HEADS)
		var branches []string
		err := walkForBranch(path, &branches)

		if err != nil {
			fmt.Printf("Error Walking for Branches: %s\n", err)
		}

		// print all branches
		for _, branch := range branches {
			// remove the path from the branches and print
			branchName := strings.TrimPrefix(branch, path+string(os.PathSeparator))
			fmt.Printf("%s\n", branchName)
		}
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// branchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// branchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
