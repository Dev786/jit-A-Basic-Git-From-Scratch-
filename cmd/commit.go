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

func createTree(path string, database *internals.Database) (*internals.Tree, error) {
	// get message from commit -m "message" flag
	var entries []*internals.Entry
	tree := internals.NewTree()
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip the root directory itself
		if filePath == path {
			return nil
		}
		// Skip .jit directory
		if info.IsDir() && filepath.Base(filePath) == ".jit" {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			data, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			// Get the file name
			fileName := filepath.Base(filePath)

			blob := internals.NewBlob(
				fileName,
				data,
			)

			database.Store(blob)

			var entry internals.Entry
			entry.SetType(blob.GetType())
			entry.SetOid(blob.GetOid())
			entry.SetSize(blob.GetSize())
			entry.SetName(blob.GetName())

			entries = append(entries, &entry)

			fmt.Printf("Stored blob for file: %s\n", filePath)
		} else {
			// Recursively create tree for directories
			subTree, err := createTree(filePath, database)
			if err != nil {
				return err
			}
			if subTree != nil {
				var entry internals.Entry
				entry.SetType(internals.TREE)
				entry.SetOid(subTree.GetOid())
				entry.SetSize(subTree.GetSize())
				entry.SetName(filepath.Base(filePath))
				entries = append(entries, &entry)
				fmt.Printf("Stored Tree for directory: %s\n", filePath)
			}
		}
		return nil
	})

	// create a new Tree object
	tree.SetEntries(entries)

	// store the tree object
	database.Store(tree)

	if err != nil {
		fmt.Printf("Error Reading Files: %s\n", err)
	}
	return tree, err
}

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		// check for root folder
		if _, err := os.Stat(internals.ROOT); os.IsNotExist(err) {
			fmt.Println("Not a jit repository, please run 'jit init' to initialize a new repository")
			return
		}
		message, err := cmd.Flags().GetString("message")
		if err != nil {
			fmt.Printf("Error Fetching Message: %s\n", err)
		}

		database, err := internals.NewDatabase(filepath.Join(internals.ROOT, internals.DB_PATH))

		if err != nil {
			fmt.Printf("Error Creating Database: %s\n", err)
		}
		rootDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error Fetching Root Directory: %s\n", err)
			return
		}
		var tree *internals.Tree

		tree, err = createTree(rootDir, database)

		if err != nil {
			fmt.Printf("Error Creating Tree: %s\n", err)
		}

		fmt.Printf("Tree: %s\n", tree.GetOid())

		commit := internals.NewCommit(tree, nil, internals.AUTHOR, internals.EMAIL, internals.COMMITTER, message)

		database.Store(commit)

		// read the head file and get the refs path
		headFile := filepath.Join(internals.ROOT, internals.HEAD)

		fmt.Printf("HEAD File: %s\n", headFile)

		// read the head file
		headContent, err := os.ReadFile(headFile)

		fmt.Printf("HEAD File Content:  %s\n", headContent)

		if err != nil {
			fmt.Printf("Error Reading HEAD: %s\n", err)
		}

		// extract refs: refs/heads/main
		ref := strings.Split(string(headContent), "ref: ")[1]

		// get the refs path
		refsPath := filepath.Join(internals.ROOT, ref)

		// write the commit oid to the refs file
		err = os.WriteFile(refsPath, []byte(commit.GetOid()), 0644)

		if err != nil {
			fmt.Printf("Error Writing Commit to Refs: %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	commitCmd.Flags().StringP("message", "m", "", "Commit Message")
}
