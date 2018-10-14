package main

import (
	"fmt"
	"os"
)

var remUsageMsg = `Usage: tagctl rem <files> -t <tag> [-t <tag>] ...

Remove all specified tags from files. Operation is atomic, either all or none
tags removed from all files. Non-existing files and tags are silently ignored.

Options:
  -h, --help		Print this message.
  -t, --tag <tag>	Tag to add; can be used multiple times

Examples:
  tagctl rem report.odt -t todo
    Remove 'todo' tag from report.odt file in current directory.
`

type remOpts struct {
	files []string
	tags  []string
}

func parseRemFlags(args []string) (*remOpts, error) {
	res := new(remOpts)
	var err error
	res.files, res.tags, err = parseFilesTags(args)
	return res, err
}

func rem(opts *remOpts) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for _, file := range opts.files {
		for _, tag := range opts.tags {
			err := db.RemoveTag(tx, tag, file)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit()
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
