FROM alpine:latest

COPY ./ ./
RUN chmod -R +x ./scripts

# install packages
RUN apk add --no-cache --update \
    postgresql-client \
    ca-certificates \
    git

# install golang
COPY --from=golang:1.17.3-alpine3.15 /usr/local/go/ /usr/local/go/
ENV GOPATH=$HOME/go
ENV PATH=/usr/local/go/bin:$GOPATH/bin:$PATH

# build server
RUN go build -v ./server.go

CMD ["./server"]
