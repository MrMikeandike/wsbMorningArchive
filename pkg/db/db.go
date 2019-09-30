package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// ReportRow is a sql row that represents a morning report
type ReportRow struct {
	RowID        int
	ReportID     string
	Title        string
	PostDateTime time.Time
	RawText      string
}

// Archive represents a row from the archive table of archived reports that are out of date since new edits
type Archive struct {
	ReportID        string
	RowID           int
	Title           string
	RawText         string
	PostDateTime    time.Time
	ArchiveDateTime time.Time
}

// ToArchive converts a ReportRow to a Archive
func (r *ReportRow) ToArchive() *Archive {
	return &Archive{
		ReportID:        r.ReportID,
		RowID:           r.RowID,
		Title:           r.Title,
		RawText:         r.RawText,
		PostDateTime:    r.PostDateTime,
		ArchiveDateTime: time.Now(),
	}
}

// LastKnownID retrieves the most recent morning report that was saved to the database, and returns the id
func LastKnownID(db *sqlx.DB) (string, error) {
	var lastKnown string
	db.Select(lastKnown, "SELECT ReportID from table ORDER BY `datetime` DESC LIMIT 1")
	return lastKnown, nil

}

// InsertNew inserts new Reports into the database
// func InsertNew(db *sqlx.DB, reports []report.Report) error {
// 	tx, err := db.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	insertString := "INSERT INTO table(row1, row2, row3, row4) VALUES (:rowid, :reportid, :title, :postdate)"
// 	stmt, err := tx.Prepare(insertString)
// 	if err != nil {
// 		return err
// 	}
// 	for _, r := range reports {
// 		row := ReportRow{
// 			RowID:        1,
// 			ReportID:     r.FullID,
// 			Title:        r.Title,
// 			PostDateTime: r.DateTime,
// 		}
// 		_, err = stmt.Exec(row)
// 		if err != nil {
// 			return err
// 		}

// 	}
// 	err = tx.Commit()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// SelectAllReports retrieves all current reports from the database.
func SelectAllReports(db *sqlx.DB) ([]ReportRow, error) {
	rows := []ReportRow{}
	queryString := "select rowid, reportid, title, postdate, rawtext from table"
	err := db.Select(rows, queryString)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// InsertArchive inserts a report into the archive table
func InsertArchive(db *sqlx.DB, r ReportRow) error {
	return nil
}
