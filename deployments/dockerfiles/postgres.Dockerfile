FROM postgres:latest
ADD deployments/mydb.sql /docker-entrypoint-initdb.d/