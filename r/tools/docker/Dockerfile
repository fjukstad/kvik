FROM golang
RUN apt-get update && apt-get install -y \
    build-essential \
    r-base 

RUN go get github.com/gorilla/mux 
RUN mkdir -p /go/src/github.com/fjukstad

VOLUME /go/src/github.com/fjukstad/kvik 

EXPOSE 8181

WORKDIR /go/src/github.com/fjukstad/kvik/r/tools
#CMD go run server.go -packages "jsonlite"
