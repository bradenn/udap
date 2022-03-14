FROM node:17-alpine3.14

WORKDIR /app

COPY . .

RUN yarn install

RUN yarn build

CMD [ "yarn", "run", "serve"]
