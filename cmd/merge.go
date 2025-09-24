package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/liamgallear/md_split/internal/style"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:   "merge [splits-directory] [output-file]",
	Short: "Merge split files back into a single markdown file",
	Long: `Merge takes a directory containing split markdown files and combines them back into a single markdown file.
The split files should be numbered (01-, 02-, 03-, etc.) as created by the split command.`,
	Args: cobra.ExactArgs(2),
	RunE: mergeMarkdown,
}

func mergeMarkdown(cmd *cobra.Command, args []string) error {
	splitsDir := args[0]
	outputFile := args[1]
	
	// Check if splits directory exists
	if _, err := os.Stat(splitsDir); os.IsNotExist(err) {
		return fmt.Errorf("splits directory does not exist: %s", splitsDir)
	}

	// Read all markdown files from the splits directory
	files, err := filepath.Glob(filepath.Join(splitsDir, "*.md"))
	if err != nil {
		return fmt.Errorf("error reading splits directory: %v", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no markdown files found in splits directory: %s", splitsDir)
	}

	// Sort files to ensure correct order (relies on numeric prefix)
	// This works because the files are named with zero-padded numbers (01-, 02-, etc.)
	var sortedFiles []string
	for i := 1; i <= 99; i++ { // Support up to 99 sections
		prefix := fmt.Sprintf("%02d-", i)
		for _, file := range files {
			if strings.Contains(filepath.Base(file), prefix) {
				sortedFiles = append(sortedFiles, file)
				break
			}
		}
	}

	if len(sortedFiles) == 0 {
		return fmt.Errorf("no properly numbered split files found (expecting files with format 01-*, 02-*, etc.)")
	}

	// Merge the content
	fmt.Println(style.Info("Merging split files..."))
	var mergedContent strings.Builder
	bar := progressbar.Default(int64(len(sortedFiles)))
	
	for i, file := range sortedFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error reading file %s: %v", file, err)
		}
		
		// Add the content
		mergedContent.Write(content)
		
		// Add spacing between sections (except for the last one)
		if i < len(sortedFiles)-1 {
			mergedContent.WriteString("\n\n")
		}
		
		fmt.Println(style.FileMerged(filepath.Base(file)))
		bar.Add(1)
	}

	// Write the merged content to output file
	if err := os.WriteFile(outputFile, []byte(mergedContent.String()), 0644); err != nil {
		return fmt.Errorf("error writing output file %s: %v", outputFile, err)
	}

	fmt.Printf("\n%s\n", style.Summary(fmt.Sprintf("Successfully merged %d files into %s", len(sortedFiles), outputFile)))
	return nil
}