# Build Operator
FROM golang:alpine
RUN mkdir -p $GOPATH/src/github.com/srleyva/aws-operator
WORKDIR /go/src/github.com/srleyva/aws-operator
ADD ./ ./ 
RUN  CGO_ENABLED=0 GOOS=linux go build main.go && cp main /main

# Inject PKI Server into container
FROM scratch
COPY --from=0 /main /bin/main
ENTRYPOINT ["main"]
