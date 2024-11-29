FROM ubuntu:20.04

#configs
ENV TZ=Europe/Minsk
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN printf '#!/bin/sh\nexit 0' > /usr/sbin/policy-rc.d
RUN export RUNLEVEL=1

#updates
RUN apt update -y --fix-missing
RUN apt install wget make iputils-ping curl -y
RUN apt install build-essential fakeroot dpkg-dev -y
RUN apt install libx11-dev -y 
RUN apt install libxcursor-dev -y
RUN apt install libxrandr-dev -y
RUN apt install libxinerama-dev -y
RUN apt install libxi-dev -y
RUN apt install libgl1-mesa-dev -y
RUN apt install pkg-config -y
RUN apt install libxxf86vm-dev -y

RUN echo 0
# install go 
COPY ./bin/go1.21.4.linux-amd64.tar.gz .

# install.sh 
RUN echo 0
COPY ./build/docker/scripts/install.sh .
RUN chmod +x install.sh 
RUN ./install.sh

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