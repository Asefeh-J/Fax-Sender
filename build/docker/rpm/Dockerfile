FROM centos:7

RUN yum update -y
RUN yum install wget make iputils-ping curl -y  

# install dev-tools
RUN yum install rpmdevtools -y 
RUN yum install rpm-build -y

# install go 
COPY ./bin/go1.21.4.linux-amd64.tar.gz .

# install.sh 
RUN echo 0
COPY ./build/docker/scripts/install.sh .
RUN chmod +x install.sh 
RUN ./install.sh

#remove 
RUN rm -rf /go1.21.4.linux-amd64.tar.gz 

# copy source files
RUN mkdir -p /app
VOLUME ["/app"]
WORKDIR /app

# set proxy for pacakges
ENV http_proxy http://10.10.10.101:64874
ENV https_proxy http://10.10.10.101:64874
ENV socks_proxy socks://10.10.10.101:8080

CMD ["/bin/bash"]