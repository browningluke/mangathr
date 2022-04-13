package database

import (
	"fmt"
	"log"
	"mangathrV2/ent"
)

func (d *Driver) CreateManga(mangaID, title, plugin string) (*ent.Manga, error) {
	u, err := d.client.Manga.
		Create().
		SetMangaID(mangaID).
		SetTitle(title).
		SetPlugin(plugin).
		Save(d.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}
