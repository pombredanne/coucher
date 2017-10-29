package revclean

import (
	"log"
	"github.com/yoav/coucher/database"
)

func CleanRevs(db *database.Database, id string, rev string, dryrun bool) (string, error) {
	resolvedRev, conflicts, err := db.GetRev(id, true)
	if err != nil {
		return "", err
	}
	if conflicts == nil {
		return "", nil
		log.Printf("Document %s has no conflicts!\n", id)
	} else {
		log.Printf("Document %s has %d conflicts:\n%v\n", id, len(conflicts), conflicts)
	}
	if rev == "" {
		rev = resolvedRev
		log.Printf("No rev specified. Will preserve resolved rev: %s\n", rev)
	}

	log.Printf("Cleaning conflicst in %s, keeping rev %s\n", id, rev)
	err = db.CleanConflicts(id, rev, dryrun)
	if err != nil {
		return "", err
	}
	return rev, nil
}
