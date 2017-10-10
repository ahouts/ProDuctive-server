FROM golang:1.9.1
ADD . /go/src/github.com/ahouts/ProDuctive-server
RUN go-wrapper download github.com/ahouts/ProDuctive-server
RUN go-wrapper install github.com/ahouts/ProDuctive-server
ENTRYPOINT /go/bin/ProDuctive-server
EXPOSE 443
