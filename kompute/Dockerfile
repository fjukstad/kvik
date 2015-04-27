FROM rocker/rstudio

MAINTAINER Bjørn Fjukstad <bjorn@cs.uit.no>

# From the rstudio docker file at 
# https://registry.hub.docker.com/u/opencpu/rstudio/dockerfile/ 
RUN \
  useradd -ms /bin/bash builder && \
  apt-get update && \
  apt-get -y dist-upgrade && \
  apt-get install -y wget make devscripts apache2-prefork-dev apache2-mpm-prefork libapreq2-dev r-base r-base-dev libapparmor-dev libcurl4-openssl-dev libprotobuf-dev libprotoc-dev xvfb xauth xfonts-base curl libssl-dev pkg-config

USER builder

RUN \
  cd ~ && \
  wget https://github.com/jeffreyhorner/rapache/archive/v1.2.6.tar.gz && \
  tar xzf v1.2.6.tar.gz && \
  cd rapache-1.2.6 && \
  dpkg-buildpackage

RUN \
  cd ~ && \
  wget https://github.com/jeroenooms/opencpu-server/archive/v1.4.6.tar.gz && \
  tar xzf v1.4.6.tar.gz && \
  cd opencpu-server-1.4.6 && \
  dpkg-buildpackage

USER root

RUN \
  cd /home/builder && \
  dpkg -i libapache2-mod-r-base_*.deb && \
  dpkg -i opencpu-lib_*.deb && \
  dpkg -i opencpu-server_*.deb

RUN \
  apt-get purge -y devscripts apache2-prefork-dev libapreq2-dev \ 
  libapparmor-dev libprotobuf-dev libprotoc-dev xvfb xauth pkg-config && \
  apt-get autoremove --purge -y && \
  deluser builder --remove-home

# Apache ports (without caching)
EXPOSE 80
EXPOSE 443
EXPOSE 8004

# Create users that can use rstudio etc
# RUN useradd -s /bin/bash -d /home/bjorn -u 501 -g 20 bjorn
# RUN echo 'bjorn:test' | chpasswd
# RUN addgroup bjorn staff
# RUN mkdir /home/bjorn && chown bjorn:staff /home/bjorn

# Define default command.
CMD service opencpu restart && rstudio-server restart && \
    tail -F /var/log/opencpu/apache_access.log
