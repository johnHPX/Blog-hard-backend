# imagem oficial da golang
FROM golang:1.18.3 as builder

# É importante que esses ARGs sejam definidos após a instrução FROM
ARG ACCESS_TOKEN_USR="johnHPX"
ARG ACCESS_TOKEN_PWD="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"

# Cria os arquivos de usuário e grupo que serão usados ​​na execução
# container para executar o processo como usuário sem privilégios.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group
# Crie um arquivo netrc usando as credenciais especificadas usando --build-arg
RUN printf "machine github.com\n\
    login ${ACCESS_TOKEN_USR}\n\
    password ${ACCESS_TOKEN_PWD}\n\
    \n\
    machine api.github.com\n\
    login ${ACCESS_TOKEN_USR}\n\
    password ${ACCESS_TOKEN_PWD}\n"\
    >> /root/.netrc
RUN chmod 600 /root/.netrc

# criando e configurando o diretorio
RUN mkdir -p /app
ADD . /app
WORKDIR /app

# configurando o go mod
RUN go mod download
RUN go mod verify

# gerando um binario da aplicação
RUN go build -o /server cmd/webapi/main.go

# imagem distroless, para redução do dockerfile
FROM gcr.io/distroless/base-debian10

# configurando o diretorio
WORKDIR /

# copiando o binario e a pasta configs para dentro do distroless
COPY --from=builder /server ./server
ADD configs ./configs

# espondo a porta da api
EXPOSE 40183

# # defindo que o usuario não root tera acesso a aplicação
USER nonroot:nonroot

# executando o binario
ENTRYPOINT ["./server"]