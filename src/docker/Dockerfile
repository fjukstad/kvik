FROM ubuntu:latest

# update pkg things
RUN apt-get update 
RUN apt-get upgrade -y

# git python wget and all that jazz 
RUN apt-get install -y build-essential python python-dev  python-pip wget git mercurial pkg-config wget tar r-base

RUN pip install numpy

# go
RUN wget https://storage.googleapis.com/golang/go1.3.1.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.3.1.linux-amd64.tar.gz
RUN rm go1.3.1.linux-amd64.tar.gz
ENV PATH $PATH:/usr/local/go/bin

# zeromq
RUN sudo apt-get install -y libtool autoconf automake
RUN wget http://download.zeromq.org/zeromq-4.0.4.tar.gz
RUN tar -xzf zeromq-4.0.4.tar.gz
RUN rm zeromq-4.0.4.tar.gz 
WORKDIR zeromq-4.0.4
RUN ./configure
RUN sudo make install 
RUN sudo ldconfig 

# back to home 
WORKDIR ~/

# zeromq bindings for python
RUN pip install pyzmq

# new ssh keys
RUN cat /dev/zero | ssh-keygen -q -N ""

# now get kvik 
WORKDIR /root
RUN git clone https://github.com/fjukstad/kvik.git
ENV GOPATH /root/kvik/src/
ENV PATH $PATH:$GOPATH/bin

# Kvik dependencies
RUN go get code.google.com/p/gorest
RUN go get -tags zmq_4_x github.com/alecthomas/gozmq
RUN go get github.com/fjukstad/gographer

RUN ls /root/kvik/src/

# Compile frontend, webapp and dataengine
WORKDIR  /root/kvik/src/src/github.com/fjukstad/webapp
RUN go install 

WORKDIR /root/kvik/src/src/github.com/fjukstad/frontend
RUN go install 

WORKDIR /root/kvik/src/src/github.com/fjukstad/dataengine
RUN go install 

# ALSO R and RPY2
RUN pip install rpy2
