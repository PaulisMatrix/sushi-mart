FROM postgres:11.21-alpine3.17

ARG PG_USER 
ARG PG_PASSWORD
ARG PGDB

ENV POSTGRES_DB=${PGDB}
ENV POSTGRES_USER=${PG_USER}
ENV POSTGRES_PASSWORD=${PG_PASSWORD}
ENV PGUSER=${POSTGRES_USER}

COPY init.sql triggers.sql /docker-entrypoint-initdb.d/
RUN sed -i -e "s/\${POSTGRES_DB}/${POSTGRES_DB}/g" -e "s/\${POSTGRES_USER}/${POSTGRES_USER}/g" /docker-entrypoint-initdb.d/init.sql

HEALTHCHECK --interval=10s --timeout=5s --retries=3 CMD [ "pg_isready -d $${POSTGRES_DB} -U $${POSTGRES_USER}"]