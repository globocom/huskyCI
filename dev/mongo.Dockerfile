FROM mongo:4.0.3

ADD dev/mongo-init.js /docker-entrypoint-initdb.d/
