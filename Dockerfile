FROM golang:1.15-buster
ENV GO111MODULE=on
WORKDIR /go/src/github.com/Nekroze/Gandalf
COPY go.* ./
RUN echo "Pulling go dependencies" \
	&& go mod download
COPY . .
RUN echo "Testing gandalf" \
	&& go test -v -cover -short github.com/Nekroze/Gandalf/...
