#!/bin/bash
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 3650 -nodes -subj "/C=US/ST=California/L=Santa Clara/O=None/CN=www.example.com"
