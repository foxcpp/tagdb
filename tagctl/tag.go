package main

import (
	"fmt"
	"os"
)

var tagUsageMsg = `Usage: tagctl tag <files> -t <tag> [-t <tag>] ...

Add all specified tags to files.

Options:
  -h, --help		Print this message.
  -t, --tag <tag>	Tag to tag; can be used multiple times

Examples:
  tagctl tag agreement.odt report.odt -t work
	Add 'work' tag to agreement.odt and report.odt files in current
	directory.

  tagctl tag report.odt -t todo
    Add 'todo' tag to report.odt file in current directory.
`

type tagOpts struct {
	files []string
	tags  []string
}

func parsetagFlags(args []string) (*tagOpts, error) {
	res := new(tagOpts)
	var err error
	res.files, res.tags, err = parseFilesTags(args)
	return res, err
}

func tag(opts *tagOpts) error {
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
			err := db.Tag(tx, t, file)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit()
}

func tagSubcmd() int {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, tagUsageMsg)
		return 1
	}
	opts, err := parsetagFlags(os.Args[2:])
	if err != nil {
		if err == ErrHelp {
			fmt.Fprintln(os.Stderr, tagUsageMsg)
			return 0
		} else {
			fmt.Fprintln(os.Stderr, "error:", err)
			return 1
		}
	}

	if err := tag(opts); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 2
	}
	return 0
}
