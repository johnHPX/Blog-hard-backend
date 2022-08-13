# imagem oficial da golang
FROM golang:1.18.3-alpine3.16 as builder

# criando e configurando o diretorio
RUN mkdir -p /app
ADD . /app
WORKDIR /app

# configurando o go mod
RUN go mod download
RUN go mod verify

# gerando um binario da aplicação
RUN go build -o /server cmd/webapi/main.go

# imagem scratch, para redução do dockerfile
FROM scratch

# configurando o diretorio
WORKDIR /

# copiando o binario e a pasta configs para dentro do distroless
COPY --from=builder /server ./server
ADD configs ./configs

# espondo a porta da api
EXPOSE 40183

# defindo que o usuario não root tera acesso a aplicação
USER nonroot:nonroot

# executando o binario
ENTRYPOINT ["./server"]