package database

import (
	"fmt"
	"log"
	"mangathrV2/ent"
	"time"
)

func (d *Driver) createManga(mangaID, title, source, mapping string) (*ent.Manga, error) {
	u, err := d.client.Manga.
		Create().
		SetMangaID(mangaID).
		SetTitle(title).
		SetSource(source).
		SetMapping(mapping).
		SetRegisteredOn(time.Now()).
		Save(d.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}

func (d *Driver) CreateManga(mangaID, title, source, mapping string) (*ent.Manga, error) {
	manga, err := d.QueryManga(mangaID)
	if err != nil {
		return d.createManga(mangaID, title, source, mapping)
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
		builder.SetTitle("")
	}

	err := builder.Exec(d.ctx)
	if err != nil {
		return fmt.Errorf("failed creating user: %w", err)
	}
	return nil
}
