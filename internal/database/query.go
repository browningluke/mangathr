package database

import (
	"fmt"
	"log"
	"mangathrV2/ent"
	"mangathrV2/ent/manga"
)

func (d *Driver) QueryManga(mangaID string) (*ent.Manga, error) {
	u, err := d.client.Manga.
		Query().
		Where(manga.MangaID(mangaID)).
		// `Only` fails if no user found,
		// or more than 1 user returned.
		Only(d.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}
	log.Println("user returned: ", u)
	return u, nil
}

func (d *Driver) QueryAllManga() ([]*ent.Manga, error) {
	u, err := d.client.Manga.
		Query().
		// `Only` fails if no user found,
		// or more than 1 user returned.
		All(d.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}
	log.Println("user returned: ", u)
	return u, nil
}
