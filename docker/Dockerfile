FROM ubuntu:latest

# update pkg things
RUN apt-get update 
RUN apt-get upgrade -y

# git python wget and all that jazz 
RUN apt-get install -y build-essential python python-dev  python-pip wget git mercurial pkg-config wget tar r-base mercurial

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
RUN cat /dev/zero | ssh-keygen -t rsa -C "" -q -N "" -f /root/.ssh/id_rsa

ENV GOPATH /root/go
ENV PATH $PATH:$GOPATH/bin

# Kvik dependencies
# First rollback gorest to specific changeset 
RUN go get code.google.com/p/gorest
WORKDIR /root/go/src/code.google.com/p/gorest
RUN hg update -r 3ba1ba16240d

RUN go get -tags zmq_4_x github.com/alecthomas/gozmq
RUN go get code.google.com/p/go.net/websocket

#RUN git clone https://github.com/fjukstad/gowebsocket.git /root/go/src/github.com/fjukstad/gowebsocket
#WORKDIR /root/go/src/github.com/fjukstad/gowebsocket
#RUN go install
RUN go get github.com/fjukstad/gowebsocket

#RUN git clone https://github.com/fjukstad/gographer.git /root/kvik/src/src/github.com/fjukstad/gographer
#WORKDIR /root/go/src/github.com/fjukstad/gographer
#RUN go install
RUN go get github.com/fjukstad/gographer

RUN go get github.com/fjukstad/gocache

RUN ls /root/go/src/github.com/fjukstad/
RUN go get github.com/fjukstad/rpcman

# now get kvik 
WORKDIR /root/go/src/github.com/fjukstad
RUN git clone -b 2.0 https://github.com/fjukstad/kvik.git

# Compile frontend, webapp and dataengine
WORKDIR  /root/go/src/github.com/fjukstad/kvik/webapp
RUN go install 

WORKDIR /root/go/src/github.com/fjukstad/kvik/frontend
RUN go install 

WORKDIR /root/go/src/github.com/fjukstad/kvik/dataengine
RUN go install 

# ALSO R things for python 
RUN pip install rpy2 singledispatch

RUN apt-get install -y default-jdk

RUN wget http://cran.r-project.org/src/contrib/rJava_0.9-6.tar.gz
RUN R CMD INSTALL rJava_0.9-6.tar.gz

RUN wget http://cran.r-project.org/src/contrib/xlsxjars_0.6.1.tar.gz
RUN R CMD INSTALL xlsxjars_0.6.1.tar.gz

RUN wget http://cran.r-project.org/src/contrib/xlsx_0.5.7.tar.gz
RUN R CMD INSTALL xlsx_0.5.7.tar.gz 

ENV LD_LIBRARY_PATH /usr/lib/jvm/java-7-openjdk-amd64/jre/lib/amd64/server:${LD_LIBRARY_PATH}

# expose port 8000 for Kvik 
EXPOSE 8000
EXPOSE 8080
EXPOSE 8888 

# supervisor for starting up kvik
RUN apt-get install -y supervisor
run mkdir -p /var/log/supervisor
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# New R script? 
# COPY data-engine.r /root/kvik/src/src/github.com/fjukstad/dataengine/
RUN ls /etc
# run supervisor! 
CMD ["/usr/bin/supervisord"]
