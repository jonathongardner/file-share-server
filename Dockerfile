FROM node:lts-alpine
# make the 'app' folder the current working directory
ENV ROOT /app
RUN mkdir -p $ROOT
WORKDIR $ROOT

COPY . .
RUN npm install

EXPOSE 3000
CMD [ "node", "app.js" ]
