# Git Diff Generator

## How to install
Download one of the release file based on your OS or un (required min. Go 1.19):
```bash
go install github.com/j03hanafi/gitdiff@latest
```

## How to use
```bash
cd /path/to/your/repo
gitdiff --from <past-commit> --to <current-commit>
```

It will generate a diff file in the current directory.

## Help
```bash
gitdiff --help
```