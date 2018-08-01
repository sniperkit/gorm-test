#!/usr/bin/env bash

keytool -import -alias selfsigned -file ./shared/certs/snk.crt -keystore truststore.jks