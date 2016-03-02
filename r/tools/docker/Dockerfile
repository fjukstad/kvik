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
    xpdf 

RUN echo 'deb http://cran.rstudio.com/bin/linux/debian jessie-cran3/' >> /etc/apt/sources.list && \
    apt-key adv --keyserver keys.gnupg.net --recv-key 381BA480 && \ 
    apt-get update && \
    apt-get install -y \
        r-base

# WORKDIR /tmp
# RUN wget https://cran.r-project.org/src/base/R-3/R-3.1.0.tar.gz
# RUN tar -xzf R-3.1.0.tar.gz
# WORKDIR /tmp/R-3.1.0
# RUN ./configure -with-x=no
# WORKDIR /tmp/R-3.1.0
# RUN make 
# RUN cp bin/R /usr/bin/R

ADD deps.r /tmp/deps.r
WORKDIR /tmp/
RUN R --vanilla < deps.r

RUN go get github.com/gorilla/mux 
RUN mkdir -p /go/src/github.com/fjukstad

VOLUME /go/src/github.com/fjukstad/kvik 

EXPOSE 8181

WORKDIR /go/src/github.com/fjukstad/kvik/r/tools
#CMD go run server.go -packages "jsonlite"
