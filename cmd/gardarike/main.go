package main

import (
	"abbysoft/gardarike-online/db/postgres"
	"abbysoft/gardarike-online/logic"
	"abbysoft/gardarike-online/server"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	version = "0.1.4"
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

func parseServerConfig(config *viper.Viper) (result server.Config, err error) {
	if config == nil {
		return result, fmt.Errorf("missing [server] section in the configuration")
	}

	if err := config.Unmarshal(&result); err != nil {
		return result, fmt.Errorf("failed to parse [server] config section: %w", err)
	}

	if len(result.RequestEndpoint) == 0 {
		return result, fmt.Errorf("you should set RequestEndpoint variable")
	}
	if len(result.EventEndpoint) == 0 {
		return result, fmt.Errorf("you should set EventEndpoint variable")
	}

	return
}

func parseGeneratorConfig(config *viper.Viper) (result logic.TerrainGeneratorConfig, err error) {
	if config == nil {
		return result, fmt.Errorf("missing [generator] section in the configuration")
	}

	config.SetDefault("Octaves", 7)
	config.SetDefault("Persistence", 2)
	config.SetDefault("ScaleFactor", 1)

	if err := config.Unmarshal(&result); err != nil {
		return result, fmt.Errorf("failed to parse [generator] config section: %w", err)
	}

	return
}

func parseLogicConfig(config *viper.Viper) (result logic.Config, err error) {
	if config == nil {
		return result, fmt.Errorf("missing [logic] section in the configuration")
	}

	config.SetDefault("AFKTimeout", time.Minute*10)
	config.SetDefault("ChatMessageMaxLength", 200)

	if err := config.Unmarshal(&result); err != nil {
		return result, fmt.Errorf("failed to parse [logic] config section: %w", err)
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

	logicConfig, err := parseLogicConfig(viper.Sub("logic"))
	if err != nil {
		log.WithError(err).Fatal("Failed to parse logic config")
	}

	s, err := server.NewServer(serverConfig, logicConfig, dbConfig, generatorConfig)
	if err != nil {
		log.WithError(err).Fatalf("Failed to start server")
	}

	log.Fatal(s.Serve())
}
