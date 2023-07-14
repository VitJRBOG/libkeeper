FROM node:20
WORKDIR /app
COPY libkeeper-app/package.json libkeeper-app/main.js ./
COPY libkeeper-app/public/ ./public
COPY libkeeper-app/views/ ./views
RUN /bin/bash -c "npm install"