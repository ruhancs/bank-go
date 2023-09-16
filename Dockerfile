#construir image
#docker build -t bank:latest .

#permissao de leitura
#chmod +x wait-for.sh
#chmod +x start.sh

#rodar app
#docker run --name bank -p 8000:8000 bank:latest
#docker run --name bank -p 8000:8000 -e GIN_MODE=release bank:latest

# Build stage
FROM golang:1.21.0-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

#utilizar somente o executavel do app
# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY .env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration

EXPOSE 8000
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]