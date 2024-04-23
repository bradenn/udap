FROM node:21.5-bookworm

WORKDIR /app

COPY . .

RUN yarn install

RUN yarn build

CMD [ "yarn", "run", "preview", "--host", "--port", "5045"]
