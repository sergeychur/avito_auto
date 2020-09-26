FROM postgres:12.4

ENV POSTGRES_USER=docker \
POSTGRES_DB=docker \
POSTGRES_PASSWORD=docker

COPY ./sql/*.sql /docker-entrypoint-initdb.d/

EXPOSE 5432

CMD ["postgres"]