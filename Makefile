LOCAL_BIN:=$(CURDIR)/bin

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t cr.selcloud.ru/manzo/test-server:v0.0.1 .
	docker login -u token -p CRgAAAAAp6JuqinTqimeRiBLwtXQX5ZOSy-616Wn cr.selcloud.ru/manzo
	docker push cr.selcloud.ru/manzo/test-server:v0.0.1
