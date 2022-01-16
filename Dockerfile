FROM golang:1.17-alpine AS server-builder
RUN apk --update add --no-cache git curl bash
RUN export GOBIN=$HOME/work/bin
WORKDIR /go/src/app
ADD . .
RUN go get -d -v ./...
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:3.15 AS optioncode-builder
RUN apk --update add --no-cache curl jq bash
WORKDIR /app
ADD mkjson.sh /app
RUN ./mkjson.sh
RUN jq -e . optioncodes.json >/dev/null

FROM alpine:3.15
COPY --from=server-builder /go/src/app/main /app/
COPY --from=optioncode-builder /app/optioncodes.json /app/
ADD static/ /app/static
WORKDIR /app
EXPOSE 8080
CMD ["./main"]