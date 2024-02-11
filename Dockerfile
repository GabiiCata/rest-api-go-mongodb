# Usar la imagen oficial de Go como imagen base
FROM golang:latest as builder

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar el archivo go.mod y go.sum para instalar las dependencias
# Asumiendo que tu proyecto usa módulos de Go
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copiar el resto del código fuente de la aplicación al contenedor
COPY . .

# Compilar la aplicación a un binario estático
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Usar una imagen de scratch para una imagen más pequeña y segura
FROM scratch

# Copiar el binario estático desde la etapa de construcción anterior
COPY --from=builder /app/main .

# Exponer el puerto que Gin escuchará
EXPOSE 3001

# Comando para ejecutar la aplicación
CMD ["./main"]
