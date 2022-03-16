NAME = mutate-ingress-for-imported-service
IMAGE_PREFIX = ghcr.io/alfianabdi
IMAGE_NAME = mutate-ingress-for-imported-service
IMAGE_VERSION = $$(git log --abbrev-commit --format=%h -s | head -n 1)

export GO111MODULE=on

app: deps
	go build -v -o $(NAME) cmd/main.go

deps:
	go get -v ./...

test: deps
	go test -v ./... -cover
	
docker:
	docker build --no-cache -t $(IMAGE_PREFIX)/$(IMAGE_NAME):$(IMAGE_VERSION) .

push:
	docker push $(IMAGE_PREFIX)/$(IMAGE_NAME):$(IMAGE_VERSION) 

.PHONY: docker push kind deploy reset