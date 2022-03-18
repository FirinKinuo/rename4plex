package anime

import (
	"errors"
	"fmt"
	"github.com/oriser/regroup"
	log "github.com/sirupsen/logrus"
	"go-plex-anime/internal/config"
	"go-plex-anime/pkg/search"
	"path/filepath"
	"strconv"
)

var cfg = config.GetConfig()
var VideoFormats = []string{".3gp", ".asf", ".avi", ".flv", ".m2ts", ".m4v", ".mkv", ".mov",
	".mp4", ".mts", ".ogg", ".vob", ".wmv", ".webm"}
var SeasonSpecials = []string{"OVA", "ONA", "OBA", "OAV"}

type FileAnime struct {
	Path      string
	Name      string
	Extension string
	Episode   int16
	Season    int8
}

func (f *FileAnime) GetAnimeFileName() string {
	fileName := fmt.Sprintf("%s.%s", f.Name, f.Extension)

	if f.Season != -1 && f.Episode != -1 {
		fileName = fmt.Sprintf("%s s%02de%02d.%s", f.Name, f.Season, f.Episode, f.Extension)
	}

	return fileName
}

func InitFileAnime(animePath string) (AnimeFile *FileAnime, err error) {
	log.Infof("Try init video %s as anime", animePath)
	if search.IsExistString(filepath.Ext(animePath), &VideoFormats) {
		for _, regex := range cfg.RegexpsAnimeData {
			log.Debugf("Try rexexp: %s", regex)
			match, err := regroup.MustCompile(regex).Groups(filepath.Base(animePath))
			if err == nil {
				log.Infof("Video %s init as Anime by regexp %s", animePath, regex)
				episode := -1
				season := 1
				if match["episode"] != "" {
					episode, err = strconv.Atoi(match["episode"])
					if err != nil {
						return nil, err
					}
				}

				if match["season"] != "" {
					if search.IsExistString(match["season"], &SeasonSpecials) {
						season = 0
					} else {
						season, err = strconv.Atoi(match["season"])
						if err != nil {
							return nil, err
						}
					}
				}

				return &FileAnime{
					Path:      filepath.FromSlash(animePath),
					Name:      match["title"],
					Extension: match["ext"],
					Episode:   int16(episode),
					Season:    int8(season),
				}, err
			}
		}
	} else {
		return nil, errors.New(fmt.Sprintf("File %s is not a video", animePath))
	}
	return nil, errors.New(fmt.Sprintf("File %s did not pass regexp", animePath))
}
