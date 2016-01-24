# Use builds from launchpad
FROM ubuntu:latest

MAINTAINER Bjørn Fjukstad <bjorn@cs.uit.no>

# Install.
RUN \
  apt-get update && \
  apt-get -y dist-upgrade && \
  apt-get install -y software-properties-common && \
  add-apt-repository -y ppa:opencpu/opencpu-1.5 && \
  apt-get update && \
  apt-get install -y opencpu \
        opencpu-cache \ 
        rstudio-server \
        apache2-utils \ 
        python3 \
        libapache2-mod-proxy-html \
        libxml2-dev \
        iptables \ 
        git \
        vim 

# Install R packages
ADD deps.R /tmp/deps.R
WORKDIR /tmp/
RUN R --vanilla < deps.R

# Apache ports (without caching)
# EXPOSE 80
# EXPOSE 443
EXPOSE 8004

# directories etc. 
RUN mkdir -p /home/rstudio
RUN mkdir -p /home/data
RUN mkdir -p /usr/local/lib/R/site-library 

# fix permissions to the shared r pkg location so that we can read/write
# to it from the docker host 
RUN chmod -R g+w /usr/local/lib/R/site-library
RUN chmod -R g+w /home/rstudio

#RUN htpasswd -b -c /etc/http.passwd biopsy@mcgill van-mi-ka-al

# Add r server config file to only accept connections on localhost
# ADD rserver.conf /etc/rstudio/rserver.conf 

# VOLUME /home
VOLUME /tmp 

# Set up passwords and that
RUN htpasswd -b -c /etc/http.passwd user password
ADD rstudio.conf /etc/apache2/sites-enabled/rstudio.conf 
ADD rstudio.conf /etc/apache2/sites-available/rstudio.conf 
ADD opencpu.conf /etc/apache2/sites-enabled/opencpu.conf 
ADD opencpu.conf /etc/apache2/sites-available/opencpu.conf 

# install path for r packages 
ADD rsession.conf /etc/rstudio/rsession.conf
ADD server.conf /etc/opencpu/server.conf

RUN useradd -g 20 -u 1000 -p $(openssl passwd -1 rstudio) rstudio

# Create users that can use rstudio etc
# RUN useradd -s /bin/bash -d /home/bjorn -u 501 -g 20 bjorn
# RUN echo 'bjorn:test' | chpasswd
# RUN addgroup bjorn staff
# RUN mkdir /home/bjorn && chown bjorn:staff /home/bjorn
#
RUN a2enmod headers

RUN chown -R 1000:1001 /usr/lib/R/library
RUN chown -R 1000:1001 /usr/share/R

# Define default command.
CMD service opencpu restart &&  \ 
    service opencpu-cache restart && \
    rstudio-server restart && \
    tail -F /var/log/opencpu/apache_access.log
