package database

import (
	"fmt"
	"log"
	"mangathrV2/ent"
	"mangathrV2/ent/manga"
)

func (d *Driver) queryManga(mangaID string, eager bool) (*ent.Manga, error) {
	mq := d.client.Manga.Query().
		Where(manga.MangaID(mangaID))
	if eager {
		mq.WithChapters()
	}

	u, err := mq.Only(d.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}
	log.Println("user returned: ", u)
	return u, nil
}

func (d *Driver) CheckMangaExistence(mangaID string) (bool, error) {
	_, err := d.queryManga(mangaID, false)
	if err == nil {
		return true, nil
	}
	return false, err
}

func (d *Driver) QueryManga(mangaID string) (*ent.Manga, error) {
	return d.queryManga(mangaID, true)
}

func (d *Driver) QueryAllManga() ([]*ent.Manga, error) {
	u, err := d.client.Manga.
		Query().
		WithChapters().
		// `Only` fails if no user found,
		// or more than 1 user returned.
		All(d.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}
	log.Println("user returned: ", u)
	return u, nil
}
