FROM postgres:10.3-alpine

COPY dump/up.sql /docker-entrypoint-initdb.d/1.sql

CMD ["postgres"]
