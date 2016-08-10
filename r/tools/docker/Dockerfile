FROM golang:1.6.0
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y \
    build-essential \
    wget \
    gcc \
    gfortran \
    libreadline-dev \
    perl \
    vim \
    libxml2-dev \ 
    xpdf 

RUN echo 'deb http://cran.rstudio.com/bin/linux/debian jessie-cran3/' >> /etc/apt/sources.list && \
    apt-key adv --keyserver keys.gnupg.net --recv-key 381BA480 && \ 
    apt-get update && \
    apt-get install -y \
        r-base

RUN R -e 'install.packages("jsonlite", repos="http://cran.rstudio.org")'

RUN go get github.com/gorilla/mux 
RUN mkdir -p /go/src/github.com/fjukstad

# when the kvik repo is publically available we'll simply do a 
# go get ... to fetch the source code. 
# RUN go get github.com/fjukstad/kvik
VOLUME /go/src/github.com/fjukstad/kvik 

EXPOSE 80

WORKDIR /go/src/github.com/fjukstad/kvik/r/tools
CMD go run server.go -port=":80" 
