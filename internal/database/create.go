package database

import (
	"fmt"
	"github.com/browningluke/mangathr/v2/ent"
	"github.com/browningluke/mangathr/v2/internal/logging"
	"time"
)

func (d *Driver) createManga(mangaID, title, source, mapping string, groups []string) (*ent.Manga, error) {
	u, err := d.client.Manga.
		Create().
		SetMangaID(mangaID).
		SetTitle(title).
		SetSource(source).
		SetMapping(mapping).
		SetRegisteredOn(time.Now()).
		SetFilteredGroups(groups).
		Save(d.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating manga: %w", err)
	}
	logging.Debugln("manga was created: ", u)
	return u, nil
}

func (d *Driver) CreateManga(mangaID, title, source, mapping string, groups []string) (*ent.Manga, error) {
	manga, err := d.QueryMangaByID(mangaID, source)
	if err != nil {
		return d.createManga(mangaID, title, source, mapping, groups)
	}
	return manga, nil
}

func (d *Driver) CreateChapter(chapterID, num, title string, manga *ent.Manga) error {
	builder := d.client.Chapter.
		Create().
		SetChapterID(chapterID).
		SetNum(num).
		SetManga(manga).
		SetRegisteredOn(time.Now())
	// TODO Add created time here

	if title != "" {
		builder.SetTitle(title)
	}

	err := builder.Exec(d.ctx)
	if err != nil {
		return fmt.Errorf("failed creating chapter: %w", err)
	}
	return nil
}
