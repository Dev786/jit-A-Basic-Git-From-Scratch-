/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"jit/internals"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func readTree(treeHash string, newFileRoot string) error {
	fmt.Fprintf(os.Stdout, "Reading Tree: %s\n", treeHash)
	// get the tree file
	treePath := filepath.Join(internals.ROOT, internals.DB_PATH, treeHash[:2], treeHash[2:])

	treeData, err := os.ReadFile(treePath)

	if err != nil {
		fmt.Printf("Error Reading Tree: %s\n", err)
	}

	// get all the files from the tree
	treeDataString := string(treeData)
	treeDataRows := strings.Split(treeDataString, "\n")

	for _, row := range treeDataRows {
		if row == "" {
			continue
		}

		rowData := strings.Split(row, " ")

		fileType := rowData[0]
		if fileType == internals.TREE {
			// get the tree hash
			treeHash := rowData[1]
			treeName := rowData[4]

			dir := filepath.Join(internals.ROOT, newFileRoot, treeName)
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Printf("Error Creating Directory: %s\n", err)
				return err
			}

			newFileRoot = filepath.Join(newFileRoot, treeName)

			readTree(treeHash, newFileRoot)
			continue
		}

		fileName := rowData[4]
		fileHash := rowData[1]

		// get the file
		filePath := filepath.Join(internals.ROOT, internals.DB_PATH, fileHash[:2], fileHash[2:])
		// fmt.Printf("File Path: %s\n", filePath)
		fileData, err := os.ReadFile(filePath)

		if len(fileData) == 0 {
			continue
		}

		if err != nil {
			// fmt.Printf("Error Reading File: %s\n", err)
			return err
		}

		// fmt.Print("File Data: ", string(fileData))

		// Extract the compressed data
		var blobType, blobId, blobSize, blobMode, blobName string
		var nItems int
		nItems, err = fmt.Sscanf(string(fileData), "%s %s %s %s %s\n", &blobType, &blobId, &blobSize, &blobMode, &blobName)

		if err != nil || nItems != 5 {
			fmt.Printf("Error Parsing Blob Data: %s\n", err)
			return err
		}

		// Find the index of the first newline to separate header from data
		newlineIndex := bytes.Index(fileData, []byte("\n"))
		if newlineIndex == -1 {
			fmt.Println("Invalid blob format: missing header newline")
			return fmt.Errorf("invalid blob format")
		}

		// Get the compressed data as bytes
		compressedDataBytes := fileData[newlineIndex+1:]

		// fmt.Printf("Blob Info - Type: %s, ID: %s, Size: %s, Mode: %s, Name: %s\n", blobType, blobId, blobSize, blobMode, blobName)
		// fmt.Printf("Compressed Data Bytes: %v\n", compressedDataBytes)

		// Directly decompress the data
		var out bytes.Buffer
		r, err := zlib.NewReader(bytes.NewReader(compressedDataBytes))
		if err != nil {
			fmt.Printf("Error Creating Decompressor: %s\n", err)
			return err
		}
		_, err = io.Copy(&out, r)
		if err != nil {
			fmt.Printf("Error Decompressing Data: %s\n", err)
			return err
		}
		r.Close()

		// Ensure the directory exists
		destFilePath := filepath.Join(internals.CHECKOUT_DATA_PATH, fileName)
		if err := os.MkdirAll(filepath.Dir(destFilePath), 0755); err != nil {
			return err
		}

		// Create the file if it does not exist
		file, err := os.OpenFile(destFilePath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Printf("Error Creating File: %s\n", err)
			return err
		}
		defer file.Close()

		// Write decompressed data to the file
		if _, err := file.Write(out.Bytes()); err != nil {
			fmt.Printf("Error Writing File: %s\n", err)
			return err
		}
	}
	return nil
}

func copyAllFileToRootFromCommit(commitHash string) error {
	// get the commit file
	commitPath := filepath.Join(internals.ROOT, internals.DB_PATH, commitHash[:2], commitHash[2:])
	commitData, err := os.ReadFile(commitPath)

	if err != nil {
		fmt.Printf("Error Reading Commit: %s\n", err)
		return err
	}

	// get the tree hash from the commit
	commitDataString := string(commitData)
	treeRow := strings.Split(commitDataString, "\n")[0]
	treeHash := strings.Split(treeRow, " ")[1]
	err = readTree(treeHash, filepath.Join(internals.ROOT, "test_dir"))
	return err
}

