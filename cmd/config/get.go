/*
Copyright © 2023 Albert David Lewandowski a.lewandowski@magellanic.ai
*/
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

// getConfigCmd represents the getConfig command
var getConfigCmd = &cobra.Command{
	Use:   "get",
	Short: "Get config",
	RunE: func(cmd *cobra.Command, args []string) error {
		configId, err := cmd.PersistentFlags().GetString("config_id")
		if err != nil {
			return err
		}
		format, err := cmd.PersistentFlags().GetString("format")
		if err != nil {
			return err
		}
		apiKey, err := cmd.Flags().GetString("api_key")
		if err != nil {
			return err
		}
		outputPath, err := cmd.PersistentFlags().GetString("output")
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
	viper.BindEnv("config_id", "CONFIG_ID")
	getConfigCmd.PersistentFlags().StringP("format", "f", "json", "format (json, yaml or dotenv)")
	viper.BindEnv("format", "FORMAT")
	getConfigCmd.PersistentFlags().StringP("output", "o", "./output.json", "output file path")
	viper.BindEnv("output", "OUTPUT")
}

func validateGetConfigParams(configId, format string) error {
	if configId == "" {
		return errors.New("config_id not set")
	}
	if !slices.Contains([]string{"json", "yaml", "dotenv"}, format) {
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
	file, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	getConfigUrl := fmt.Sprintf("%s%s", apiUrl, apiPath)

	payload := &getConfigPayload{Id: configId, Key: apiKey, Format: format}
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
	_, err = io.Copy(file, getConfigRes.Body)
	if err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := getConfigRes.Body.Close(); err != nil {
		return err
	}
	return nil
}
