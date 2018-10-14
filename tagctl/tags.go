package main

import (
	"fmt"
	"os"
)

var tagsUsageMsg = `Usage: tagctl tags

Print all known tags, one name per line.

Options:
  -h, --help		Print this message.
`

func parseTagsFlags(args []string) (err error) {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return ErrHelp
		}
	}
	return nil
}

func tags() error {
	db, err := getDB()
	if err != nil {
		return err
	}

	tags, err := db.Tags()
	if err != nil {
		return err
	}

	for _, tag := range tags {
		fmt.Println(tag)
	}
	return nil
}

func tagsSubcmd() int {
	args := []string{}
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}
	err := parseTagsFlags(args)
	if err != nil {
		if err == ErrHelp {
			fmt.Fprintln(os.Stderr, tagsUsageMsg)
			return 0
		} else {
			fmt.Fprintln(os.Stderr, "error:", err)
			return 1
		}
	}

	if err := tags(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 2
	}
	return 0
}
