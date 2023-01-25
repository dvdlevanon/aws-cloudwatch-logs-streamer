
# Brief

A small utility that read event logs from AWS Cloudwatch and write them to stdout.

# Motivation

Assuming you have a central location where all logs are stored, like Loki or Elastic. You want to be able to send Cloudwatch logs to the central location as well.

One of the ways to do so is using this tool which print the logs to stdout. And then logs collector like Promtail, Fluentd etc can collect and send the logs to the central location.

# Build

`make build`

# Run

This tool uses AWS SDK for accessing Cloudwatch, you should configure your environment before running this tool. For more information see [how to specifying credentials and region](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html)

Once aws credentials and region configured, run the streamer with this command:

`./build/aws-cloudwatch-logs-streamer -g <cloudwatch group name>`

Run with `--help` to get all command line flags and options

```
Usage:
  aws-cloudwatch-logs-streamer [flags]

Flags:
      --config string            config file (default is $HOME/.aws-cloudwatch-logs-streamer.yaml)
  -g, --groupname string         Cloudwatch group name to stream from
  -h, --help                     help for aws-cloudwatch-logs-streamer
  -i, --interval int             Check for new events every X milliseconds (default 1000)
      --squash                   Remove new lines from the original log line
  -s, --streamname stringArray   Cloudwatch stream name to stream from, can be specified multiple times - by default stream from all available streams
```

# Docker

Build and run a docker image using those commands:
```
docker make

docker run \
	-e "AWS_REGION=<aws_region>" \
	-e "AWS_ACCESS_KEY_ID=<aws_key>" \
	-e "AWS_SECRET_ACCESS_KEY=<aws_secret>" \
	-e "LOG_STREAMER_GROUPNAME=/aws/rds/instance/sightd-production/postgresql" \
	-it aws-cloudwatch-logs-streamer:latest
```
