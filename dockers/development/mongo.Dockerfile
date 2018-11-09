FROM mongo:4.0.3

ADD dockers/development/mongo-init.js /docker-entrypoint-initdb.d/