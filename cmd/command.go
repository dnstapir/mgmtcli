/*
 * Copyright (c) 2024 Johan Stenstam, johan.stenstam@internetstiftelsen.se
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dnstapir/tapir"
	// "github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

var PopCmd = &cobra.Command{
	Use:   "pop",
	Short: "Prefix command, only usable via sub-commands",
}

var PopMqttCmd = &cobra.Command{
	Use:   "mqtt",
	Short: "Prefix command, only usable via sub-commands",
}

var PopStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Instruct TAPIR-POP to stop",
	Run: func(cmd *cobra.Command, args []string) {
		resp := SendCommandCmd(tapir.CommandPost{
			Command: "stop",
		})
		if resp.Error {
			fmt.Printf("%s\n", resp.ErrorMsg)
		}

		fmt.Printf("%s\n", resp.Msg)
	},
}

func init() {
	// rootCmd.AddCommand(PopCmd)
	// PopCmd.AddCommand(PopStopCmd, PopMqttCmd)
}

func SendCommandCmd(data tapir.CommandPost) tapir.CommandResponse {
	_, buf, _ := api.RequestNG(http.MethodPost, "/command", data, true)

	var cr tapir.CommandResponse

	err := json.Unmarshal(buf, &cr)
	if err != nil {
		log.Fatalf("Error from json.Unmarshal: %v\n", err)
	}
	return cr
}
