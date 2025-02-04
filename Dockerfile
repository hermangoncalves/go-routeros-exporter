# Etapa 1: Build (Compila o binário Go)
FROM golang:1.23 AS builder

WORKDIR /app

# Copia os arquivos de dependência e baixa os módulos
COPY go.mod go.sum ./
RUN go mod download

# Copia o código-fonte (incluindo o main.go na raiz)
COPY . .

# Compila o binário com otimização
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main .

# Etapa 2: Criar a imagem final mínima usando `distroless`
FROM gcr.io/distroless/static:nonroot

WORKDIR /

# Copia o binário da fase anterior
COPY --from=builder /app/main /

# Define um usuário não-root para maior segurança
USER nonroot

# Expor a porta (caso necessário)
EXPOSE 8080

# Define o comando de execução
CMD ["/main"]