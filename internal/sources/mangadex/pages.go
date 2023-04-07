package mangadex

import (
	"encoding/json"
	"fmt"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/manga"
	"github.com/browningluke/mangathrV2/internal/rester"
	"github.com/browningluke/mangathrV2/internal/utils"
	"math"
	"strconv"
	"time"
)

/*
	-- Pages --
*/

func (m *Scraper) addPagesToChapter(chapter *manga.Chapter) *logging.ScraperError {
	resInterface := rester.New().Get(
		fmt.Sprintf("%s/at-home/server/%s", APIROOT, chapter.ID),
		map[string]string{},
		[]rester.QueryParam{},
	).DoWithHelperFunc(4, "200ms", func(res rester.Response, err error) {
		logging.Errorln(res.StatusCode)

		if res.StatusCode == 429 {
			header := res.Headers

			nextRetryTimeInt := header["X-Ratelimit-Retry-After"][0]
			nextRetryTime, err := strconv.ParseInt(nextRetryTimeInt, 10, 64)
			now := time.Now().Unix()

			if err != nil {
				// Since we can't propagate an error (inside a helper func), we have to assume the worst case.
				nextRetryTime = now + 10 // Wait 10 seconds
			}

			timeDiff := nextRetryTime - now
			logging.Warningln(fmt.Sprintf("Time right now: %d", now))
			logging.Warningln(fmt.Sprintf("Sleeping for: %d", timeDiff))

			time.Sleep(time.Duration(timeDiff) * time.Second)
		}
	})

	jsonString := resInterface.(string)

	var chapterResp chapterResponse

	err := json.Unmarshal([]byte(jsonString), &chapterResp)
	if err != nil {
		return &logging.ScraperError{
			Error:   err,
			Message: "An error occurred while getting chapter pages",
			Code:    0,
		}
	}

	length := len(chapterResp.Chapter.Data)
	digits := int(math.Floor(math.Log10(float64(length))) + 1)

	logging.Debugln("Chapter pages: ", jsonString)

	addPages := func(slice []string, key string) {
		for i, p := range slice {
			chapter.AddPage(
				fmt.Sprintf("%s/%s/%s/%s",
					chapterResp.BaseUrl, key, chapterResp.Chapter.Hash, p),
				utils.PadString(fmt.Sprintf("%d", i+1), digits),
			)
		}
	}

	if config.DataSaver {
		addPages(chapterResp.Chapter.DataSaver, "data-saver")
	} else {
		addPages(chapterResp.Chapter.Data, "data")
	}

	return nil
}
