ARG NODE_IMAGE=node:20.4-alpine

FROM ${NODE_IMAGE} AS base
RUN apk --no-cache add dumb-init
RUN mkdir -p /home/node/app && chown node:node home/node/app
WORKDIR /home/node/app
USER node
RUN mkdir tmp

FROM base AS dependencies
COPY --chown=node:node ./package*.json ./
RUN npm ci
COPY --chown=node:node . .
RUN ls -al

FROM dependencies AS build
RUN node ace build --production

FROM base AS production
COPY --chown=node:node ./package*.json ./
COPY --chown=node:node ./.env ./
RUN npm ci --production
COPY --chown=node:node --from=build /home/node/app/build .
EXPOSE 3333
CMD [ "dumb-init", "node", "server.js" ]
