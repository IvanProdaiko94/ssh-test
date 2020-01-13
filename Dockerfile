FROM golang:1.12.13 as builder

COPY . $GOPATH/src/github.com/IvanProdaiko94/ssh-test
WORKDIR $GOPATH/src/github.com/IvanProdaiko94/ssh-test

# install jq
RUN JQ_URL="https://circle-downloads.s3.amazonaws.com/circleci-images/cache/linux-amd64/jq-latest" \
  && curl --silent --show-error --location --fail --retry 3 --output /usr/bin/jq $JQ_URL \
  && chmod +x /usr/bin/jq \
  && jq --version

# install go-swagger
RUN download_url=$(curl -s https://api.github.com/repos/go-swagger/go-swagger/releases/latest | \
      		jq -r '.assets[] | select(.name | contains("'"$(uname | tr '[:upper:]' '[:lower:]')"'_amd64")) | .browser_download_url') && \
      	curl -o /usr/local/bin/swagger -L'#' "$download_url" && \
        chmod +x /usr/local/bin/swagger

# generate templates
RUN swagger validate ./swagger.yaml && swagger generate server -f ./swagger.yaml

# install dependencies
RUN GO111MODULE=on go mod download && GO111MODULE=on go mod vendor

# build an app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /go/bin/svc cmd/tic-tac-toe-server/main.go

## final build stage
FROM scratch

COPY --from=builder /go/bin/svc /svc
COPY ./cert /cert

EXPOSE 8080

CMD ["./svc"]