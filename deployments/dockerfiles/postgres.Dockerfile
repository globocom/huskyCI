FROM postgres:latest
ADD deployments/huskyci.sql /docker-entrypoint-initdb.d/