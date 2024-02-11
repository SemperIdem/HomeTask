FROM golang:1.20
WORKDIR /go/src/homeTask
COPY ./ ./
RUN make build
CMD ("./homeTask")