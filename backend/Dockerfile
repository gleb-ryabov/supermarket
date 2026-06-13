FROM golang:1.25.1 AS build
COPY . /backend
WORKDIR /backend
RUN go build -o server ./cmd/server
RUN go build -o migrate ./cmd/migrate

# --- API ---
FROM alpine:3.20 AS api
WORKDIR /backend 
COPY --from=build /backend/server .
CMD [ "./server" ]

# --- MIGRATIONS ---
FROM alpine:3.20 AS migrate
WORKDIR /backend
COPY --from=build /backend/migrate .
CMD [ "./migrate" ]