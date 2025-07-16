# Git Diff Generator

A command-line tool that generates a comprehensive CSV file comparing files between two Git commits. This tool helps you analyze changes between commits by providing detailed information about file types, sizes, and modification dates.

## Features

- Compare files between any two Git commits
- Generate a structured CSV report with detailed file information
- Track file additions, modifications, and deletions
- Compare file sizes between commits
- Include custom remarks in the generated report
- Cross-platform support (Linux, macOS, Windows)

## Installation

### Option 1: Using Go Install (requires Go 1.19+)

```bash
go install github.com/j03hanafi/gitdiff@latest
```

### Option 2: Download Pre-built Binary

1. Visit the [Releases](https://github.com/j03hanafi/gitdiff/releases) page
2. Download the appropriate binary for your operating system:
   - Linux: `gitdiff_Linux_x86_64.tar.gz` or `gitdiff_Linux_arm64.tar.gz`
   - macOS: `gitdiff_Darwin_x86_64.tar.gz` or `gitdiff_Darwin_arm64.tar.gz`
   - Windows: `gitdiff_Windows_x86_64.zip` or `gitdiff_Windows_arm64.zip`
3. Extract the archive and place the binary in your PATH

## Usage

### Basic Usage

Navigate to your Git repository and run:

```bash
gitdiff --from <past-commit> --to <current-commit>
```

Example:

```bash
cd /path/to/your/repo
gitdiff --from abc123 --to def456
```

### Adding Remarks

You can add a custom remark to the generated CSV file:

```bash
gitdiff --from <past-commit> --to <current-commit> --remark "Sprint 42 changes"
```

### Command-line Options

| Option | Description | Required |
|--------|-------------|----------|
| `--from` | The base commit hash to compare from | Yes |
| `--to` | The target commit hash to compare to | Yes |
| `--remark` | Custom remark to include in the CSV output | No |
| `--help` | Display help information | No |

## Output Format

The tool generates a CSV file named `diff_<from-commit>_<to-commit>.csv` with the following structure:

1. Header row with commit identifiers (first 5 characters of each commit hash)
2. Column headers for both commits:
   - No (sequential number)
   - File Name (path to the file)
   - File Type (file extension)
   - Date Modified (commit date in "DD MMM YYYY" format)
   - File Size (in KB)
   - Remark (custom remark if provided)

Example CSV output:

```
,from,abc12,,,to,def45
No,File Name,File Type,Date Modified,File Size (KB),No,File Name,File Type,Date Modified,File Size (KB),Remark
1,src/main.go,.go,01 Jan 2023,12.50,1,src/main.go,.go,15 Jan 2023,13.20,Sprint 42 changes
2,README.md,.md,01 Jan 2023,1.20,2,README.md,.md,15 Jan 2023,2.40,Sprint 42 changes
```

## How It Works

1. The tool uses Git commands to identify files that changed between the specified commits
2. It checks out each commit to analyze file details (size, type, etc.)
3. It compiles all information into a structured CSV report
4. The CSV file is saved in the current directory

## Requirements

- Git must be installed and accessible in your PATH
- If building from source: Go 1.19 or later

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is open source and available under the [MIT License](LICENSE).

## Acknowledgments

- Built with [Go](https://golang.org/)
- Released using [GoReleaser](https://github.com/goreleaser/goreleaser)