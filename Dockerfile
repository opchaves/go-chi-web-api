FROM node:alpine AS dev_client
WORKDIR /app
COPY ./web/package.json ./web/package-lock.json .
RUN npm install
COPY ./web/ .
CMD ["npm", "run", "dev"]

FROM node:alpine AS build_client
WORKDIR /app
COPY ./web/package.json ./web/package-lock.json .
RUN npm ci
COPY ./web/ .
RUN npm run build
CMD ["npm", "run", "preview"]

FROM golang:1.22-alpine3.19 as base
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache make git
# Create appuser.
ENV USER=appuser
ENV UID=10001 

# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

# GOPATH is set in the go image
# GOPATH/bin is added to the path there.
WORKDIR $GOPATH/src/webapp
COPY go.mod go.sum .
RUN go mod download
RUN go mod verify
COPY . .
COPY --from=build_client /app/dist/* ./internal/web/build/
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o $GOPATH/bin/server ./cmd/server/main.go
EXPOSE 8080
CMD [ "make", "start" ]

FROM scratch as prod
# Import the user and group files from the builder.
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group
COPY --from=base /go/bin/server /bin/server
ARG APP_ENV=production
EXPOSE 8080
CMD ["/bin/server"]
