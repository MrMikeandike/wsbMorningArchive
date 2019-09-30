package db

import (
	"time"

	"github.com/MrMikeandike/wsbMorningArchive/pkg/report"
	"github.com/jmoiron/sqlx"
)

// ReportRow is a sql row that represents a morning report
type ReportRow struct {
	RowID    int
	ReportID string
	Title    string
	PostDate time.Time
}

// LastKnownID retrieves the most recent morning report that was saved to the database, and returns the id
func LastKnownID(db *sqlx.DB) (string, error) {
	var lastKnown string
	db.Select(lastKnown, "SELECT ReportID from table ORDER BY `datetime` DESC LIMIT 1")
	return lastKnown, nil

}

// InsertNew inserts new Reports into the database
func InsertNew(db *sqlx.DB, reports []report.Report) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	insertString := "INSERT INTO table(row1, row2, row3, row4) VALUES (:rowid, :reportid, :title, :postdate)"
	for _, r := range reports {
		row := ReportRow{
			RowID:    1,
			ReportID: r.FullID,
			Title:    r.Title,
			PostDate: r.DateTime,
		}
		_, err = tx.Exec(insertString, row)
		if err != nil {
			return err
		}

	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
