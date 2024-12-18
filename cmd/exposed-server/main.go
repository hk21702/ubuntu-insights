package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port    int
	Verbose bool
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	rootCmd := getCommands()
	rootCmd.Execute()
}

func getCommands() *cobra.Command {
	var cfgFile string
	var port int
	var verbose bool

	rootCmd := &cobra.Command{
		Use:   "ubuntu-insights-exposed-server",
		Short: "Start the exposed server for Ubuntu Insights",
		Long:  "Start the exposed server for Ubuntu Insights, the internet exposed component which receives data from the client and does basic validation before forwarding it to the ingest server",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := initConfig(cfgFile)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to read config file")
			}

			zerolog.SetGlobalLevel(zerolog.InfoLevel)
			if config.Verbose || verbose {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
				log.Debug().Msg("Verbose logging enabled")
			}

			if port != 0 {
				config.Port = port
			}
			log.Info().Msgf("Starting server on port %d", config.Port)
			log.Debug().Msgf("Starting server with config: %+v", config)
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Optional config file, if not provided, defaults will be used. Values in the config file will be overridden by command line flags")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "port to run the server on")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")

	return rootCmd
}

func initConfig(filename string) (ServerConfig, error) {
	var config ServerConfig

	if filename == "" {
		return ServerConfig{Port: 8080, Verbose: false}, nil
	}

	viper.SetConfigFile(filename)
	err := viper.ReadInConfig()
	if err != nil {
		return ServerConfig{}, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return ServerConfig{}, err
	}

	return config, nil
}
