package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"projectx-server/game"
	"projectx-server/model/postgres"
)

const (
	version = "0.1.3"
)

var (
	flagVersion bool = false
)

func setupLogging() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func setupConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("configs/")

	viper.SetDefault("generator.octaves", 7)
	viper.SetDefault("generator.persistence", 2)
	viper.SetDefault("generator.scaleFactor", 1)

	return viper.ReadInConfig()
}

func setupFlags() {
	flag.BoolVar(&flagVersion, "version", false, "print version and exit")

	flag.Parse()
}

func parseDBConfig(config *viper.Viper) (result postgres.Config, err error) {
	if config == nil {
		return result, fmt.Errorf("missing [server] section in the configuration")
	}

	if err := config.Unmarshal(&result); err != nil {
		return result, fmt.Errorf("failed to parse [server] config section")
	}

	return
}

func parseServerConfig(config *viper.Viper) (result Config, err error) {
	if config == nil {
		return result, fmt.Errorf("missing [server] section in the configuration")
	}

	if err := config.Unmarshal(&result); err != nil {
		return result, fmt.Errorf("failed to parse [server] config section")
	}

	if len(result.RequestEndpoint) == 0 {
		return result, fmt.Errorf("you should set RequestEndpoint variable")
	}
	if len(result.EventEndpoint) == 0 {
		return result, fmt.Errorf("you should set EventEndpoint variable")
	}

	return
}

func parseGeneratorConfig(config *viper.Viper) (result game.TerrainGeneratorConfig, err error) {
	if config == nil {
		return result, fmt.Errorf("missing [generator] section in the configuration")
	}

	if err := config.Unmarshal(&result); err != nil {
		return result, fmt.Errorf("failed to parse [generator] config section")
	}

	return
}

func main() {
	setupFlags()

	if flagVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	setupLogging()
	log.Printf("ProjectX UDPServer v%s", version)

	if err := setupConfig(); err != nil {
		log.WithError(err).Fatal("Failed to init configuration")
	}

	serverConfig, err := parseServerConfig(viper.Sub("server"))
	if err != nil {
		log.WithError(err).Fatal("Failed to parse server config")
	}

	dbConfig, err := parseDBConfig(viper.Sub("db"))
	if err != nil {
		log.WithError(err).Fatal("Failed to parse db config")
	}

	generatorConfig, err := parseGeneratorConfig(viper.Sub("generator"))
	if err != nil {
		log.WithError(err).Fatal("Failed to parse generator config")
	}

	server, err := NewServer(serverConfig, dbConfig, generatorConfig)
	if err != nil {
		log.WithError(err).Fatalf("Failed to start server on %s", os.Args[1])
	}

	log.Fatal(server.Serve())
}
