package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"go-plex-anime/internal/config"
	"go-plex-anime/internal/filesystem/anime"
	"go-plex-anime/internal/filesystem/movement"
	"go-plex-anime/pkg/search"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Opts struct {
	Path         string `short:"p" long:"path" description:"Path to Anime folder or file" required:"true"`
	SaveOriginal bool   `long:"save-original" description:"Keep the original file (Slower)" required:"false"`
	SymLink      bool   `long:"try-symlink" description:"Use symlink instead of copy (Faster, but may not work on all OS)" required:"false"`
}

var opts Opts

func replaceAnimeFile(animePath string, saveOriginal bool, useSymLink bool) (err error) {
	if _, err := os.Stat(animePath); err == nil && search.IsExistString(filepath.Ext(animePath), &anime.VideoFormats) {
		if animeFile, err := anime.InitFileAnime(opts.Path); err == nil {
			if _, err := movement.MoveAnimeToPlex(animeFile, saveOriginal, useSymLink); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func replaceAllAnimeInFolder(animeFolderPath string, saveOriginal bool, useSymLink bool) (err error) {
	if videoFiles, err := ioutil.ReadDir(animeFolderPath); err == nil {
		for _, videoFile := range videoFiles {
			if animeFile, err := anime.InitFileAnime(filepath.FromSlash(fmt.Sprintf("%s/%s", animeFolderPath, videoFile.Name()))); err == nil {
				if _, err := movement.MoveAnimeToPlex(animeFile, saveOriginal, useSymLink); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	} else {
		return err
	}

	return nil
}

func main() {
	config.InitLogger()
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Debugf("Unable to determine arguments, reason: %s", err)
		os.Exit(-1)
	}
	if animePath, err := os.Stat(opts.Path); err == nil {
		if animePath.IsDir() {
			if err := replaceAllAnimeInFolder(opts.Path, opts.SaveOriginal, opts.SymLink); err != nil {
				log.Error(err)
			}
		} else {
			if err := replaceAnimeFile(opts.Path, opts.SaveOriginal, opts.SymLink); err != nil {
				log.Error(err)
			}
		}
	} else {
		log.Fatal(err)
	}
}
