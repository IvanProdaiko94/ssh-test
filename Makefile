export GO111MODULE=on

SHELL=/bin/bash
IMAGE_TAG := $(shell git rev-parse HEAD)
IMAGE_NAME := "host/company_name/project_name"
CERTIFICATE_DIR := "./cert"
GO=${GO_VARS} go

# https://goswagger.io/install.html installation guide
.PHONY: install_swagger
install_swagger:
	download_url=$(curl -s https://api.github.com/repos/go-swagger/go-swagger/releases/latest | \
  		jq -r '.assets[] | select(.name | contains("'"$(uname | tr '[:upper:]' '[:lower:]')"'_amd64")) | .browser_download_url')
	curl -o /usr/local/bin/swagger -L'#' "$download_url"
	chmod +x /usr/local/bin/swagger

.PHONY: generate
generate:
	swagger validate ./swagger.yaml
	swagger generate server -f ./swagger.yaml

.PHONY: run
run:
	$(GO) run cmd/tic-tac-toe-server/main.go \
		--tls-host localhost \
		--tls-port 8080 \
		--tls-certificate $(CERTIFICATE_DIR)/RootCA.crt \
		--tls-key $(CERTIFICATE_DIR)/RootCA.key

.PHONY: deps
deps:
	$(GO) mod download
	$(GO) mod vendor

.PHONY: unit_test
unit_test:
	$(GO) test -v -cover ./...

.PHONY: build
build:
	$(GO) build cmd/tic-tac-toe-server/main.go

.PHONY: dockerise
dockerise:
	docker build -t "${IMAGE_NAME}:${IMAGE_TAG}" -f Dockerfile .

.PHONY: generate_certificate
generate_certificate:
	mkdir -p $(CERTIFICATE_DIR)
	openssl req -x509 -nodes -new -sha256 -days 1024 \
		-newkey rsa:2048 -keyout $(CERTIFICATE_DIR)/RootCA.key \
		-out $(CERTIFICATE_DIR)/RootCA.pem -subj "/C=US/CN=Example-Root-CA"
	openssl x509 -outform pem -in $(CERTIFICATE_DIR)/RootCA.pem -out $(CERTIFICATE_DIR)/RootCA.crt