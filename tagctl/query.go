package main

import (
	"errors"
	"fmt"
	"os"
)

var queryUsageMsg = `Usage: tagctl query <tag>

List files by tag, one absolute path per line.

Options:
  -h, --help		Print this message.

Examples:
  tagctl query work		List all work-related files
`

type queryOpts struct {
	expr string
}

func parseQueryFlags(args []string) (opts *queryOpts, err error) {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return nil, ErrHelp
		}
	}
	if len(args) != 1 {
		return nil, errors.New("exactly one argument is required: tag")
	}
	opts = new(queryOpts)
	opts.expr = args[0]
	return
}

func query(opts *queryOpts) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	files, err := db.TaggedFiles(opts.expr)
	if err != nil {
		return err
	}

	for _, file := range files {
		fmt.Println(file)
	}
	return nil
}

func querySubcmd() int {
	if len(os.Args) < 3 {
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
