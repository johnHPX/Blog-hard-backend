# BlogHard
Um sistema de gerenciamento de conteúdo(CMS) feito com GO, postgress e docker.<br/>Futuramente irei desenvolver o front-end para está API.

Toda documentação necessaria está na pasta "docs" e nesse link: https://drive.google.com/drive/u/0/folders/1z8u4znIMH_JeN_aGTfhnGvXcb_ZesLO8

## Principais funcionalidades implementadas neste projeto.
<ol>
   <li>Arquitetura de software baseada em camadas.</li>
   <p>Arquitetura limpa e de facil manutenção.</p>
   <li>Algoritmo de recuperação de senha.</li>
   <p>São 3 estágios:<br/>(1) O usuário tenta recuperar a senha pelo email, ao digitar o sistema verifica se o e-mail está cadastrado e se pertence aquele usuário, depois é enviado um código de 6 dígitos para o email do usuário.<br/>(2) O usuário digita e enviar o código para API, ela verifica se está correto e caso sim, retorna um token especial para recuperar a senha.<br/>(3) com o token especial o usuário pode alterar a senha.</p>
   <li>Migrations e PostgreSQL.</li>
   <p>Utilizando o PostgreSqL, um dos bancos de dados mais robustos da atualidade junto com um script capaz de realizar migrações para o mesmo. Assim, quando houver alterações nas tabelas, a restruturação do banco de dados será bem mais simples. Pois com apenas dois ou mais comandos, é possível recria-lo.</p>
   <li>Autenticação JWT Token.</li>
   <p>Utilizei a biblioteca do JWT do Golang, para criar autenticações seguras, assinadas com uma chave secreta. Implementei a autenticação de forma que a API gere dois tokens, o access token (atoken) e o refresh token (rtoken). O atoken é retornado, quando um usuário loga no sistema e tem duração de 4 horas. Já o refresh token é salvo no banco de dados e tem duração de 1 semana. quando o atoken expira e o usuário acessa um endpoint que exige autenticação, a API verificar se esse usuário tem um rtoken, caso sim ele gera um novo atoken e um novo rtoken. Dessa forma, desde que rtoken do usuário não expire, ele não precisará se logar novamente na aplicação, pois a API sempre fará essa verificação e caso não tenha problemas, deixará o usuário acessar os endpoints que exigem autenticação, assim melhorando a experiência do usuário.</p>
   <li>Paginação.</li>
   <p>Todos os endpoints de listar enitdades, tem o mecanismo de paginação, onde o frontend pode definir um começo (offset) um fim (limit) e os próximos valores (page), todos os valores sendo passados por querys request.</p>
   <li>Notificações pelo Gmail para usuários administradores.</li>
   <p>Introduzir um algoritmo que enviar emails quando algum usuário interage com o blog. A interação poder ocorrer de duas formas: O usuário pode comentar em alguma postagem ou pode responder um comentário de um outro usuário. E quando ocorrer essa interação, a API informa a todos o(s) usuário(s) administradore(s).
</p>
   <li>Validator lib.</li>  
   <p>Eu fiz essa lib com Golang para validar os dados das requisições. Aqui está o repósitorio: https://github.com/johnHPX/validator-hard</p>
   <li>DockerFile.</li>
   <p>Criei um DockerFile, assim como em um projeto anterior, configurei de acordo com os requisitos do projeto. O dockerfile criar um binário da aplicação e o executa dentro de uma DockerImage bem menor chamada distroless, na qual o seu tamanho é muito pequeno(26 mb), o que ajudou muito o deploy para a produção.<p>
   <li>Documentação.</li>
   <p>Por fim, eu criei uma documentação para o projeto, mostrando como funciona a arquitetura, os endpoints e a modelagem do banco de dados.<p>
</ol>
