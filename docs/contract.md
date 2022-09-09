# estrutura das requisições e respostas

request e response serão enviadas por JSON.

#### - _BODY_

exemplos de JSON a serem enviadoas pelo body's Request

```
## object

{
  "attribute: "value",
  "attribute": [{
    "attribute: "value",
  }]
}

## array

[
  {
    "attribute: "value",
  },
  {
    "attribute: "value",
  }
]

```

#### - _QUERIES_

passando dados pela URL.

/example?attribute=value&attribute=value

# CONCEITOS

  | Nome             | Descrição                                                                                        |
  | ---------------- | ------------------------------------------------------------------------------------------------ |
  | request          | diz se os dados em sua maioria serão enviados pelos 3 tipos de paramentros(body, queries ou url) |
  | type             | diz o tipo de estrutura que será enviada. exemplo: objeto, array...                              |
  | method           | diz o verbo http que a rota utiliza. exemplo (post, get, put e delete)                           |
  | toke is requered | diz se é necessario um token de acesso para utilização desse endpoint                            |
  | attribute name   | diz o nome do atributo                                                                           |
  | type value       | diz a tipagem do atributo                                                                        |
  | size             | diz o valoz maximo de caracteres que o atributo pode ter                                         |
  | is it requerid?  | diz se o atributo é obrigatorio                                                                  |
  | type send        | diz a maneira que o atributo será enviado. exemplo: body, queries ou url paraments               |
  | description      | a descrição do atributo                                                                          |
  | status           | o status do response. exemplo: 200, 201...                                                       |

<hr>
<h1> USER Routes </h1>


## 1. [HOST:PORT]/user/store

criando uma novo usuario

#### - _Request_

| request | type   | method | token is required |
| ------- | ------ | ------ | ----------------- |
| body    | object | POST   | not               |

| attribute name | type value | size  | is it required? | type send      | description                                      |
| -------------- | ---------- | ----- | --------------- | -------------- | ------------------------------------------------ |
| `name`         | `string`   | `255` | `true`          | body paraments | nome do usuario                                  |
| `telephone`    | `string`   | `13`  | `true`          | body paraments | telefone do usuario                              |
| `nick`         | `string`   | `255` | `true`          | body paraments | nick do usuario                                  |
| `email`        | `string`   | `255` | `true`          | body paraments | email do usuario                                 |
| `secret`       | `string`   | `255` | `true`          | body paraments | senha do usuario                                 |
| `mid`          | `string`   | `-`   | `false`         | body paraments | mensagem da resposta caso o codigo http seja 200 |

#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |

| attribute name | type value | description                                      |
| -------------- | ---------- | ------------------------------------------------ |
| `mid`          | `string`   | mensagem da resposta caso o codigo http seja 200 |

## 2. [HOST:PORT]/user/adm/store

criando usuário admin.
somente usuario admin pode utilizar essa rota.

#### - _Request_

| request | type   | method | token is required |
| ------- | ------ | ------ | ----------------- |
| body    | object | POST   | yes               |

| attribute name | type value | size  | is it required? | type send      | description                                      |
| -------------- | ---------- | ----- | --------------- | -------------- | ------------------------------------------------ |
| `name`         | `string`   | `255` | `true`          | body paraments | nome do usuario                                  |
| `telephone`    | `string`   | `13`  | `true`          | body paraments | telefone do usuario                              |
| `nick`         | `string`   | `255` | `true`          | body paraments | nick do usuario                                  |
| `email`        | `string`   | `255` | `true`          | body paraments | email do usuario                                 |
| `secret`       | `string`   | `255` | `true`          | body paraments | senha do usuario                                 |
| `mid`          | `string`   | `-`   | `false`         | body paraments | mensagem da resposta caso o codigo http seja 200 |

#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |

| attribute name | type value | description                                      |
| -------------- | ---------- | ------------------------------------------------ |
| `mid`          | `string`   | mensagem da resposta caso o codigo http seja 200 |

## 3. [HOST:PORT]/user/list

listando todos os usuarios
somente usuario admin pode utilizar essa rota.

#### - _Request_

| request | type | method | token is required |
| ------- | ---- | ------ | ----------------- |
| queries | -    | GET    | yes               |

| attribute name | type value | size | is it required? | type send         | description                                         |
| -------------- | ---------- | ---- | --------------- | ----------------- | --------------------------------------------------- |
| `offset`       | `int`      | `-`  | `false`         | queries paraments | deslocamento inicial dos dados trazidos             |
| `limit`        | `int`      | `-`  | `false`         | queries paraments | limite padrão de quantos dados serão trazidos       |
| `page`         | `int`      | `-`  | `false`         | queries paraments | o numero da pagina na qual os dados estão agrupados |
| `mid`          | `string`   | `-`  | `false`         | queries paraments | mensagem da resposta caso o codigo http seja 200    |

#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |

| attribute name | type value | description                                      |
| -------------- | ---------- | ------------------------------------------------ |
| `count`        | `int`      | numero total de dados do banco dessa query       |
| `users`        | `[]User`   | array de users                                   |
| `mid`          | `string`   | mensagem da resposta caso o codigo http seja 200 |

