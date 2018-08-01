#!/usr/bin/env bash

openssl req \
  -newkey rsa:2048 -nodes -keyout ./shared/certs/snk.key \
  -x509 -sha256 -days 365 -out ./shared/certs/snk.crt