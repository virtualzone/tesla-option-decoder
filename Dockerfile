FROM golang:1.16-alpine AS server-builder
RUN apk --update add --no-cache git curl bash
RUN export GOBIN=$HOME/work/bin
WORKDIR /go/src/app
ADD . .
RUN go get -d -v ./...
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:3.14
COPY --from=server-builder /go/src/app/main /app/
ADD static/ /app/static
ADD optioncodes.json /app/
WORKDIR /app
EXPOSE 8080
CMD ["./main"]