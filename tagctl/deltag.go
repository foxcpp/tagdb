package main

import (
	"errors"
	"fmt"
	"os"
)

var deltagUsageMsg = `Usage: tagctl deltag <tag>

Remove tag from all files, effectively deleting it.

Options:
  -h, --help		Print this message.
`

type deltagOpts struct {
	tag string
}

func parseDeltagFlags(args []string) (opts *deltagOpts, err error) {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return nil, ErrHelp
		}
	}
	if len(args) != 1 {
		return nil, errors.New("exactly one argument is required: tag")
	}
	opts = new(deltagOpts)
	opts.tag = args[0]
	return
}

func deltag(opts *deltagOpts) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if err = db.ForgetTag(tx, opts.tag); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func deltagSubcmd() int {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, deltagUsageMsg)
		return 1
	}
	opts, err := parseDeltagFlags(os.Args[2:])
	if err != nil {
		if err == ErrHelp {
			fmt.Fprintln(os.Stderr, deltagUsageMsg)
			return 0
		} else {
			fmt.Fprintln(os.Stderr, "error:", err)
			return 1
		}
	}

	if err := deltag(opts); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 2
	}
	return 0
}
