# AWS CloudWatch Logs Streamer

![Build Status](https://img.shields.io/badge/build-passing-success)
![License](https://img.shields.io/badge/license-MIT-blue)

## Table of Contents
- [Overview](#overview)
- [Motivation](#motivation)
- [Requirements](#requirements)
- [Build](#build)
- [Usage](#usage)
- [Docker Support](#docker-support)
- [Contributing](#contributing)
- [License](#license)

## Overview

AWS CloudWatch Logs Streamer is a utility that reads event logs from AWS CloudWatch and writes them to `stdout`. This utility is useful for integrating CloudWatch logs with centralized logging solutions like Loki, Elasticsearch, or any system that can collect logs from `stdout`.

## Motivation

In environments where logs are aggregated to a centralized location, such as Loki or Elastic, it can be essential to also stream AWS CloudWatch logs. This tool facilitates that integration by printing CloudWatch logs to `stdout`, making it easy for log collectors like Promtail and Fluentd to send those logs to your centralized logging system.

## Requirements

- Make sure your AWS credentials and region are configured properly. For more information, see [AWS SDK Configuration Guide](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html).

## Build

To build the project, run:

```bash
make build
```

## Usage

After AWS credentials and the region are configured, run the streamer with the following command:

```bash
./build/aws-cloudwatch-logs-streamer -g <cloudwatch-group-name>
```

For a full list of command-line flags and options:

```bash
aws-cloudwatch-logs-streamer --help
```

### Command-Line Flags
```bash
Usage:
  aws-cloudwatch-logs-streamer [flags]

Flags:
      --config string            Config file (default is $HOME/.aws-cloudwatch-logs-streamer.yaml)
  -g, --groupname string         CloudWatch group name to stream from
  -h, --help                     Help for aws-cloudwatch-logs-streamer
  -i, --interval int             Check for new events every X milliseconds (default 1000)
      --squash                   Remove new lines from the original log line
  -s, --streamname stringArray   CloudWatch stream name to stream from (can be specified multiple times; streams from all available streams by default)
```

## Docker Support

To build and run the Docker image, use the following commands:

```bash
make docker

docker run \
  -e "AWS_REGION=<aws_region>" \
  -e "AWS_ACCESS_KEY_ID=<aws_key>" \
  -e "AWS_SECRET_ACCESS_KEY=<aws_secret>" \
  -e "LOG_STREAMER_GROUPNAME=/aws/rds/instance/sightd-production/postgresql" \
  -it aws-cloudwatch-logs-streamer:latest
```

## Contributing

Feel free to submit pull requests or create issues to improve the project.

## License

This project is licensed under the MIT License. See the LICENSE.md file for details.
