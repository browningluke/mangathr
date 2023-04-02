package database

import (
	"fmt"
	"github.com/browningluke/mangathrV2/ent"
	"github.com/browningluke/mangathrV2/ent/chapter"
	"github.com/browningluke/mangathrV2/ent/manga"
	"github.com/browningluke/mangathrV2/internal/logging"
)

func (d *Driver) UpdateManga(mangaUpdate *ent.MangaUpdateOne) (*ent.Manga, error) {
	manga, err := mangaUpdate.Save(d.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating manga: %w", err)
	}
	logging.Debugln("manga was updated: ", manga)
	return manga, nil
}

func (d *Driver) DeleteManga(m *ent.Manga) error {
	_, err := d.client.
		Chapter.Delete().
		Where(chapter.HasMangaWith(manga.MangaID(m.MangaID))).
		Exec(d.ctx)
	if err != nil {
		return fmt.Errorf("failed deleting chapters for manga %w", err)
	}

	err = d.client.Manga.DeleteOne(m).Exec(d.ctx)
	if err != nil {
		return fmt.Errorf("failed deleting manga %w", err)
	}
	logging.Debugln("manga was deleted: ", m)
	return nil
}
