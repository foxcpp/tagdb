package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/foxcpp/tagdb"
)

var (
	generalUsageMsg = `Usage: tagctl <subcommand> [options...]

Control and query tags associated with files.

Options:
  -v, --version		Print version and exit
  -h, --help		Print usage help and exit

Subcommands:
  query		List files by tag
  add   	Add tag(s) to file(s)
  rem   	Remove tag(s) from file(s)

Use 'tagctl subcmd -h' to get more detailed description and usage hints for
particular subcommand.

Database to use can be specified using TAGDB environment variable. Default is
~/.tag.db.
`

	queryUsageMsg = `Usage: tagctl query <tag>

List files by tag, one absolute path per line.

Options:
  -h, --help		Print this message.

Examples:
  tagctl query work		List all work-related files`

	addUsageMsg = `Usage: tagctl add <files> -t <tag> [-t <tag>] ...

Add all specified tags to files. Operation is atomic, either all or none
tags added to all files.

Options:
  -h, --help		Print this message.
  -t, --tag <tag>	Tag to add; can be used multiple times

Examples:
  tagctl add agreement.odt report.odt -t work
	Add 'work' tag to agreement.odt and report.odt files in current
	directory.

  tagctl add report.odt -t todo
    Add 'todo' tag to report.odt file in current directory.
`

	remUsageMsg = `Usage: tagctl rem <files> -t <tag> [-t <tag>] ...

Remove all specified tags from files. Operation is atomic, either all or none
tags removed from all files. Non-existing files and tags are silently ignored.

Options:
  -h, --help		Print this message.
  -t, --tag <tag>	Tag to add; can be used multiple times

Examples:
  tagctl rem report.odt -t todo
    Remove 'todo' tag from report.odt file in current directory.
`
)

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

func querySubcmd() int {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, queryUsageMsg)
		return 1
	}
	opts, err := parseQueryFlags(os.Args[2:])
	if err != nil {
		if err == ErrHelp {
			fmt.Fprintln(os.Stderr, queryUsageMsg)
			return 0
		} else {
			fmt.Fprintln(os.Stderr, "error:", err)
			return 1
		}
	}

	if err := query(opts); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 2
	}
	return 0
}

func addSubcmd() int {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, addUsageMsg)
		return 1
	}
	opts, err := parseAddFlags(os.Args[2:])
	if err != nil {
		if err == ErrHelp {
			fmt.Fprintln(os.Stderr, addUsageMsg)
			return 0
		} else {
			fmt.Fprintln(os.Stderr, "error:", err)
			return 1
		}
	}

	if err := add(opts); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 2
	}
	return 0
}

func remSubcmd() int {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, remUsageMsg)
		return 1
	}
	opts, err := parseRemFlags(os.Args[2:])
	if err != nil {
		if err == ErrHelp {
			fmt.Fprintln(os.Stderr, remUsageMsg)
			return 0
		} else {
			fmt.Fprintln(os.Stderr, "error:", err)
			return 1
		}
	}

	if err := rem(opts); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 2
	}
	return 0
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
	} else if subCmd == "add" {
		exitCode = addSubcmd()
	} else if subCmd == "rem" {
		exitCode = remSubcmd()
	} else if subCmd == "help" || subCmd == "-h" || subCmd == "--help" {
		fmt.Fprintln(os.Stderr, generalUsageMsg)
	} else {
		fmt.Fprintln(os.Stderr, "error: unknown subcommand:", subCmd, "\n")
		exitCode = 1
	}
}
