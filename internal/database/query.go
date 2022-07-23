package database

import (
	"fmt"
	"github.com/browningluke/mangathrV2/ent"
	"github.com/browningluke/mangathrV2/ent/manga"
	"github.com/browningluke/mangathrV2/internal/logging"
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
	logging.Debugln("manga returned: ", u)
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
		All(d.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}
	logging.Debugln("manga returned: ", u)
	return u, nil
}
