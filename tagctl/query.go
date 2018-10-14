package main

import (
	"errors"
	"fmt"
)

type queryOpts struct {
	expr string
}

func parseQueryFlags(args []string) (opts *queryOpts, err error) {
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

	files, err := db.FileList(opts.expr)
	if err != nil {
		return err
	}

	for _, file := range files {
		fmt.Println(file)
	}
	return nil
}
