FROM golang:1.14 AS builder
ENV PORT 8080
EXPOSE 8080

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

# FROM alpine:3.10
# COPY --from=builder /go/bin/regionalmap /usr/local/bin/
# RUN chmod +x /usr/local/bin/regionalmap

CMD ["regionalmap"]