/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"aws-cloudwatch-logs-streamer/pkg/logstreamer"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var interval int
var groupName string
var streamNames []string
var squashLines bool

var rootCmd = &cobra.Command{
	Use:   "aws-cloudwatch-logs-streamer",
	Short: "Read cloudwatch log events and print to stdout",
	Long: `Read cloudwatch log events and print to stdout.
Then log collectors like Fluend, Promtail etc will collect and forward it to Elastic, Loki etc..

This is a tricky way to stream cloudwatch logs to a central logging infrastructure`,
	Run: func(cmd *cobra.Command, args []string) {
		stream()
	},
}

func stream() {
	streamer, err := logstreamer.New(groupName, interval)
	if err != nil {
		panic(err)
	}

	if len(streamNames) == 0 {
		streamNames, err = streamer.GetLogStreams()
		if err != nil {
			panic(err)
		}
	}

	events := make(chan string)
	for _, stream := range streamNames {
		go streamer.Read(stream, events)
	}

	for event := range events {
		if squashLines {
			event = strings.Replace(event, "\n", "", -1)
		}

		fmt.Println(event)
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aws-cloudwatch-logs-streamer.yaml)")

	rootCmd.Flags().StringVarP(&groupName, "groupname", "g", "", "Cloudwatch group name to stream from")
	rootCmd.Flags().StringArrayVarP(&streamNames, "streamname", "s", make([]string, 0), "Cloudwatch stream name to stream from, can be specified multiple times - by default stream from all available streams")
	rootCmd.Flags().IntVarP(&interval, "interval", "i", 1000, "Check for new events every X milliseconds")
	rootCmd.Flags().BoolVar(&squashLines, "squash", false, "Remove new lines from the original log line")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".aws-cloudwatch-logs-streamer")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
