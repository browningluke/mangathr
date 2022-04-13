package main

import (
	"context"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"mangathrV2/ent"
	"mangathrV2/ent/manga"
	"mangathrV2/internal/argparse"
	"mangathrV2/internal/commands/download"
	"mangathrV2/internal/config"
	"mangathrV2/internal/utils"
)

func main() {
	// Load config object, returns Config struct
	var c config.Config
	if err := c.Load("./examples/config.yml"); err != nil {
		utils.RaiseError(err)
	}
	fmt.Println(c)

	// Load argparse object, returns ArgParse struct
	var a argparse.Argparse
	if err := a.Parse(); err != nil {
		utils.RaiseError(err)
	}
	fmt.Println(a)

	switch a.Command {
	case "download":
		fmt.Println("Downloading", a.Download)
		download.Run(&a.Download, &c)
		break
	case "register":
		fmt.Println("Registering", a.Register)
	}

	// Merge Config & ArgParse (ArgParse priority) into ProgramOptions
	// Call (download|register|update|manage|config).go > run(ProgramOptions po) to start program execution
}

func register() {
	client, err := ent.Open("sqlite3", "file:examples/db.sqlite?cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer func(client *ent.Client) {
		err := client.Close()
		if err != nil {
			panic(err)
		}
	}(client)
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	ctx := context.Background()

	c, err := QueryManga(ctx, client)
	if err != nil {
		fmt.Println(err)
	}

	all, err := c.QueryChapters().All(ctx)
	if err != nil {
		return
	}

	fmt.Println(all)
	//_, err = CreateManga(ctx, client)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//_, err = QueryManga(ctx, client)
	//if err != nil {
	//	fmt.Println(err)
	//}
}

func CreateManga(ctx context.Context, client *ent.Client) (*ent.Manga, error) {
	c, err := client.Chapter.
		Create().
		SetChapterID("asfhjiou12f").
		SetNum("1").
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}

	u, err := client.Manga.
		Create().
		SetMangaID("123414afsfaksjfgiaq").
		SetTitle("komi-san").
		SetPlugin("mangadex").
		AddChapters(c).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}

func QueryManga(ctx context.Context, client *ent.Client) (*ent.Manga, error) {
	u, err := client.Manga.
		Query().
		Where(manga.MangaID("123414afsfaksjfgiaq")).
		// `Only` fails if no user found,
		// or more than 1 user returned.
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}
	log.Println("user returned: ", u)
	return u, nil
}
