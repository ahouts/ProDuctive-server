#!/bin/bash
if [ -f ./key.pem ]; then
    echo "./key.pem already exists..."
    exit -1
fi
if [ -f ./cert.pem ]; then
    echo "./cert.pem already exists..."
    exit -1
fi
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 3650 -nodes -subj "/C=US/ST=California/L=Santa Clara/O=None/CN=www.example.com"
chmod 700 key.pem cert.pem

