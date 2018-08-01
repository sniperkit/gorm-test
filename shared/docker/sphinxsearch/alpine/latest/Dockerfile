####################################
## Builder - image arguments
####################################
ARG RUNNER_ALPINE_VERSION=${RUNNER_ALPINE_VERSION:-"3.7"}
ARG RUNNER_IMAGE_NAME=${RUNNER_IMAGE_NAME:-"alpine:${RUNNER_ALPINE_VERSION}"}

####################################
## Runner
####################################
FROM ${RUNNER_IMAGE_NAME}

# WORKDIR /usr/bin
# COPY --from=build /go/bin .
RUN apk add --no-cache sphinx touch

RUN touch /etc/sphinx/sphinx.conf \
    && mkdir -p /var/run/sphinx/ \
    && mkdir -p /var/log/sphinx/ \
    && touch /var/log/sphinx/searchd.log

EXPOSE 9306 9312
VOLUME ["/var/lib/sphinx", "/etc/sphinx"]

COPY scripts/entrypoint.sh /opt/sphinxsearch/entrypoint.sh
RUN chmod +x /opt/sphinxsearch/entrypoint.sh

ENTRYPOINT ["/opt/sphinxsearch/entrypoint.sh"]
CMD ["searchd", "--nodetach"]