package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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

var splitCmd = &cobra.Command{
	Use:   "split [file]",
	Short: "Split a markdown file based on H2 headings",
	Long: `Split takes a markdown file and splits it into separate files based on H2 headings.
Each split file will be prefixed with a section number (01, 02, 03, etc.) and saved in a 'splits' directory relative to the original file.`,
	Args: cobra.ExactArgs(1),
	RunE: splitMarkdown,
}

var mergeCmd = &cobra.Command{
	Use:   "merge [splits-directory] [output-file]",
	Short: "Merge split files back into a single markdown file",
	Long: `Merge takes a directory containing split markdown files and combines them back into a single markdown file.
The split files should be numbered (01-, 02-, 03-, etc.) as created by the split command.`,
	Args: cobra.ExactArgs(2),
	RunE: mergeMarkdown,
}

func splitMarkdown(cmd *cobra.Command, args []string) error {
	inputFile := args[0]
	
	// Check if input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", inputFile)
	}

	// Read the markdown file
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Parse the content and split by H2 headings
	sections, err := parseMarkdownSections(string(content))
	if err != nil {
		return fmt.Errorf("error parsing markdown: %v", err)
	}

	if len(sections) == 0 {
		return fmt.Errorf("no H2 headings found in the file")
	}

	// Create splits directory
	inputDir := filepath.Dir(inputFile)
	splitsDir := filepath.Join(inputDir, "splits")
	
	if err := os.MkdirAll(splitsDir, 0755); err != nil {
		return fmt.Errorf("error creating splits directory: %v", err)
	}

	// Write each section to a separate file
	for i, section := range sections {
		filename := fmt.Sprintf("%02d-%s.md", i+1, sanitizeFilename(section.title))
		filepath := filepath.Join(splitsDir, filename)
		
		if err := os.WriteFile(filepath, []byte(section.content), 0644); err != nil {
			return fmt.Errorf("error writing file %s: %v", filename, err)
		}
		
		fmt.Printf("Created: %s\n", filename)
	}

	fmt.Printf("Successfully split markdown into %d files in %s\n", len(sections), splitsDir)
	return nil
}

type section struct {
	title   string
	content string
}

func parseMarkdownSections(content string) ([]section, error) {
	var sections []section
	lines := strings.Split(content, "\n")
	
	// Regex to match H2 headings (## heading)
	h2Regex := regexp.MustCompile(`^##\s+(.+)$`)
	
	var currentSection *section
	var currentContent []string
	
	for _, line := range lines {
		if h2Regex.MatchString(line) {
			// Save the previous section if it exists
			if currentSection != nil {
				currentSection.content = strings.Join(currentContent, "\n")
				sections = append(sections, *currentSection)
			}
			
			// Start a new section
			matches := h2Regex.FindStringSubmatch(line)
			title := strings.TrimSpace(matches[1])
			currentSection = &section{title: title}
			currentContent = []string{line} // Include the heading in the content
		} else if currentSection != nil {
			// Add line to current section
			currentContent = append(currentContent, line)
		}
		// If no current section and not an H2, ignore the line (content before first H2)
	}
	
	// Save the last section
	if currentSection != nil {
		currentSection.content = strings.Join(currentContent, "\n")
		sections = append(sections, *currentSection)
	}
	
	return sections, nil
}

func sanitizeFilename(filename string) string {
	// Replace spaces with hyphens and remove invalid characters
	reg := regexp.MustCompile(`[^\w\-_.]`)
	sanitized := reg.ReplaceAllString(filename, "")
	sanitized = strings.ReplaceAll(sanitized, " ", "-")
	sanitized = strings.ToLower(sanitized)
	
	// Limit length to avoid filesystem issues
	if len(sanitized) > 50 {
		sanitized = sanitized[:50]
	}
	
	return sanitized
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
	var mergedContent strings.Builder
	
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
		
		fmt.Printf("Merged: %s\n", filepath.Base(file))
	}

	// Write the merged content to output file
	if err := os.WriteFile(outputFile, []byte(mergedContent.String()), 0644); err != nil {
		return fmt.Errorf("error writing output file %s: %v", outputFile, err)
	}

	fmt.Printf("Successfully merged %d files into %s\n", len(sortedFiles), outputFile)
	return nil
}