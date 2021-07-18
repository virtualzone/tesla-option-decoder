FROM golang:1.16-alpine AS server-builder
RUN apk --update add --no-cache git curl bash
RUN export GOBIN=$HOME/work/bin
WORKDIR /go/src/app
ADD . .
RUN ./mkjson.sh
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o main .
#RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main .

FROM arm64v8/alpine
COPY --from=server-builder /go/src/app/main /app/
COPY --from=server-builder /go/src/app/optioncodes.json /app/
ADD static/ /app/static
WORKDIR /app
EXPOSE 8080
CMD ["./main"]