| User        | type     | description        |
| ----------- | -------- | ------------------ |
| `personID`  | `string` | id da pessoa       |
| `userID`    | `string` | id do usuario      |
| `name`      | `string` | nome da pessoa     |
| `telephone` | `string` | telefone da pessoa |
| `nick`      | `string` | nick do usuario    |
| `email`     | `string` | email do usuario   |
| `kind`      | `string` | tipo do usuario    |

## 4. [HOST:PORT]/user/list/name/{name}

listando usuarios pelo nome
somente usuario adimin tem acesso a essa rota

#### - _Request_

| request | type | method | token is required |
| ------- | ---- | ------ | ----------------- |
| queries | -    | GET    | yes               |

| attribute name | type value | size | is it required? | type send         | description                                         |
| -------------- | ---------- | ---- | --------------- | ----------------- | --------------------------------------------------- |
| `offset`       | `int`      | `-`  | `false`         | queries paraments | deslocamento inicial dos dados trazidos             |
| `limit`        | `int`      | `-`  | `false`         | queries paraments | limite padrão de quantos dados são trazidos         |
| `page`         | `int`      | `-`  | `false`         | queries paraments | o numero da pagina na qual os dados estão agrupados |
| `mid`          | `string`   | `-`  | `false`         | queries paraments | mensagem da resposta caso o codigo http seja 200    |

#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |

| attribute name | type value | description                                      |
| -------------- | ---------- | ------------------------------------------------ |
| `count`        | `int`      | numero total de dados do banco dessa query       |
| `users`        | `[]User`   | array de users                                   |
| `mid`          | `string`   | mensagem da resposta caso o codigo http seja 200 |

| User        | type     | description        |
| ----------- | -------- | ------------------ |
| `personID`  | `string` | id da pessoa       |
| `userID`    | `string` | id do usuario      |
| `name`      | `string` | nome da pessoa     |
| `telephone` | `string` | telefone da pessoa |
| `nick`      | `string` | nick do usuario    |
| `email`     | `string` | email do usuario   |
| `kind`      | `string` | tipo do usuario    |

## 5. [HOST:PORT]/user/find/id/{id}

buscar um usuario pelo id.
somente admins podem usar essa rota.

#### - _Request_

| request | type | method | token is required |
| ------- | ---- | ------ | ----------------- |
| queries | -    | GET    | yes               |

| attribute name | type value | size | is it required? | type send         | description                                      |
| -------------- | ---------- | ---- | --------------- | ----------------- | ------------------------------------------------ |
| `id`           | `string`   | `-`  | `true`          | url paraments     | id do usuario                                    |
| `mid`          | `string`   | `-`  | `false`         | queries paraments | mensagem da resposta caso o codigo http seja 200 |

#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |

| attribute name | type value | description                                      |
| -------------- | ---------- | ------------------------------------------------ |
| `personID`     | `string`   | id da pessoa                                     |
| `userID`       | `string`   | id do usuario                                    |
| `name`         | `string`   | nome da pessoa                                   |
| `telephone`    | `string`   | telefone da pessoa                               |
| `nick`         | `string`   | nick do usuario                                  |
| `email`        | `string`   | email do usuario                                 |
| `kind`         | `string`   | tipo do usuario                                  |
| `mid`          | `string`   | mensagem da resposta caso o codigo http seja 200 |

## 6. [HOST:PORT]/user/update/id/{id}

atualizar um usuario por id.

#### - _Request_

| request | type   | method | token is required |
| ------- | ------ | ------ | ----------------- |
| body    | object | PUT    | yes               |

| attribute name | type value | size  | is it required? | type send      | description                                      |
| -------------- | ---------- | ----- | --------------- | -------------- | ------------------------------------------------ |
| `id`           | `string`   | `36`  | `true`          | url paraments  | id do usuario                                    |
| `name`         | `string`   | `255` | `true`          | body paraments | nome do usuario                                  |
| `telephone`    | `string`   | `13`  | `true`          | body paraments | telefone do usuario                              |
| `nick`         | `string`   | `255` | `true`          | body paraments | nick do usuario                                  |
| `email`        | `string`   | `255` | `true`          | body paraments | email do usuario                                 |
| `kind`         | `string`   | `10`  | `true`          | body paraments | kind do usuario                                  |
| `mid`          | `string`   | `-`   | `false`         | body paraments | mensagem da resposta caso o codigo http seja 200 |

#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |

| attribute name | type value | description                                      |
| -------------- | ---------- | ------------------------------------------------ |
| `mid`          | `string`   | mensagem da resposta caso o codigo http seja 200 |

## 7. [HOST:PORT]/user/remove/id/{id}

remover um usuario por id.

#### - _Request_

| request | type | method | token is required |
| ------- | ---- | ------ | ----------------- |
| queries | -    | DELETE | yes               |

