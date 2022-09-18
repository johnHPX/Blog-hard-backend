# Arquitetura da Aplicação

Essa aplicação utiliza os conceitos de MVC(model,view,controller) para sua organização.

## CMD
    A Pasta CMD é reponsavel por guadar os arquivos com os pacote MAIN. é atravez dela que a aplicação será executada.
    atualmente existem dois arquivos MAIN, o migration e o main da aplicação.
    o migration é um arquivo responsavel pelas configurações do banco de dados, ppodendo criar ou deletar tabelas.
    main é o arquivo principal de fato, e nele é chamado todas as configuraçõe e funções da API e por onde é executada a Aplicação.

## CONFIGS
    Nessa pasta fica todas as configurações da API. é um arquivo yaml que contem: o nome do projeto, a porta, as configuraçõe do banco de dados e a conta de contato.

## DOCS 
    Aqui fica toda a documentação do projeto.

## MIGRATIONS
    aqui ficas as migrações do banco de dados, todas as informações das tabelas.

## INTERNAL
    aqui reside todo o código principal da aplicação. essa pasta contém 4 sub-pastas importantes.

### appl
    dentro da pasta appl, tem a pasta service, onde ela é responsavel pela lógica da aplicação, seria o "controller".

### domain
    dentro da pasta domain, tem a pasta models, onde ela é responsavel pelo os medelos das entidades, que seria um objeto que representa uma tabela do banco de dados.

### infra
    dentro da pasta infra, temos o repository e o utils. a pasta repository fica com a responsabilidade de lidar com o banco de dados, fazendo consultas e retornando querys, seria o "model". a pasta utils tem a responsabilidade de fornecer funções que seja uteis para o desenvolvimento da API. por exemplo a pasta responseAPI, que contém um objecto de erros da api junto com a função de exibir json, seja dados da aplicação ou erros.

### interf
    dentro da pasta interf, tem a resource e a route. resource é responsavel por exibir os dados de resquest e response da api, além de conter a função principal do endpoits, responsavel pela criação do endpoint, ela seria o "view". a pasta routes é responsavel por gerenciar as rotas da api, podendo definir os metdos, função principal, path e se é necessario o token de autenticação.

# como funciona a lógica da API

    resource -> service -> repository
    resource <- service <- repository

    é exatamente igual a lógica do padrão MVC.