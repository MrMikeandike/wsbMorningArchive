package report

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/MrMikeandike/wsbMorningArchive/pkg/db"

	"github.com/jmoiron/sqlx"

	"github.com/turnage/graw/reddit"
)

var (
	// RedditUser is the user who posts the morning reports
	RedditUser = "/u/Teenoh"
	// SubReddit is the subreddit that the posts are made on
	SubReddit = "wallstreetbets"
)

// Report contains the important fields of the reddit post
type Report struct {
	FullID    string
	DateTime  time.Time
	Title     string
	rawString string
}

// ReportRow method will return a db.ReportRow struct from the Report struct.
func (r *Report) ReportRow(rowID int) *db.ReportRow {
	return &db.ReportRow{
		ReportID:     r.FullID,
		Title:        r.Title,
		RowID:        rowID,
		PostDateTime: r.DateTime,
		RawText:      r.rawString,
	}
}

// Page represents one HTTP GETs worth of reddit posts. The first id, last id, and post count can be used to send another request for the next page.
type Page struct {
	Reports   []Report
	FirstID   string
	LastID    string
	PostCount int
}

// GetPage fetches a single page of reports
func GetPage(bot reddit.Bot, params map[string]string) (Page, error) {
	harvest, err := bot.ListingWithParams(RedditUser, params)
	if err != nil {
		return Page{}, err
	}
	resultLEN := len(harvest.Posts)
	if resultLEN == 0 {
		return Page{
			PostCount: resultLEN,
		}, nil
	}
	var reports []Report
	for _, v := range harvest.Posts {
		// TODO: handle deleted posts labeld with v.Deleted
		if !isReport(v) {
			continue
		}
		r := Report{
			FullID:    v.Name,
			Title:     v.Title,
			rawString: v.SelfText,
			DateTime:  time.Unix(int64(v.CreatedUTC), 0),
		}
		reports = append(reports, r)
	}
	return Page{
		Reports:   reports,
		FirstID:   harvest.Posts[0].Name,
		LastID:    harvest.Posts[resultLEN-1].Name,
		PostCount: resultLEN,
	}, nil

}

func isReport(p *reddit.Post) bool {
	if p.Subreddit != SubReddit {
		return false
	}
	ok, _ := regexp.MatchString("[Mm][Oo][Rr][Nn][Ii][Nn][Gg] *?[Bb][Rr][Ii][Ee][Ff]", p.Title)
	if !ok {
		return false
	}

	return true
}

// GetPageBefore retrieves a single page of morning reports using a single post as a anchor
// Note that reddit has the before/after backwards in their api. The function says before, but it
// uses the "after" property of the api.
func GetPageBefore(bot reddit.Bot, before string, count string) (Page, error) {
	params := map[string]string{
		"raw_json": "1",
		"limit":    "100",
		"after":    before,
		"count":    count,
	}
	return GetPage(bot, params)
}

// GetPageAfter retrieves a single page of morning reports using a single post as a anchor
// Note that reddit has the before/after backwards in their api. The function says after, but it
// uses the "before" property of the api.
func GetPageAfter(bot reddit.Bot, after string, count string) (Page, error) {
	params := map[string]string{
		"raw_json": "1",
		"limit":    "100",
		"before":   after,
		"count":    count,
	}
	return GetPage(bot, params)
}

// GetAll retrieves all morning reports from reddit. Impliments GetPage.
func GetAll(bot reddit.Bot) ([]Report, error) {
	var (
		count      = 0
		before     = ""
		allReports []Report
		page       Page
		err        error
	)
	// safety net for now
	i := 0

	for {
		i++
		if i > 6 {
			return nil, fmt.Errorf("Looped forever for some reason")
		}
		page, err = GetPageBefore(bot, before, strconv.Itoa(count))
		if err != nil {
			return nil, err
		}
		if page.PostCount == 0 {
			return allReports, nil
		}
		allReports = append(allReports, page.Reports...)
		count = count + page.PostCount
		before = page.LastID
	}
}

// GetSpecific retrieves posts by their specific id
func GetSpecific(bot reddit.Bot, ids []string) (Page, error) {
	params := map[string]string{
		"raw_json": "1",
		"limit":    "100",
		//"before":   "",
		//"count":    "0",
	}

	harvest, err := bot.ListingWithParams("/by_id/"+strings.Join(ids, ","), params)
	if err != nil {
		return Page{}, err
	}
	resultLEN := len(harvest.Posts)
	if resultLEN == 0 {
		return Page{
			PostCount: resultLEN,
		}, nil
	}
	var reports []Report
	for _, v := range harvest.Posts {
		r := Report{
			FullID:    v.Name,
			Title:     v.Title,
			rawString: v.SelfText,
			DateTime:  time.Unix(int64(v.CreatedUTC), 0),
		}
		reports = append(reports, r)
	}
	return Page{
		Reports:   reports,
		FirstID:   harvest.Posts[0].Name,
		LastID:    harvest.Posts[resultLEN-1].Name,
		PostCount: resultLEN,
	}, nil

}

// ScanNew does stuff
func ScanNew(bot reddit.Bot, sqlDB *sqlx.DB) error {
	after, err := db.LastKnownID(sqlDB)
	if err != nil {
		return err
	}
	count := 0
	var reports []Report
	var page Page
	for {
		page, err = GetPageAfter(bot, after, strconv.Itoa(count))
		if err != nil {
			return err
		}
		if page.PostCount == 0 {
			break
		}
		reports = append(reports, page.Reports...)
		count = page.PostCount
		after = page.LastID
	}
	if len(reports) == 0 {
		return nil
	}

	return nil
}
