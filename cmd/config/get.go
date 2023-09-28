// Package config /*
// Copyright © 2023 Magellanic <contact@magellanic.ai>
package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var getConfigCmd = &cobra.Command{
	Use:   "get",
	Short: "Get config",
	RunE: func(cmd *cobra.Command, args []string) error {
		configId, err := cmd.PersistentFlags().GetString("config_id")
		if err != nil {
			return err
		}
		format, err := cmd.PersistentFlags().GetString("config_format")
		if err != nil {
			return err
		}
		apiKey, err := cmd.Flags().GetString("api_key")
		if err != nil {
			return err
		}
		outputPath, err := cmd.PersistentFlags().GetString("config_output_file")
		if err != nil {
			return err
		}
		apiUrl, err := cmd.Flags().GetString("api_url")
		if err != nil {
			return err
		}
		if err = validateGetConfigParams(configId, format); err != nil {
			return err
		}
		return getConfig(apiKey, configId, format, outputPath, apiUrl)
	},
}

func init() {
	configCmd.AddCommand(getConfigCmd)
	getConfigCmd.PersistentFlags().StringP("config_id", "c", "", "config id")
	viper.BindEnv("mgl_config_id", "MGL_CONFIG_ID")
	getConfigCmd.PersistentFlags().StringP("config_format", "f", "json", "format (json, yaml, dotenv or react_env)")
	viper.BindEnv("mgl_format", "MGL_CONFIG_FORMAT")
	getConfigCmd.PersistentFlags().StringP("config_output_file", "o", "./output", "output file path")
	viper.BindEnv("mgl_output", "MGL_CONFIG_OUTPUT_FILE")
}

func validateGetConfigParams(configId, format string) error {
	if configId == "" {
		return errors.New("config_id not set")
	}
	if !slices.Contains([]string{"json", "yaml", "dotenv", "react_env"}, format) {
		return errors.New("format must be equal one of the following: json, yaml, dotenv")
	}
	return nil
}

type getConfigPayload struct {
	Id     string `json:"id"`
	Key    string `json:"key"`
	Format string `json:"format"`
}

func getConfig(apiKey, configId, format, outputPath, apiUrl string) error {

	getConfigUrl := fmt.Sprintf("%s%s", apiUrl, apiPath)

	var responseFormat string
	if format == "react_env" {
		responseFormat = "json"
	} else {
		responseFormat = format
	}
	payload := &getConfigPayload{Id: configId, Key: apiKey, Format: responseFormat}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(payloadBytes)

	getConfigReq, err := http.NewRequest(http.MethodPost, getConfigUrl, reader)
	if err != nil {
		return err
	}
	getConfigReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	defer client.CloseIdleConnections()

	getConfigRes, err := client.Do(getConfigReq)
	if err != nil {
		return err
	}
	if getConfigRes.StatusCode >= 400 {
		buf := new(strings.Builder)
		_, err = io.Copy(buf, getConfigRes.Body)
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("request error: %d\npayload: \n%s", getConfigRes.StatusCode, buf.String()))
	}
	file, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	if format == "react_env" {
		responseBytes, err := io.ReadAll(getConfigRes.Body)
		if err != nil {
			return err
		}
		var out bytes.Buffer
		err = json.Indent(&out, responseBytes, "", "  ")
		if err != nil {
			return err
		}
		str := out.String()
		str = "window.env = " + str
		_, err = file.WriteString(str)
		if err != nil {
			return err
		}
	} else {
		_, err = io.Copy(file, getConfigRes.Body)
		if err != nil {
			return err
		}
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := getConfigRes.Body.Close(); err != nil {
		return err
	}
	return nil
}
