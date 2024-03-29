FROM golang:1.22-alpine AS server-builder
RUN apk --update add --no-cache git curl bash
RUN export GOBIN=$HOME/work/bin
WORKDIR /go/src/app
ADD . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o main .

FROM alpine:3.19 AS optioncode-builder
RUN apk --update add --no-cache curl jq bash
WORKDIR /app
ADD mkjson.sh /app
RUN ./mkjson.sh
RUN jq -e . optioncodes.json >/dev/null

FROM gcr.io/distroless/static-debian12
COPY --from=server-builder /go/src/app/main /app/
COPY --from=optioncode-builder /app/optioncodes.json /app/
ADD static/ /app/static
WORKDIR /app
EXPOSE 8080
USER 65532:65532
CMD ["./main"]