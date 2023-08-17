#!/usr/bin/env bash

set -xe

openssl genrsa -out conf/ssl/ca.key 2048

openssl req -new -x509 -days 365 -key conf/ssl/ca.key \
  -subj "/C=US/O=genie/OU=dev/CN=ca" \
  -out conf/ssl/ca.cert

openssl req -newkey rsa:2048 -nodes -keyout conf/ssl/server.key \
  -subj "/C=US/O=genie/OU=dev/CN=localhost" \
  -out conf/ssl/server.csr

openssl x509 -req \
  -extfile <(printf "subjectAltName=DNS:localhost") \
  -days 365 \
  -in conf/ssl/server.csr \
  -CA conf/ssl/ca.cert \
  -CAkey conf/ssl/ca.key \
  -CAcreateserial \
  -out conf/ssl/server.crt

rm conf/ssl/localhost.crt
cat conf/ssl/server.crt conf/ssl/ca.cert > conf/ssl/localhost.crt
