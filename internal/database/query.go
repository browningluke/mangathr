package database

import (
	"fmt"
	"github.com/browningluke/mangathr/v2/ent"
	"github.com/browningluke/mangathr/v2/ent/manga"
	"github.com/browningluke/mangathr/v2/ent/predicate"
	"github.com/browningluke/mangathr/v2/internal/logging"
)

func (d *Driver) queryManga(eager bool, ps ...predicate.Manga) (*ent.Manga, error) {
	mq := d.client.Manga.Query().
		Where(ps...)
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
	_, err := d.queryManga(false, manga.MangaID(mangaID))
	if err == nil {
		return true, nil
	}
	return false, err
}

func (d *Driver) CheckMangaExistenceByPredicate(ps ...predicate.Manga) (bool, error) {
	_, err := d.queryManga(false, ps...)
	if err == nil {
		return true, nil
	}
	return false, err
}

func (d *Driver) QueryMangaByID(mangaID, source string) (*ent.Manga, error) {
	return d.queryManga(true, manga.MangaID(mangaID), manga.Source(source))
}

func (d *Driver) QueryMangaByTitle(title, source string) (*ent.Manga, error) {
	return d.queryManga(true, manga.TitleEqualFold(title), manga.Source(source))
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
