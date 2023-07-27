FROM golang:latest

WORKDIR $GOPATH/google-photo-viewer
COPY . $GOPATH/google-photo-viewer
RUN go build

EXPOSE 80
ENTRYPOINT ["./google-photo-viewer"]
CMD ["./google-photo-viewer"]
