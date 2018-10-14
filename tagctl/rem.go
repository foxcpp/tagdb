package main

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
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
