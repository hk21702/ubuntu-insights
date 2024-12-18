package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func stringPtr(s string) *string {
	return &s
}

func TestGlobalLogLevel(t *testing.T) {
	testCases := []struct {
		name            string
		args            []string
		configContent   *string
		expectedVerbose bool
	}{
		{
			name:            "default values",
			args:            []string{},
			configContent:   nil,
			expectedVerbose: false,
		},
		{
			name:            "verbose flag",
			args:            []string{"--verbose"},
			configContent:   nil,
			expectedVerbose: true,
		},
		{
			name:            "verbose flag and config file",
			args:            []string{"--verbose"},
			configContent:   stringPtr("verbose: false"),
			expectedVerbose: true,
		},
		{
			name:            "config file",
			args:            []string{},
			configContent:   stringPtr("verbose: true"),
			expectedVerbose: true,
		},
		{
			name:            "short flag",
			args:            []string{"-v"},
			configContent:   nil,
			expectedVerbose: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.configContent != nil {
				file, err := os.CreateTemp("", "config-*.yaml")
				assert.NoError(t, err)
				defer os.Remove(file.Name())

				_, err = file.WriteString(*tc.configContent)
				assert.NoError(t, err)
				file.Close()

				tc.args = append(tc.args, "--config", file.Name())
			}

			rootCmd := getCommands()
			rootCmd.SetArgs(tc.args)

			err := rootCmd.Execute()
			assert.NoError(t, err)

			if tc.expectedVerbose {
				assert.Equal(t, zerolog.DebugLevel, zerolog.GlobalLevel())
			} else {
				assert.Equal(t, zerolog.InfoLevel, zerolog.GlobalLevel())
			}

		})
	}
}

func TestGetCommands(t *testing.T) {
	testCases := []struct {
		name            string
		args            []string
		expectedPort    int
		expectedVerbose bool
	}{
		{
			name:            "default values",
			args:            []string{},
			expectedPort:    8080,
			expectedVerbose: false,
		},
		{
			name:            "custom port",
			args:            []string{"--port", "9090"},
			expectedPort:    9090,
			expectedVerbose: false,
		},
		{
			name:            "verbose mode",
			args:            []string{"--verbose"},
			expectedPort:    8080,
			expectedVerbose: true,
		},
		{
			name:            "custom port and verbose mode",
			args:            []string{"--port", "9090", "--verbose"},
			expectedPort:    9090,
			expectedVerbose: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rootCmd := getCommands()
			rootCmd.SetArgs(tc.args)

			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)

			err := rootCmd.Execute()
			assert.NoError(t, err)

			port, err := rootCmd.PersistentFlags().GetInt("port")
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedPort, port)

			verbose, err := rootCmd.PersistentFlags().GetBool("verbose")
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedVerbose, verbose)
		})
	}
}

func TestInitConfig(t *testing.T) {
	testCases := []struct {
		name          string
		useDefault    bool
		configContent *string
		errorExpected bool
		exectedOutput ServerConfig
	}{
		{
			name:          "default, don't use config file",
			useDefault:    true,
			configContent: nil,
			errorExpected: false,
			exectedOutput: ServerConfig{Port: 8080, Verbose: false},
		},
		{
			name:          "valid config file - all values",
			useDefault:    false,
			configContent: stringPtr("port: 8081\nverbose: true"),
			errorExpected: false,
			exectedOutput: ServerConfig{Port: 8081, Verbose: true},
		},
		{
			name:          "valid config file - missing verbose",
			useDefault:    false,
			configContent: stringPtr("port: 8081"),
			errorExpected: false,
			exectedOutput: ServerConfig{Port: 8081, Verbose: false}},
		{
			name:          "invalid config file - bad value",
			useDefault:    false,
			configContent: stringPtr("port: abc"),
			errorExpected: true,
			exectedOutput: ServerConfig{}},
		{
			name:          "valid config file - empty file",
			useDefault:    false,
			configContent: stringPtr(""),
			errorExpected: false,
			exectedOutput: ServerConfig{}},
		{
			name:          "invalid config file - cant' find file",
			useDefault:    false,
			configContent: nil,
			errorExpected: true,
			exectedOutput: ServerConfig{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			configPath := ""
			if tc.configContent != nil {
				file, err := os.CreateTemp("", "config-*.yaml")
				assert.NoError(t, err)
				defer os.Remove(file.Name())

				_, err = file.WriteString(*tc.configContent)
				assert.NoError(t, err)
				file.Close()

				configPath = file.Name()
			} else if !tc.useDefault {
				configPath = "-invalid filepath-"
			}

			config, err := initConfig(configPath)

			if tc.errorExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.exectedOutput, config)

		})

	}
}
