#!/usr/bin/env bash
mkdir ~/ssl/
openssl genrsa -des3 -out ~/ssl/rootCA.key 2048
openssl req -x509 -new -nodes -key ~/ssl/rootCA.key -sha256 -days 1024 -out ~/ssl/rootCA.pem
