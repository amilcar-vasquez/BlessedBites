FROM node:20-alpine AS build
WORKDIR /app

COPY frontend/package.json ./package.json
RUN npm install

COPY frontend .
RUN npm run build

FROM node:20-alpine
WORKDIR /app
COPY --from=build /app .
EXPOSE 3000
CMD ["npm", "run", "preview", "--", "--host", "0.0.0.0", "--port", "3000"]
