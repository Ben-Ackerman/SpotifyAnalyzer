FROM node:11.12.0-alpine as build-stage
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
RUN npm install
ENV PATH /usr/src/app/node_modules/.bin:$PATH
COPY ./spotify-analyzer-frontend/package.json /usr/src/app/package.json
RUN npm install react-scripts@latest -g
COPY ./spotify-analyzer-frontend/ /usr/src/app
RUN CI=true npm test
RUN npm run build

FROM nginx:latest
EXPOSE 80
#COPY src /var/www
COPY --from=build-stage /usr/src/app/build /var/www
COPY nginx/nginx.conf /etc/nginx/nginx.conf