# Etapa de construcción
FROM golang:1.21-alpine AS builder

# Instalar git, necesario para fetch de dependencias
RUN apk add --no-cache git

WORKDIR /usr/src/app

# Pre-copiar/cache go.mod y go.sum para pre-descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código de la aplicación y construir un binario estático
COPY . .
RUN CGO_ENABLED=0 go build -v -o /usr/local/bin/app -ldflags '-s -w' .

# Etapa final
FROM alpine:latest

# Copiar el binario estático desde la etapa de construcción
COPY --from=builder /usr/local/bin/app /usr/local/bin/app

# Set GIN_MODE to release to run Gin in production mode
ENV GIN_MODE=release

# Exponer el puerto que tu aplicación Gin usará
EXPOSE 5050

# Comando para ejecutar la aplicación
CMD ["app"]
