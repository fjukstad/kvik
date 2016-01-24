# Use builds from launchpad
FROM ubuntu:latest

MAINTAINER Bjørn Fjukstad <bjorn@cs.uit.no>

# Install.
RUN \
  apt-get update && \
  apt-get -y dist-upgrade && \
  apt-get install -y \ 
        apache2 \ 
        apache2-utils \ 
        libapache2-mod-proxy-html \
        libxml2-dev \
        iptables \ 
        git \
        vim \
        wget \
        curl 

EXPOSE 80
EXPOSE 443
EXPOSE 8004

ADD load-balancer.conf /etc/apache2/sites-enabled/load-balancer.conf 
ADD load-balancer.conf /etc/apache2/sites-available/load-balancer.conf 


RUN a2enmod headers
RUN a2enmod proxy  
RUN a2enmod proxy_http  
RUN a2enmod proxy_balancer 
RUN a2enmod lbmethod_byrequests

RUN rm -rf /etc/apache2/sites-enabled/000-default.conf
RUN rm -rf /etc/apache2/sites-available/000-default.conf

# Define default command.
CMD service apache2 restart # && \ 
# tail -F /var/log/opencpu/apache_access.log
#CMD ["/bin/bash"]
