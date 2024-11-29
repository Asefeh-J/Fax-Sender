#!/bin/bash

echo "going to install go"
rm -rf /usr/loca/go && tar -C /usr/local -xvf go1.21.4.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc

echo "DONE!"