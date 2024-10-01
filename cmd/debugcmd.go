/*
 * Copyright (c) 2024 Johan Stenstam, johan.stenstam@internetstiftelsen.se
 */
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dnstapir/tapir"
	"github.com/ryanuber/columnize"

	"github.com/invopop/jsonschema"
	"github.com/spf13/cobra"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Prefix command to various debug tools; do not use in production",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("debug called")

		var tm tapir.TagMask = 345
		fmt.Printf("%032b num tags: %d\n", tm, tm.NumTags())
	},
}

var debugMqttStatsCmd = &cobra.Command{
	Use:   "mqtt-stats",
	Short: "Return the MQTT stats counters from the MQTT Engine",
	Run: func(cmd *cobra.Command, args []string) {
		resp := SendDebugCmd(tapir.DebugPost{
			Command: "mqtt-stats",
		})
		if resp.Error {
			fmt.Printf("%s\n", resp.ErrorMsg)
		}
		if resp.Msg != "" {
			fmt.Printf("%s\n", resp.Msg)
		}

		var out = []string{"MQTT Topic|Msgs|Last MQTT Message|Time since last msg"}
		for topic, count := range resp.MqttStats.MsgCounters {
			t := resp.MqttStats.MsgTimeStamps[topic]
			out = append(out, fmt.Sprintf("%s|%d|%s|%v\n", topic, count, t.Format(timelayout), time.Since(t).Round(time.Second)))
		}
		fmt.Printf("%s\n", columnize.SimpleFormat(out))
	},
}

var popcomponent, popstatus string

var debugUpdatePopStatusCmd = &cobra.Command{
	Use:   "update-pop-status",
	Short: "Update the status of a TAPIR-POP component, to trigger a status update over MQTT",
	Run: func(cmd *cobra.Command, args []string) {
		switch popstatus {
		case "ok", "warn", "fail":
		default:
			fmt.Printf("Invalid status: %s\n", popstatus)
			os.Exit(1)
		}

		resp := SendDebugCmd(tapir.DebugPost{
			Command:   "send-status",
			Component: popcomponent,
			Status:    tapir.StringToStatus[popstatus],
		})
		if resp.Error {
			fmt.Printf("%s\n", resp.ErrorMsg)
		}
		if resp.Msg != "" {
			fmt.Printf("%s\n", resp.Msg)
		}
	},
}

var zonefile string
var Listname string

var debugGenerateSchemaCmd = &cobra.Command{
	Use:   "generate-schema",
	Short: "Experimental: Generate the JSON schema for the current data structures",
	Run: func(cmd *cobra.Command, args []string) {

		reflector := &jsonschema.Reflector{
			DoNotReference: true,
		}
		schema := reflector.Reflect(&tapir.WBGlist{}) // WBGlist is only used as a example.
		schemaJson, err := schema.MarshalJSON()
		if err != nil {
			fmt.Printf("Error marshalling schema: %v\n", err)
			os.Exit(1)
		}
		var prettyJSON bytes.Buffer

		// XXX: This doesn't work. It isn't necessary that the response is JSON.
		err = json.Indent(&prettyJSON, schemaJson, "", "  ")
		if err != nil {
			fmt.Printf("Error indenting schema: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%v\n", string(prettyJSON.Bytes()))
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)
	debugCmd.AddCommand(debugMqttStatsCmd)
	debugCmd.AddCommand(debugGenerateSchemaCmd, debugUpdatePopStatusCmd)

	debugUpdatePopStatusCmd.Flags().StringVarP(&popcomponent, "component", "c", "", "Component name")
	debugUpdatePopStatusCmd.Flags().StringVarP(&popstatus, "status", "s", "", "Component status (ok, warn, fail)")
}

type DebugResponse struct {
	Msg      string
	Data     interface{}
	Error    bool
	ErrorMsg string
}

func SendDebugCmd(data tapir.DebugPost) tapir.DebugResponse {
	_, buf, _ := api.RequestNG(http.MethodPost, "/debug", data, true)

	var dr tapir.DebugResponse

	var pretty bytes.Buffer
	err := json.Indent(&pretty, buf, "", "   ")
	if err != nil {
		fmt.Printf("JSON parse error: %v", err)
	}
	//	fmt.Printf("Received %d bytes of data: %v\n", len(buf), pretty.String())
	//	os.Exit(1)

	err = json.Unmarshal(buf, &dr)
	if err != nil {
		log.Fatalf("Error from json.Unmarshal: %v\n", err)
	}
	return dr
}
