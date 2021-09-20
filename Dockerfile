FROM golang:1.16.3

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go clean --modcache
RUN go mod download
RUN go get github.com/pilu/fresh
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./main.go

CMD [ "fresh" ]