FROM golang:1.25.1 AS build
COPY . /backend
WORKDIR /backend
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o migrate ./cmd/migrate

# --- API ---
FROM alpine:3.20 AS api
WORKDIR /backend 
COPY --from=build /backend/server .
EXPOSE 8080
CMD [ "./server" ]

# --- MIGRATIONS ---
FROM alpine:3.20 AS migrate
WORKDIR /backend
COPY --from=build /backend/migrate .
COPY --from=build /backend/migrations ./migrations
CMD [ "./migrate" ]