# md_split
Tool to split and merge markdown files based on H2 headings

## Installation

1. Clone the repository:
```bash
git clone https://github.com/liamgallear/md_split.git
cd md_split
```

2. Build the tool:
```bash
go build -o md_split
```

### Pre-built Binaries

Alternatively, you can download pre-built binaries from the [Releases page](https://github.com/liamgallear/md_split/releases). Binaries are available for:
- Linux (AMD64, ARM64)
- macOS (AMD64, ARM64) 
- Windows (AMD64, ARM64)

## Usage

The tool provides two main commands: `split` and `merge`.

### Split Command

Split a markdown file into separate files based on H2 headings:

```bash
./md_split split [markdown_file]
```

#### Example

```bash
./md_split split document.md
```

This will:
- Parse `document.md` for H2 headings (`## Section Title`)
- Create a `splits` directory in the same location as the input file
- Generate numbered files like `01-section-title.md`, `02-another-section.md`, etc.
- Each split file contains the H2 heading and all content until the next H2 heading

### Merge Command

Merge split files back into a single markdown file:

```bash
./md_split merge [splits-directory] [output-file]
```

#### Example

```bash
./md_split merge ./splits merged_document.md
```

This will:
- Read all numbered markdown files from the `splits` directory
- Combine them in the correct order (01-, 02-, 03-, etc.)
- Output the merged content to `merged_document.md`

### Features

- ✅ **Split**: Splits based on H2 headings (`## Section Title`)
- ✅ **Merge**: Combines split files back into a single document
- ✅ Numbers files sequentially (01, 02, 03, etc.)
- ✅ Sanitizes filenames (removes special characters, converts spaces to hyphens)
- ✅ Creates `splits` directory relative to the input file
- ✅ Preserves all markdown formatting in split files
- ✅ Includes subsections (H3, H4, etc.) with their parent H2 section
- ✅ Error handling for missing files and files without H2 headings
- ✅ Maintains correct section order when merging

### Help

```bash
./md_split --help
./md_split split --help
./md_split merge --help
```
