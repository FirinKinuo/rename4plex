package main

import (
	"errors"
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

func checkFileExist(filePath string) (err error) {
	if _, err = os.Stat(filePath); err != nil {
		return err
	}

	return nil
}

func checkIsFileVideo(filePath string) (err error) {
	if search.IsExistString(filepath.Ext(filePath), &anime.VideoFormats) != true {
		return errors.New(fmt.Sprintf("File %s - does not video!", filepath.Ext(filePath)))
	}
	return nil
}

func replaceAnimeFile(animePath string, saveOriginal bool, useSymLink bool) (err error) {
	if err = checkFileExist(animePath); err != nil {
		return err
	}

	if err = checkIsFileVideo(animePath); err != nil {
		return err
	}

	animeFile, err := anime.InitFileAnime(animePath)

	if err != nil {
		return err
	}

	_, err = movement.MoveAnimeToPlex(animeFile, saveOriginal, useSymLink)

	if err != nil {
		return err
	}

	return nil
}

func replaceAllAnimeInFolder(animeFolderPath string, saveOriginal bool, useSymLink bool) (err error) {
	videoFiles, err := ioutil.ReadDir(animeFolderPath)

	if err != nil {
		return nil
	}

	for _, videoFile := range videoFiles {
		err = replaceAnimeFile(filepath.Join(animeFolderPath, videoFile.Name()), saveOriginal, useSymLink)

		if err != nil {
			return err
		}
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

	animePath, err := os.Stat(opts.Path)

	if err != nil {
		log.Error(err)
	}

	if animePath.IsDir() {
		err = replaceAllAnimeInFolder(opts.Path, opts.SaveOriginal, opts.SymLink)
	} else {
		err = replaceAnimeFile(opts.Path, opts.SaveOriginal, opts.SymLink)
	}

	if err != nil {
		log.Error(err)
	}
}
