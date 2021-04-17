# Letsgo cli app

## Limitations
- The app assume that shell environment is zshrc
- Have to manually run the export commands or source zshrc file

## Command help
```
This cli app does the following
	1. Create the bin/, src/ and pkg/ directories if not present in current directory if flag is present
	2. Update GOPATH and PATH
	3. Update the zshrc file if flag is present

Usage:
  letsgo [flags]

Flags:
  -h, --help    help for letsgo
  -m, --mkdir   crate bin/, src/, pkg/ directories
  -z, --zshrc   update the ~/.zshrc file with new GOPATH and PATH
```

## To-Do
- do not create directories if present for e.g. src
- make this work for multiple shell environments
- create shell configuration file if not present
- avoid having to run manual export or source ~/.zshrc command