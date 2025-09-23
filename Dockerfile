# --- ESTÁGIO 1: Build ---
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copia e baixa as dependências de forma otimizada
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código
COPY . .

# Compila um binário estático para a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main .

# --- ESTÁGIO 2: Final ---
FROM alpine:latest

# Define o diretório de trabalho
WORKDIR /app

# Cria um grupo e um usuário específicos para a aplicação
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copia o binário compilado e define o novo usuário como proprietário
COPY --from=builder --chown=appuser:appgroup /app/main .

# Muda para o usuário não-root
USER appuser

# Expõe a porta que a aplicação escuta
EXPOSE 8080

# Comando para iniciar a aplicação
CMD ["./main"]