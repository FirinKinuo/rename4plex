package movement

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"go-plex-anime/internal/config"
	"go-plex-anime/internal/filesystem/anime"
	"go-plex-anime/pkg/file"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var cfg = config.GetConfig()

func MoveAnimeToPlex(animeFile *anime.FileAnime, saveOriginal bool, useSymLink bool) (animePath string, err error) {
	animeFolder := strings.ReplaceAll(path.Join(cfg.DirPlexAnimeLibrary, animeFile.Name), "_", " ")
	newAnimeFile := strings.ReplaceAll(path.Join(animeFolder, animeFile.GetAnimeFileName()), "_", " ")

	log.Infof("Rename anime file %s to %s", animeFile.Path, newAnimeFile)

	if _, err := os.Stat(animeFolder); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(animeFolder, os.ModePerm); err != nil {
			return "", err
		}
	}

	if _, err := os.Stat(newAnimeFile); errors.Is(err, os.ErrNotExist) {
		if saveOriginal {
			if useSymLink {
				log.Infof("Create SymLink attempt %s to %s", animeFile.Path, newAnimeFile)
				if err := os.Link(animeFile.Path, newAnimeFile); err != nil {
					return "", err
				}
			} else {
				log.Infof("Copy attempt %s to folder %s", animeFile.GetAnimeFileName(), animeFolder)

				tmpFile, err := ioutil.ReadFile(animeFile.Path)
				if err != nil {
					return "", err
				}

				if err := ioutil.WriteFile(newAnimeFile, tmpFile, 0644); err != nil {
					return "", err
				}
			}
		} else {
			log.Infof("Replace attempt %s to %s", animeFile.Path, newAnimeFile)
			if err := file.MoveFile(animeFile.Path, filepath.FromSlash(newAnimeFile)); err != nil {
				return "", err
			}
		}
	}

	log.Infof("Success move %s", newAnimeFile)
	return newAnimeFile, nil
}
