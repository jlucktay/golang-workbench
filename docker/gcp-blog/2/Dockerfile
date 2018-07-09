# https://github.com/GoogleCloudPlatform/nodejs-docs-samples/tree/master/containerengine/hello-world

FROM node:alpine
WORKDIR /app
COPY package.json /app/package.json
RUN npm install --production
COPY server.js /app/server.js
EXPOSE 8080
CMD npm start
