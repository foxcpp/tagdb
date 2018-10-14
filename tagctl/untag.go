package main

import (
	"fmt"
	"os"
)

var untagUsageMsg = `Usage: tagctl untag <files> -t <tag> [-t <tag>] ...

Remove all specified tags from files. Non-existing files and tags are silently
ignored.

Options:
  -h, --help		Print this message.
  -t, --tag <tag>	Tag to add; can be used multiple times

Examples:
  tagctl untag report.odt -t todo
    Remove 'todo' tag from report.odt file in current directory.
`

type untagOpts struct {
	files []string
	tags  []string
}

func parseuntagFlags(args []string) (*untagOpts, error) {
	res := new(untagOpts)
	var err error
	res.files, res.tags, err = parseFilesTags(args)
	return res, err
}

func untag(opts *untagOpts) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for _, file := range opts.files {
		for _, t := range opts.tags {
			err := db.Untag(tx, t, file)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit()
}

func untagSubcmd() int {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, untagUsageMsg)
		return 1
	}
	opts, err := parseuntagFlags(os.Args[2:])
	if err != nil {
		if err == ErrHelp {
			fmt.Fprintln(os.Stderr, untagUsageMsg)
			return 0
		} else {
			fmt.Fprintln(os.Stderr, "error:", err)
			return 1
		}
	}

	if err := untag(opts); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 2
	}
	return 0
}
