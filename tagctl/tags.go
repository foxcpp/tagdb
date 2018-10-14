package main

import (
	"fmt"
	"os"
)

var tagsUsageMsg = `Usage: tagctl tags [file]

Print tags on file (or all tags if no argument passed).

Options:
  -h, --help		Print this message.
`

type tagsOpts struct {
	file string
}

func parseTagsFlags(args []string) (*tagsOpts, error) {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return nil, ErrHelp
		}
	}
	if len(args) == 1 {
		canonPath, err := canonicalPath(args[0])
		if err != nil {
			return nil, err
		}
		return &tagsOpts{canonPath}, nil
	} else if len(args) == 0 {
		return &tagsOpts{}, nil
	} else {
		return nil, ErrHelp
	}
}

func tags(opts *tagsOpts) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	var tags []string
	if len(opts.file) != 0 {
		tags, err = db.TagsOnFile(opts.file)
	} else {
		tags, err = db.Tags()
	}
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
	opts, err := parseTagsFlags(args)
	if err != nil {
		if err == ErrHelp {
			fmt.Fprintln(os.Stderr, tagsUsageMsg)
			return 0
		} else {
			fmt.Fprintln(os.Stderr, "error:", err)
			return 1
		}
	}

	if err := tags(opts); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 2
	}
	return 0
}
