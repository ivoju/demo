package config

import (
	"log"

	pg "github.com/demo/pkg/v1.0/postgres"
	"github.com/kenshaw/envcfg"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Env    *envcfg.Envcfg
	PgConn *pg.PgConn
}

// setup sets up the environment.
func New() (*Config, *logrus.Logger, error) {
	// create config and logger
	env, err := envcfg.New()
	if err != nil {
		return nil, nil, err
	}

	logger := logrus.New()

	// force all writes to regular log to logger
	log.SetOutput(logger.Writer())
	log.SetFlags(0)

	// configure logging for environment
	tf := new(logrus.TextFormatter)
	//tf.ForceColors = logrus.IsTerminal()
	tf.FullTimestamp = true
	logger.Formatter = tf

	logger.Info("[CONFIG] Setup complete")

	logger.Infof("[POSTGRES] Setup initializing")

	pgConn, err := pg.New(env)
	if err != nil {
		return nil, nil, err
	}

	logger.Infof("[POSTGRES] Setup complete")

	return &Config{
		Env:    env,
		PgConn: pgConn,
	}, logger, nil
}
