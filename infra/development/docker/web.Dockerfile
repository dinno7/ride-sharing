FROM node:24.13.0-alpine3.22

WORKDIR /app

COPY web/package*.json ./

RUN npm --registry="https://repo.hmirror.ir/npm" install

COPY web ./

RUN npm run build

EXPOSE 3000

CMD ["npm", "start"]
