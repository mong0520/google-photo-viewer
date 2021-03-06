FROM golang:latest

WORKDIR $GOPATH/src/mong0520/google-photo-viewer
COPY . $GOPATH/src/mong0520/google-photo-viewer
RUN GO111MODULE=on go build

EXPOSE 8080
ENTRYPOINT ["./google-photo-viewer"]
CMD ["./google-photo-viewer"]
