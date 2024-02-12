# Usar la imagen oficial de Go como imagen base
FROM golang:latest as builder

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar el archivo go.mod y go.sum para instalar las dependencias
# Asumiendo que tu proyecto usa módulos de Go
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Cambiar al directorio donde se encuentra el main.go
WORKDIR /app/cmd/server

# Copiar el resto del código fuente de la aplicación al contenedor
COPY . /app

# Compilar la aplicación a un binario estático
# Asegúrate de que esta ruta refleje la ubicación de tus fuentes para la compilación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Usar una imagen de scratch para una imagen más pequeña y segura
FROM scratch

# Copiar el binario estático desde la etapa de construcción anterior
COPY --from=builder /app/cmd/server/main .

# Copiar también el archivo .env si es necesario para la configuración en tiempo de ejecución
# Si tienes otros archivos de configuración, también deberás copiarlos aquí
COPY --from=builder /app/internal/env/docker.env ./.env

# Exponer el puerto que Gin escuchará
EXPOSE 3001

# Comando para ejecutar la aplicación
CMD ["./main"]
