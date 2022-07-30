package mangadex

import (
	"encoding/json"
	"fmt"
	"github.com/browningluke/mangathrV2/internal/downloader"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/rester"
	"github.com/browningluke/mangathrV2/internal/utils"
	"math"
	"strconv"
	"time"
)

/*
	-- Pages --
*/

func (m *Scraper) getChapterPages(id string) []downloader.Page {
	resInterface := rester.New().Get(
		fmt.Sprintf("%s/at-home/server/%s", APIROOT, id),
		map[string]string{},
		[]rester.QueryParam{},
	).DoWithHelperFunc(4, "200ms", func(res rester.Response, err error) {
		logging.Errorln(res.StatusCode)

		if res.StatusCode == 429 {
			header := res.Headers

			nextRetryTimeInt := header["X-Ratelimit-Retry-After"][0]
			nextRetryTime, err := strconv.ParseInt(nextRetryTimeInt, 10, 64)
			if err != nil {
				panic(err)
			}

			now := time.Now().Unix()
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
		panic(err)
	}

	length := len(chapterResp.Chapter.Data)
	digits := int(math.Floor(math.Log10(float64(length))) + 1)

	logging.Debugln("Chapter pages: ", jsonString)

	getPages := func(slice []string, key string) []downloader.Page {
		var pages []downloader.Page
		for i, chapter := range slice {
			pages = append(pages, downloader.Page{
				Url: fmt.Sprintf("%s/%s/%s/%s",
					chapterResp.BaseUrl, key, chapterResp.Chapter.Hash, chapter),
				Filename: fmt.Sprintf("%s%s",
					utils.PadString(fmt.Sprintf("%d", i+1), digits),
					utils.GetImageExtension(chapter)),
			})
		}
		return pages
	}

	var pages []downloader.Page

	if m.config.DataSaver {
		pages = getPages(chapterResp.Chapter.DataSaver, "data-saver")
	} else {
		pages = getPages(chapterResp.Chapter.Data, "data")
	}

	return pages
}
