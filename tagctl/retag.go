package main

import (
	"errors"
	"fmt"
	"os"
)

var retagUsageMsg = `Usage: tagctl retag <oldname> <newname>

Rename/merge tag oldname into newname.

Options:
  -h, --help		Print this message.

Examples:
  tagctl query work		List all work-related files`

type renameOpts struct {
	from string
	to   string
}

func parseRetagFlags(args []string) (opts *renameOpts, err error) {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return nil, ErrHelp
		}
	}

	if len(args) < 1 {
		return nil, errors.New("missing required argument: oldname")
	}
	if len(args) < 2 {
		return nil, errors.New("missing required argument: newname")
	}

	opts = new(renameOpts)
	opts.from = args[0]
	opts.to = args[1]
	return
}

func retag(opts *renameOpts) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if err := db.RenameTag(tx, opts.from, opts.to); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func retagSubcmd() int {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, retagUsageMsg)
		return 1
	}
	opts, err := parseRetagFlags(os.Args[2:])
	if err != nil {
		if err == ErrHelp {
			fmt.Fprintln(os.Stderr, retagUsageMsg)
			return 0
		} else {
			fmt.Fprintln(os.Stderr, "error:", err)
			return 1
		}
	}

	if err := retag(opts); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 2
	}
	return 0
}
