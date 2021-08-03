package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/goccy/go-yaml"
	"github.com/l3uddz/missarr/build"
	"github.com/l3uddz/missarr/sonarr"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type config struct {
	Sonarr sonarr.Config `yaml:"sonarr"`
}

var (
	// CLI
	cli struct {
		globals

		// flags
		Config    string `type:"path" default:"${config_file}" short:"c" env:"APP_CONFIG" help:"Config file path"`
		Log       string `type:"path" default:"${log_file}" short:"l" env:"APP_LOG" help:"Log file path"`
		Verbosity int    `type:"counter" default:"0" short:"v" env:"APP_VERBOSITY" help:"Log level verbosity"`

		// commands
		Sonarr SonarrCmd `cmd help:"Search sonarr for missing content"`
	}
)

type globals struct {
	Version versionFlag `name:"version" help:"Print version information and quit"`
	Update  updateFlag  `name:"update" help:"Update if newer version is available and quit"`
}

func main() {
	// cli
	ctx := kong.Parse(&cli,
		kong.Name("missarr"),
		kong.Description("Search for missing content in Sonarr/Radarr in a controlled fashion"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Summary: true,
			Compact: true,
		}),
		kong.Vars{
			"version":     fmt.Sprintf("%s (%s@%s)", build.Version, build.GitCommit, build.Timestamp),
			"config_file": filepath.Join(defaultConfigDirectory("missarr", "config.yml"), "config.yml"),
			"log_file":    filepath.Join(defaultConfigDirectory("missarr", "config.yml"), "activity.log"),
		},
	)

	if err := ctx.Validate(); err != nil {
		fmt.Println("Failed parsing cli:", err)
		return
	}

	// logger
	logger := log.Output(io.MultiWriter(zerolog.ConsoleWriter{
		TimeFormat: time.Stamp,
		Out:        os.Stderr,
		NoColor:    runtime.GOOS == "windows",
	}, zerolog.ConsoleWriter{
		TimeFormat: time.Stamp,
		Out: &lumberjack.Logger{
			Filename:   cli.Log,
			MaxSize:    5,
			MaxAge:     14,
			MaxBackups: 5,
		},
		NoColor: true,
	}))

	switch {
	case cli.Verbosity == 1:
		log.Logger = logger.Level(zerolog.DebugLevel)
	case cli.Verbosity > 1:
		log.Logger = logger.Level(zerolog.TraceLevel)
	default:
		log.Logger = logger.Level(zerolog.InfoLevel)
	}

	// config
	log.Trace().Msg("Initialising config")
	file, err := os.Open(cli.Config)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed opening config")
		return
	}
	defer file.Close()

	cfg := config{}
	decoder := yaml.NewDecoder(file, yaml.Strict())
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Error().Msg("Failed decoding configuration")
		log.Error().Msg(err.Error())
		return
	}

	// process cli
	if err := ctx.Run(&cfg); err != nil {
		log.Error().
			Err(err).
			Msg("Failed running command")
	}
}