| attribute name | type value | size | is it required? | type send         | description                                      |
| -------------- | ---------- | ---- | --------------- | ----------------- | ------------------------------------------------ |
| `id`           | `string`   | `36` | `true`          | url paraments     | id do usuario                                    |
| `mid`          | `string`   | `-`  | `false`         | queries paraments | mensagem da resposta caso o codigo http seja 200 |

#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |

| attribute name | type value | description                                      |
| -------------- | ---------- | ------------------------------------------------ |
| `mid`          | `string`   | mensagem da resposta caso o codigo http seja 200 |

# 8. [HOST:PORT]/user/login

fazer login no sistema

#### - _Request_

| request | type   | method | token is required |
| ------- | ------ | ------ | ----------------- |
| body    | object | POST   | not               |

| attribute name | type value | size  | is it required? | type send      | description                                      |
| -------------- | ---------- | ----- | --------------- | -------------- | ------------------------------------------------ |
| `nick`         | `string`   | `255` | `true`          | body paraments | email ou nick do usuario                         |
| `password`     | `string`   | `255` | `true`          | body paraments | senha do usuario                                 |
| `mid`          | `string`   | `-`   | `false`         | body paraments | mensagem da resposta caso o codigo http seja 200 |

#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |

| attribute name | type value | description                                      |
| -------------- | ---------- | ------------------------------------------------ |
| `token`        | `string`   | token de acesso a aplicação                      |
| `mid`          | `string`   | mensagem da resposta caso o codigo http seja 200 |

# 9. [HOST:PORT]/user/recor/email

primeiro estagio de recuperação de senha.
enviar um codigo valido ao email do usuario.

#### - _Request_

| request | type   | method | token is required |
| ------- | ------ | ------ | ----------------- |
| body    | object | POST   | not               |

| attribute name | type value | size  | is it required? | type send      | description                                      |
| -------------- | ---------- | ----- | --------------- | -------------- | ------------------------------------------------ |
| `email`        | `string`   | `255` | `true`          | body paraments | email do usuario                                 |
| `mid`          | `string`   | `-`   | `false`         | body paraments | mensagem da resposta caso o codigo http seja 200 |

#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |

| attribute name | type value | description                                      |
| -------------- | ---------- | ------------------------------------------------ |
| `mid`          | `string`   | mensagem da resposta caso o codigo http seja 200 |

# 10. [HOST:PORT]/user/verific/code

segundo estagio de recuperação de senha.
conferi o codigo enviado e devover um token especial.

#### - _Request_

| request | type   | method | token is required |
| ------- | ------ | ------ | ----------------- |
| body    | object | POST   | not               |

| attribute name | type value | size | is it required? | type send      | description                                      |
| -------------- | ---------- | ---- | --------------- | -------------- | ------------------------------------------------ |
| `code`         | `string`   | `6`  | `true`          | body paraments | codigo                                           |
| `mid`          | `string`   | `-`  | `false`         | body paraments | mensagem da resposta caso o codigo http seja 200 |

#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |

| attribute name | type value | description                                      |
| -------------- | ---------- | ------------------------------------------------ |
| `token`        | `string`   | token especial para recuperção de senha          |
| `mid`          | `string`   | mensagem da resposta caso o codigo http seja 200 |

# 11. [HOST:PORT]/user/password/recovery

terceiro estagio de recuperação de senha.
atualizar senha do usuario.

#### - _Request_

| request | type   | method | token is required |
| ------- | ------ | ------ | ----------------- |
| body    | object | POST   | yes               |

| attribute name | type value | size  | is it required? | type send      | description                                      |
| -------------- | ---------- | ----- | --------------- | -------------- | ------------------------------------------------ |
| `newPassword`  | `string`   | `255` | `true`          | body paraments | nova senha do usuario                            |
| `mid`          | `string`   | `-`   | `false`         | body paraments | mensagem da resposta caso o codigo http seja 200 |

#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |

| attribute name | type value | description                                      |
| -------------- | ---------- | ------------------------------------------------ |
| `mid`          | `string`   | mensagem da resposta caso o codigo http seja 200 |

# 12. [HOST:PORT]/user/password/update

atualizar senha do usuario.

#### - _Request_

| request | type   | method | token is required |
| ------- | ------ | ------ | ----------------- |
| body    | object | POST   | yes               |

| attribute name | type value | size  | is it required? | type send      | description                                      |
| -------------- | ---------- | ----- | --------------- | -------------- | ------------------------------------------------ |
| `oldPassword`  | `string`   | `255` | `true`          | body paraments | antiga senha do usuario                          |
| `newPassword`  | `string`   | `255` | `true`          | body paraments | nova senha do usuario                            |
| `mid`          | `string`   | `-`   | `false`         | body paraments | mensagem da resposta caso o codigo http seja 200 |

#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |

| attribute name | type value | description                                      |
| -------------- | ---------- | ------------------------------------------------ |
| `mid`          | `string`   | mensagem da resposta caso o codigo http seja 200 |


<hr>
<h1> POST Routes </h1>



the end!
made by Jonatas.
