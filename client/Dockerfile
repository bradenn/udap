FROM node:17-alpine3.14

WORKDIR /app

COPY . .

RUN npm ci

RUN npm run build

CMD [ "npm", "run", "serve"]
