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
	interval := viper.GetInt("interval")
	streamNames := viper.GetStringSlice("streamname")
	squashLines := viper.GetBool("squash")
	groupName := viper.GetString("groupname")
	if groupName == "" {
		fmt.Fprintf(os.Stderr, "No group name provided\n")
		return
	}

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

	rootCmd.Flags().StringP("groupname", "g", "", "Cloudwatch group name to stream from")
	rootCmd.Flags().StringArrayP("streamname", "s", make([]string, 0), "Cloudwatch stream name to stream from, can be specified multiple times - by default stream from all available streams")
	rootCmd.Flags().IntP("interval", "i", 1000, "Check for new events every X milliseconds")
	rootCmd.Flags().Bool("squash", false, "Remove new lines from the original log line")
	viper.BindPFlags(rootCmd.Flags())
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

	viper.SetEnvPrefix("log_streamer")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
