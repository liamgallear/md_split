package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "md_split",
	Short: "A tool to split and merge markdown files based on H2 headings",
	Long: `md_split is a CLI tool for managing markdown files by splitting them into sections based on H2 headings
and merging split files back together.

Available commands:
  split   Split a markdown file into separate files based on H2 headings
  merge   Merge split files back into a single markdown file`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(splitCmd)
	rootCmd.AddCommand(mergeCmd)
}