# BlogHard - backend
Um sistema de gerenciamento de conteúdo(CMS) feito com GO, postgress e docker.

Toda documentação necessaria está na pasta "docs" e nesse link: https://drive.google.com/drive/u/0/folders/1z8u4znIMH_JeN_aGTfhnGvXcb_ZesLO8

## Funcionalidades principais adicionadas nesse projeto.
<ol>
   <li>Padrão MVC(model, view, controller)</li>
   <p>Introduzi pela segunda vez, porém de forma mais aprimorada o padrão MVC</p>
   <li>Algoritmo de Recuperação de Senha Avançado.</li>
   <p>São 3 estágios. 1. o usuário tenta recuperar a senha pelo email, ao digitar o sistema verifica se o e-mail está cadastrado e se pertence aquele usuário, depois ele envia um código de 6 dígitos para o email do usuário. 2. o usuário digita e enviar o código para API, ela verifica se está correto e caso sim, retorna um token especial para recuperar a senha. 3. com o token especial o usuário pode alterar a senha.</p>
   <li>Migrations e PostgreSQL</li>
   <p>Utilizei o banco postgresql junto do migration para para manipulá-lo. Fiz relacionamento avançado entre as tabelas do banco de dados.</p>
   <li>Autenticação JWT Token</li>
   <p>Utilizei a biblioteca do JWT do Golang, para criar autenticações seguras, assinadas com uma chave secreta. Implementei a autenticação de forma que a API gere dois tokens, o access token (atoken) e o refresh token (rtoken). O atoken é retornado, quando um usuário loga no sistema e tem duração de 4 horas. Já o refresh token é salvo no banco de dados e tem duração de 1 semana. quando o atoken expira e o usuário acessa um endpoint que exige autenticação, a API verificar se esse usuário tem um rtoken, caso sim ele gera um novo atoken e um novo rtoken. Dessa forma, desde que rtoken do usuário não expire, ele não precisará se logar novamente na aplicação, pois a API sempre fará essa verificação e caso não tenha problemas, deixará o usuário acessar os endpoints que exigem autenticação, assim melhorando a experiência do usuário.</p>
   <li>Paginação em rotas de listar</li>
   <p>Todos os endpoints de listar, tem o mecanismo de paginação, onde o frontend pode definir um começo (offset) um fim (limit) e os próximos valores (page), todos os valores sendo passados por querys request.</p>
   <li>Notificações via Gmail para usuários admins</li>
   <p>introduzir um algoritmo para enviar emails quando algum usuário interage com o blog, sendo, comentando ou respondendo comentários. o algoritmo envia um email comunicando ao(s) adm (s).</p>
   <li>Utilizei uma biblioteca criada por mim para validar as requisições.</li>  
   <p>a biblioteca validator, eu fiz com Golang com o intuito de validar os dados das requisições. eu hospedei a biblioteca de forma privada no meu github.</p>
   <li>DockerFile para produção.</li>
   <p>Criei um dockerfile, assim como em um projeto anterior, configurei de acordo com os requisitos do projeto. O dockerfile criar um binário da aplicação e o executa dentro de uma docker image menor chamada distroless, onde o tamanho do deploy diminuiu muito, chegando a ficar em torno de 26 mb.<p>
   <li>Documentação do backend</li>
   <p>Por fim, eu criei uma documentação para o projeto, mostrando como funciona a arquitetura, os endpoints, a modelagem do banco e os erros da API.</p>
</ol>
