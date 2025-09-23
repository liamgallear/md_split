# md_split
Tool to split markdown files based on the H2 headings

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

## Usage

```bash
./md_split [markdown_file]
```

### Example

```bash
./md_split document.md
```

This will:
- Parse `document.md` for H2 headings (`## Section Title`)
- Create a `splits` directory in the same location as the input file
- Generate numbered files like `01-section-title.md`, `02-another-section.md`, etc.
- Each split file contains the H2 heading and all content until the next H2 heading

### Features

- ✅ Splits based on H2 headings (`## Section Title`)
- ✅ Numbers files sequentially (01, 02, 03, etc.)
- ✅ Sanitizes filenames (removes special characters, converts spaces to hyphens)
- ✅ Creates `splits` directory relative to the input file
- ✅ Preserves all markdown formatting in split files
- ✅ Includes subsections (H3, H4, etc.) with their parent H2 section
- ✅ Error handling for missing files and files without H2 headings

### Help

```bash
./md_split --help
```
