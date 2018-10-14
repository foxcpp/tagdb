package main

import (
	"fmt"
	"os"
)

var addUsageMsg = `Usage: tagctl add <files> -t <tag> [-t <tag>] ...

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

type addOpts struct {
	files []string
	tags  []string
}

func parseAddFlags(args []string) (*addOpts, error) {
	res := new(addOpts)
	var err error
	res.files, res.tags, err = parseFilesTags(args)
	return res, err
}

func add(opts *addOpts) error {
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
			err := db.AddTag(tx, tag, file)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit()
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
