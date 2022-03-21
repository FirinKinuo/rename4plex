# Anime Replacer

Replacer for names and locations of anime files downloaded trackers. Usually the format of the names and the location of
the downloaded files does not fit the plex standard, for this you need to rename them and place them in the required
folder.

`anime-replacer` helps you quickly organize the correct directories and filenames. Together with the ability to run
scripts in torrent clients, automate this process

## Install

Fetch the [latest release](https://github.com/FirinKinuo/go-plex-anime/releases) for your platform

Then create required folder for config and copy config sample

```shell
sudo mkdir -p /etc/go-plex-anime
sudo cp ./configs/config.yaml.sample /etc/go-plex-anime/config.yaml
```

And create folder for logs

```shell
sudo mkdir -p /var/log/go-plex-anime
```

**or** install by Makefile instruction

```shell
sudo make install
```

## Building

For your platform

```shell
go mod download
CGO_ENABLED=0 go build -o anime-replacer cmd/main/main.go
```

**or** by Makefile

```shell
make build
```

**or** build for all tested platforms

```shell
make build-all
``` 

## Usage

`anime-replacer` requires only the path to the video or to the directory with the video

```shell
anime-replacer -p path/to/dir
```

```shell
anime-replacer -p path/to/file.avi
```

By default, `anime-replacer` will completely move the file to the correct directory, if you need to save the file
(For example, if you are a good person, and you are seed torrents), you need to add
the `--save-original` flag, which will copy file, which takes much longer and takes up extra disk space. To avoid this,
you need to add the `--try-symlink` flag, but this may not work on all platforms

```shell
anime-replacer --save-original --try-symlink -p path/to/file.avi 
```

### Config file

Default configuration file location `/etc/go-plex-anime/config.yaml`

An example configuration location `configs/config.yaml.sample`

#### Configuration Options:

| Option             | Description                                           | value                                   |
|--------------------|-------------------------------------------------------|-----------------------------------------|
| debug              | Enable debug mode                                     | bool                                    |
| log_level          | Logging level                                         | string `(degub, info, error, critical)` |
| log_path           | Path to the log file                                  | string                                  |
| regexps_anime_data | List of regular expressions to search for anime data* | list[string]                            |

Regular expressions must have groups:
- title
- season
- episode
- ext

`(?i)(?P<title>.*)_(?P<season>\d*)_\[(?P<episode>\d*)\]_\[anilibria_tv.*\.(?P<ext>.*)`
### Options

| Option                 | Description                                                      | Required |
|------------------------|------------------------------------------------------------------|----------|
| `-p` `--path <string>` | Path to Anime folder or file                                     | True     |
| `--save-original `     | Keep the original file (Slower)                                  | False    |
| `--try-symlink`        | Use symlink instead of copy (Faster, but may not work on all OS) | False    |

