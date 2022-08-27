# imagem oficial da golang
FROM golang:1.18.3 as builder

# It is important that these ARG's are defined after the FROM statement
ARG ACCESS_TOKEN_USR="johnHPX"
ARG ACCESS_TOKEN_PWD="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"
# git is required to fetch go dependencies
# RUN apk add --no-cache ca-certificates git
# Create the user and group files that will be used in the running 
# container to run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group
# Create a netrc file using the credentials specified using --build-arg
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

# imagem scratch, para redução do dockerfile
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