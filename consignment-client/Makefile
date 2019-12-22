build:
	GOOS=linux GOARCH=amd64 go build
	docker build -t consignment-client .

run:
	docker run -e MICRO_REGISTRY=mdns consignment-client
