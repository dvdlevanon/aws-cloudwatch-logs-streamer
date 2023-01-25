
build:
	mkdir -p build && go build -o build/aws-cloudwatch-logs-streamer

docker:
	docker build -t aws-cloudwatch-logs-streamer .

clean:
	rm -rf build
