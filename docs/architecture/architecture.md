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