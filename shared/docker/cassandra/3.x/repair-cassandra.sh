#!/bin/sh

unset JVM_OPTS;

JAVA_OPTS="-Xmx512M -Xms512M" nodetool repair -pr -full
