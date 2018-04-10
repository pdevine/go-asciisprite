FROM golang:alpine3.7
WORKDIR /project
COPY whale.go .
COPY main.go .
RUN apk add --no-cache git
RUN go get github.com/nsf/termbox-go && go get github.com/pdevine/go-asciisprite
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o whale main.go whale.go

FROM scratch
COPY --from=0 /project/whale /whale
CMD ["/whale"]
