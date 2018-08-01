FROM nginx:1.15.1-alpine

RUN apk add --no-cache --no-progress inotify-tools nano bash tree

# COPY config/nginx.conf /etc/nginx/nginx.conf
COPY scripts/entrypoint.sh /opt/nginx/entrypoint.sh

EXPOSE 80 443

ENTRYPOINT ["/opt/nginx/entrypoint.sh"]
CMD ["watch"]