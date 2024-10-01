/*
 * Copyright (c) 2024 Johan Stenstam, johan.stenstam@internetstiftelsen.se
 */

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dnstapir/tapir"
)

const DefaultMgmtCfgFile = "/etc/dnstapir/tapir-mgmt.yaml"

var imr string
var servername, certname, cfgFile, Prog string

var api *tapir.ApiClient

type Config struct {
	Services Services
}

type Services struct {
}

var rootCmd = &cobra.Command{
	Use:   "tapir-cli",
	Short: "CLI  utility used to interact with TAPIR-POP, i.e. the TAPIR Policy Processor",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var standalone bool

func init() {
	Prog = "tapir-mgmt"
	cobra.OnInitialize(RootInitConfig)
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&standalone, "standalone", "", true, "Run in standalone mode")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		fmt.Sprintf("config file (default is %s)", DefaultMgmtCfgFile))
	rootCmd.PersistentFlags().BoolVarP(&tapir.GlobalCF.Verbose, "verbose", "v", false, "Verbose mode")
	rootCmd.PersistentFlags().BoolVarP(&tapir.GlobalCF.Debug, "debug", "d", false, "Debugging output")
	rootCmd.PersistentFlags().BoolVarP(&tapir.GlobalCF.ShowHdr, "headers", "H", false, "Show column headers")
	rootCmd.PersistentFlags().BoolVarP(&tapir.GlobalCF.UseTLS, "tls", "", true, "Use TLS API connections")
}

var validate *validator.Validate

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("unable to unmarshal the config %v", err)
	}

	validate = validator.New()
	if err := validate.Struct(&config); err != nil {
		log.Fatalf("Missing required attributes in config %s:\n%v\n", viper.ConfigFileUsed(), err)
	}
}

// initConfig reads in config file and ENV variables if set.
func RootInitConfig() {
	if standalone {
		// In standalone mode we don't need to connect to TAPIR-SLOGGER, will not read any config files etc.
		fmt.Printf("Running in standalone mode; no config files, etc.\n")
		return
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		viper.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in. Terminate on all errors.
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config '%s': %v\n", cfgFile, err)
		}
	} else {
		switch Prog {
		case "tapir-mgmt":
			certname = "tapir-mgmt"
			servername = "tapir-slogger"
			viper.SetConfigFile(DefaultMgmtCfgFile)
			viper.AutomaticEnv() // read in environment variables that match

			// If a config file is found, read it in.
			if err := viper.ReadInConfig(); err != nil {
				fmt.Printf("Error reading config '%s': %v\n", viper.ConfigFileUsed(), err)
			}
			if tapir.GlobalCF.Debug {
				fmt.Println("Using config file:", viper.ConfigFileUsed())
			}

			viper.SetConfigFile(tapir.DefaultTapirCliCfgFile)
			viper.AutomaticEnv() // read in environment variables that match

			// If a config file is found, read it in.
			if err := viper.MergeInConfig(); err != nil {
				fmt.Printf("Error reading config '%s': %v\n", viper.ConfigFileUsed(), err)
			}
			if tapir.GlobalCF.Debug {
				fmt.Println("Using config file:", viper.ConfigFileUsed())
			}

		default:
			fmt.Printf("Unknown value for Prog: \"%s\"\n", Prog)
			os.Exit(1)
		}
	}

	validate = validator.New() // We need to initialize the Validate object in libcli!
	//	if err := validate.Struct(&config); err != nil {
	//		log.Fatalf("Missing required attributes in config %s:\n%v\n", viper.ConfigFileUsed(), err)
	//	}

	baseurl := viper.GetString("cli." + servername + ".url")
	if baseurl == "" {
		log.Fatalf("Error: missing config key: cli.%s.url", servername)
	}
	if tapir.GlobalCF.UseTLS {
		baseurl = viper.GetString("cli." + servername + ".tlsurl")
		if baseurl == "" {
			log.Fatalf("Error: missing config key: cli.%s.tlsurl", servername)
		}
	}

	var err error
	api = &tapir.ApiClient{
		BaseUrl:    baseurl,
		ApiKey:     viper.GetString("cli." + servername + ".apikey"),
		AuthMethod: "X-API-Key",
	}

	if tapir.GlobalCF.UseTLS { // default = true
		cd := viper.GetString("certs.certdir")
		if cd == "" {
			log.Fatalf("Error: missing config key: certs.certdir")
		}
		cert := cd + "/" + certname
		tlsConfig, err := tapir.NewClientConfig(viper.GetString("certs.cacertfile"),
			cert+".key", cert+".crt")
		if err != nil {
			log.Fatalf("Error: Could not set up TLS: %v", err)
		}
		// Must set this in the lab environment, as we don't know what addresses
		// put in the server cert IP SANs. Alternative would be to add a custom
		// VerifyConnection or VerifyPerrCertificate function in the TLS config?
		tlsConfig.InsecureSkipVerify = true
		err = api.SetupTLS(tlsConfig)
	} else {
		err = api.Setup()
	}

	if err != nil {
		log.Fatalf("Error setting up api client: %v", err)
	}

	tapir.GlobalCF.Api = api
}
