package main

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
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
