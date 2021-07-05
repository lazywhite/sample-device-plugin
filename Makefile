IMG ?= "sample-device-plugin:v1"

docker-build:
	docker build . -t ${IMG}

build:
	CGO_ENABLED=0 go build -o sample-device-plugin