func handleNewBranch(branch string) error {
	// check if the branch exists in the path
	newBranchPath := filepath.Join(internals.ROOT, internals.REF_PATH, internals.REFS_HEADS, branch)

	refsPath := filepath.Join(internals.REF_PATH, internals.REFS_HEADS)

	// Check if the file exists in the path or the path exists
	if _, err := os.Stat(newBranchPath); os.IsNotExist(err) {
		fmt.Printf("Creating Branch: %s\n", branch)

		// Create the directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(newBranchPath), os.ModePerm); err != nil {
			fmt.Printf("Error Creating Directory: %s\n", err)
			return err
		}

		// Create the new branch file
		file, err := os.Create(newBranchPath)
		if err != nil {
			fmt.Printf("Error Creating Branch File: %s\n", err)
			return err
		}
		defer file.Close()

		// change the HEAD to point to the new branch
		headPath := filepath.Join(internals.ROOT, internals.HEAD)

		// read the HEAD file
		headFileData, err := os.ReadFile(headPath)
		if err != nil {
			fmt.Printf("Error Reading HEAD: %s\n", err)
			return err
		}

		// replace the refs: from the headFileData
		headRef := strings.Replace(string(headFileData), "ref: ", "", 1)

		// read the headRef file
		headRefData, err := os.ReadFile(filepath.Join(internals.ROOT, headRef))

		// add this headRefData to the new branch file
		if err != nil {
			fmt.Printf("Error Reading HEAD Ref: %s\n", err)
			return err
		}

		// add refsPath to HEAD file
		newHeadFileContent := fmt.Sprintf("ref: %s", filepath.Join(refsPath, branch))
		err = os.WriteFile(headPath, []byte(newHeadFileContent), 0644)

		if err != nil {
			fmt.Printf("Error Writing HEAD: %s\n", err)
			return err
		}

		err = os.WriteFile(newBranchPath, headRefData, 0644)

		if err != nil {
			fmt.Printf("Error Writing Branch: %s\n", err)
			return err
		}

	} else if err != nil {
		fmt.Printf("Error Checking Branch: %s\n", err)
		return err
	}

	return nil
}

func handleBranchChange(branch string) error {
	// point the HEAD to the new branch
	path := filepath.Join(internals.ROOT, internals.REF_PATH, internals.REFS_HEADS, branch)

	// check if the branch exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Branch does not exist")
		return err
	} else if err != nil {
		fmt.Printf("Error Checking Branch: %s\n", err)
		return err
	}

	newHeadFileContent := fmt.Sprintf("ref: %s", filepath.Join(internals.REF_PATH, internals.REFS_HEADS, branch))

	// write the new HEAD file
	err := os.WriteFile(filepath.Join(internals.ROOT, internals.HEAD), []byte(newHeadFileContent), 0644)

	if err != nil {
		fmt.Printf("Error Writing HEAD: %s\n", err)
		return err
	}

	// get the commit hash from the branch
	branchData, err := os.ReadFile(path)

	if err != nil {
		fmt.Printf("Error Reading Branch: %s\n", err)
		return err
	}

	commitHash := string(branchData)

	copyAllFileToRootFromCommit(commitHash)

	return nil
}

// checkoutCmd represents the checkout command
var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// get the value from b
		branch, err := cmd.Flags().GetString("branch")
		// get the value from args
		var branchToChange string = ""
		if len(args) > 0 {
			branchToChange = args[0]
		}

		if err != nil {
			fmt.Printf("Error Fetching Branch: %s\n", err)
		}

		if branch == "" && branchToChange == "" {
			fmt.Println("Branch name is required")
			return
		}

		if branch != "" {
			fmt.Printf("Creating and Checking out to new branch: %s\n", branch)
			// if the branch exists

			branchPath := filepath.Join(internals.ROOT, internals.REF_PATH, internals.REFS_HEADS, branch)
			if _, err := os.Stat(branchPath); err == nil {
				fmt.Println("Branch already exists")
				return
			}
			handleNewBranch(branch)
		} else {
			fmt.Printf("Checking out to branch: %s\n", branchToChange)
			// check the HEAD if the branch to change == HEAD
			headPath := filepath.Join(internals.ROOT, internals.HEAD)
			headData, err := os.ReadFile(headPath)

			if err != nil {
				fmt.Printf("Error Reading HEAD: %s\n", err)
				return
			}

			headBranch := strings.Replace(string(headData), "ref: refs/heads/", "", 1)

			fmt.Printf("HEAD Branch: %s, Branch To Change: %s\n", headBranch, branchToChange)

			if headBranch == branchToChange {
				fmt.Println("Already on the branch")
				return
			}
			// change branch
			handleBranchChange(branchToChange)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkoutCmd)
	checkoutCmd.Flags().StringP("branch", "b", "", "Create and checkout to a new branch")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
