# Build Operator
FROM golang:alpine
RUN mkdir -p $GOPATH/src/github.com/srleyva/aws-operator
WORKDIR /go/src/github.com/srleyva/aws-operator
ADD ./ ./ 
RUN  CGO_ENABLED=0 GOOS=linux go build main.go && cp main /main
RUN apk update && apk add ca-certificates

# Inject Binary into container
FROM scratch
COPY --from=0 /etc/ssl/certs /etc/ssl/certs
COPY --from=0 /main /bin/main
ENTRYPOINT ["main"]
