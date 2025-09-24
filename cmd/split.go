package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/liamgallear/md_split/internal/style"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var splitCmd = &cobra.Command{
	Use:   "split [file]",
	Short: "Split a markdown file based on H2 headings",
	Long: `Split takes a markdown file and splits it into separate files based on H2 headings.
Each split file will be prefixed with a section number (01, 02, 03, etc.) and saved in a 'splits' directory relative to the original file.`,
	Args: cobra.ExactArgs(1),
	RunE: splitMarkdown,
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
		return fmt.Errorf("no content found in the file")
	}

	// Create splits directory
	inputDir := filepath.Dir(inputFile)
	splitsDir := filepath.Join(inputDir, "splits")
	
	if err := os.MkdirAll(splitsDir, 0755); err != nil {
		return fmt.Errorf("error creating splits directory: %v", err)
	}

	// Write each section to a separate file
	fmt.Println(style.Info("Writing split files..."))
	bar := progressbar.Default(int64(len(sections)))
	
	for i, section := range sections {
		var filename string
		if section.title == "preamble" {
			filename = "00-preamble.md"
		} else {
			// Adjust index for non-preamble sections
			sectionNum := i + 1
			if len(sections) > 0 && sections[0].title == "preamble" {
				sectionNum = i // Don't increment if first section is preamble
			}
			filename = fmt.Sprintf("%02d-%s.md", sectionNum, sanitizeFilename(section.title))
		}
		filepath := filepath.Join(splitsDir, filename)
		
		if err := os.WriteFile(filepath, []byte(section.content), 0644); err != nil {
			return fmt.Errorf("error writing file %s: %v", filename, err)
		}
		
		fmt.Println(style.FileCreated(filename))
		bar.Add(1)
	}

	fmt.Printf("\n%s\n", style.Summary(fmt.Sprintf("Successfully split markdown into %d files in %s", len(sections), splitsDir)))
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
	var preH2Content []string
	
	for _, line := range lines {
		if h2Regex.MatchString(line) {
			// If we have pre-H2 content, save it as the first section
			if currentSection == nil && len(preH2Content) > 0 {
				// Only save if there's non-whitespace content
				preContent := strings.TrimSpace(strings.Join(preH2Content, "\n"))
				if preContent != "" {
					preSection := section{title: "preamble", content: strings.Join(preH2Content, "\n")}
					sections = append(sections, preSection)
				}
			}
			
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
		} else {
			// Content before first H2 - collect it
			preH2Content = append(preH2Content, line)
		}
	}
	
	// Handle case where there's only pre-H2 content and no H2 sections
	if currentSection == nil && len(preH2Content) > 0 {
		// Only save if there's non-whitespace content
		preContent := strings.TrimSpace(strings.Join(preH2Content, "\n"))
		if preContent != "" {
			preSection := section{title: "preamble", content: strings.Join(preH2Content, "\n")}
			sections = append(sections, preSection)
		}
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
	filename = strings.ReplaceAll(filename, " ", "-")
	reg := regexp.MustCompile(`[^\w\-_.]`)
	sanitized := reg.ReplaceAllString(filename, "")
	sanitized = strings.ToLower(sanitized)
	
	// Limit length to avoid filesystem issues
	if len(sanitized) > 50 {
		sanitized = sanitized[:50]
	}
	
	return sanitized
}