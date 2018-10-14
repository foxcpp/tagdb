package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/foxcpp/tagdb"
)

var generalUsageMsg = `Usage: tagctl <subcommand> [options...]

Control and query tags associated with files.

Options:
  -v, --version     Print version and exit
  -h, --help        Print usage help and exit

Subcommands:
  query	    List files by tag
  tag       Add tag(s) to file(s)
  untag     Remove tag(s) from file(s)
  retag	    Rename/merge tag
  tags	    List tags on file (or all tags)
  deltag    Delete tag.

Use 'tagctl subcmd -h' to get more detailed description and usage hints for
particular subcommand.

Notes:

Database to use can be specified using TAGDB environment variable. Default is
~/.tag.db.

All operations are atomic unless otherwise stated.
`

var (
	ErrHelp = errors.New("help requested")
)

func parseFilesTags(args []string) (files []string, tags []string, err error) {
	files = []string{}
	tags = []string{}

	nextIsTag := false
	for _, arg := range args {
		if arg[0] == '-' {
			if nextIsTag {
				// No argument for --tag:
				//	--tag --foo
				return nil, nil, errors.New("missing value for -t/--tag flag")
			}
			if arg == "-h" || arg == "--help" {
				return nil, nil, ErrHelp
			}
			if arg == "-t" || arg == "--tag" {
				nextIsTag = true
				continue
			}
		}
		if nextIsTag {
			tags = append(tags, arg)
			nextIsTag = false
		} else {
			canonPath, err := canonicalPath(arg)
			if err != nil {
				return nil, nil, err
			}
			files = append(files, canonPath)
		}
	}
	return files, tags, nil
}

func main() {
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, generalUsageMsg)
		exitCode = 1
		return
	}

	for _, arg := range os.Args[1:] {
		if arg == "-v" || arg == "--version" {
			fmt.Fprintln(os.Stderr, "tagdb", tagdb.Version)
			return
		}
	}

	subCmd := os.Args[1]
	if subCmd == "query" {
		exitCode = querySubcmd()
	} else if subCmd == "tag" {
		exitCode = tagSubcmd()
	} else if subCmd == "untag" {
		exitCode = untagSubcmd()
	} else if subCmd == "retag" {
		exitCode = retagSubcmd()
	} else if subCmd == "tags" {
		exitCode = tagsSubcmd()
	} else if subCmd == "deltag" {
		exitCode = deltagSubcmd()
	} else if subCmd == "help" || subCmd == "-h" || subCmd == "--help" {
		fmt.Fprintln(os.Stderr, generalUsageMsg)
	} else {
		fmt.Fprintln(os.Stderr, "error: unknown subcommand:", subCmd, "\n")
		exitCode = 1
	}
}
