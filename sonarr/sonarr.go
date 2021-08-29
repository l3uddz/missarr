package sonarr

import (
	"database/sql"
	"fmt"
	"github.com/l3uddz/missarr/logger"
	"github.com/l3uddz/missarr/migrate"
	"github.com/l3uddz/missarr/util"
	"github.com/rs/zerolog"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	URL     string `yaml:"url"`
	APIKey  string `yaml:"api_key"`
	Timeout int    `yaml:"timeout"`

	Verbosity string `yaml:"verbosity"`
}

type Client struct {
	apiURL     string
	apiHeaders map[string]string

	store *datastore
	http  *http.Client
	log   zerolog.Logger
}

func New(c *Config, db *sql.DB, mg *migrate.Migrator) (*Client, error) {
	l := logger.New(c.Verbosity).With().
		Str("pvr_type", "sonarr").
		Logger()

	// set config defaults
	if c.Timeout == 0 {
		c.Timeout = 90
	}

	// set api url
	apiURL := ""
	if strings.Contains(strings.ToLower(c.URL), "/api") {
		apiURL = c.URL
	} else {
		apiURL = util.JoinURL(c.URL, "api", "v3")
	}

	// set api headers
	apiHeaders := map[string]string{
		"X-Api-Key": c.APIKey,
	}

	// store
	store, err := newDatastore(db, mg)
	if err != nil {
		return nil, err
	}

	// create client
	cli := &Client{
		apiURL:     apiURL,
		apiHeaders: apiHeaders,

		store: store,
		http:  util.NewRetryableHttpClient(time.Duration(c.Timeout)*time.Second, nil, &l),
		log:   l,
	}

	// validate api access
	ss, err := cli.getSystemStatus()
	if err != nil {
		return nil, fmt.Errorf("validate api: %w", err)
	}

	cli.log.Info().
		Str("pvr_version", ss.Version).
		Msg("Initialised")

	return cli, nil
}
