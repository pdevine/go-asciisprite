FROM golang:alpine3.7
WORKDIR /project
COPY pants.go .
RUN apk add --no-cache git
RUN go get github.com/nsf/termbox-go && go get github.com/pdevine/go-asciisprite
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o pants pants.go

FROM scratch
COPY --from=0 /project/pants /pants
CMD ["/pants"]
