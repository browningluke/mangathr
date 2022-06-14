package database

import (
	"fmt"
	"mangathrV2/ent"
	"mangathrV2/internal/logging"
)

func (d *Driver) UpdateManga(mangaUpdate *ent.MangaUpdateOne) (*ent.Manga, error) {
	manga, err := mangaUpdate.Save(d.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating manga: %w", err)
	}
	logging.Debugln("manga was updated: ", manga)
	return manga, nil
}
