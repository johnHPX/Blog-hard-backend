# estrutura das requisições e respostas

tanto as requisições quanto as respostas serão dadas por JSON.

#### - _BODY_

exemplos de JSON a serem enviados pelo corpo da requisição.

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

#### - _queries_

passando dados pela URL.

/example?attribute=value&attribute=value



# Rotas

## 1. [HOST:PORT]/categoria

criando um nova categoria

#### - _Request_

| request | type   | method |
| ------- | ------ | ------ |
| body    | object | POST   |

| attribute name | type     | size  | is it required? | type send      | description                                      |
| -------------- | -------- | ----- | --------------- | -------------- | ------------------------------------------------ |
| `kind`         | `string` | `255` | `true`          | body paraments | tipo da categoria                                |
| `mid`          | `string` | `-`   | `false`         | body paraments | mensagem da resposta caso o codigo http seja 200 |


#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |

| attribute name | type     | description                                      |
| -------------- | -------- | ------------------------------------------------ |
| `mid`          | `string` | mensagem da resposta caso o codigo http seja 200 |


## 2. [HOST:PORT]/categorias

listando categorias

#### - _Request_

| request | type | method |
| ------- | ---- | ------ |
| queries | -    | GET    |


| attribute name | type     | size | is it required? | description                                      |
| -------------- | -------- | ---- | --------------- | ------------------------------------------------ |
| `mid`          | `string` | `-`  | `false`         | mensagem da resposta caso o codigo http seja 200 |


#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |


| attribute name | type          | description                                      |
| -------------- | ------------- | ------------------------------------------------ |
| `count`        | `int`         | numero linhas trazidas do database               |
| `categorias`   | `[]Categoria` | array de categorias                              |
| `mid`          | `string`      | mensagem da resposta caso o codigo http seja 200 |

| Categoria | type     | description       |
| --------- | -------- | ----------------- |
| `id`      | `string` | id da categoria   |
| `kind`    | `string` | tipo da categoria |


## 3. [HOST:PORT]/categoria/{id}/find

buscando uma categoria por id de categoria

#### - _Request_

| request | type | method |
| ------- | ---- | ------ |
| queries | -    | GET    |


| attribute name | type     | size | is it required? | description                                      |
| -------------- | -------- | ---- | --------------- | ------------------------------------------------ |
| `mid`          | `string` | `-`  | `false`         | mensagem da resposta caso o codigo http seja 200 |

#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |

| attribute name | type     | description                                      |
| -------------- | -------- | ------------------------------------------------ |
| `id`           | `string` | id da categoria                                  |
| `kind`         | `string` | tipo da categoria                                |
| `mid`          | `string` | mensagem da resposta caso o codigo http seja 200 |

## 4. [HOST:PORT]/categoria/{id}/update

atualizando uma categoria por id de categoria

#### - _Request_

| request | type   | method |
| ------- | ------ | ------ |
| body    | object | PUT    |

| attribute name | type     | size  | is it required? | description                                      |
| -------------- | -------- | ----- | --------------- | ------------------------------------------------ |
| `kind`         | `string` | `255` | `true`          | tipo da categoria                                |
| `mid`          | `string` | `-`   | `false`         | mensagem da resposta caso o codigo http seja 200 |


#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |

| attribute name | type     | description                                      |
| -------------- | -------- | ------------------------------------------------ |
| `mid`          | `string` | mensagem da resposta caso o codigo http seja 200 |

## 5. [HOST:PORT]/categoria/{id}/remove

deletando uma categoria por id de categoria

#### - _Request_

| request | type | method |
| ------- | ---- | ------ |
| queries | -    | DELETE |

| attribute name | type     | size | is it required? | description                                      |
| -------------- | -------- | ---- | --------------- | ------------------------------------------------ |
| `mid`          | `string` | `-`  | `false`         | mensagem da resposta caso o codigo http seja 200 |

#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |

| attribute name | type     | description                                      |
| -------------- | -------- | ------------------------------------------------ |
| `mid`          | `string` | mensagem da resposta caso o codigo http seja 200 |

## 6. [HOST:PORT]/noticia

criando um nova notícia

#### - _Request_

| request | type   | method |
| ------- | ------ | ------ |
| body    | object | POST   |

| attribute name | type         | size  | is it required? | description                                      |
| -------------- | ------------ | ----- | --------------- | ------------------------------------------------ |
| `titulo`       | `string`     | `255` | `true`          | titulo da notícia                                |
| `conteudos`    | `[]conteudo` | `-`   | `true`          | topicos da notícia                               |
| `categoria`    | `string`     | `100` | `true`          | categoria da notícia                             |
| `mid`          | `string`     | `-`   | `false`         | mensagem da resposta caso o codigo http seja 200 |

| conteudo    | type     | size   | is it required? | description                            |
| ----------- | -------- | ------ | --------------- | -------------------------------------- |
| `subTitulo` | `string` | `255`  | `true`          | subtitulo da notícia(titulo do topico) |
| `texto`     | `string` | `5000` | `true`          | texto do topico                        |


#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |


| attribute name | type     | description                                      |
| -------------- | -------- | ------------------------------------------------ |
| `mid`          | `string` | mensagem da resposta caso o codigo http seja 200 |


## 7. [HOST:PORT]/noticias

listando todas as noticias

#### - _Request_

| request | type | method |
| ------- | ---- | ------ |
| queries | -    | GET    |


| attribute name | type     | size | is it required? | description                                      |
| -------------- | -------- | ---- | --------------- | ------------------------------------------------ |
| `mid`          | `string` | `-`  | `false`         | mensagem da resposta caso o codigo http seja 200 |


#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |


| attribute name | type         | description                                      |
| -------------- | ------------ | ------------------------------------------------ |
| `count`        | `int`        | numero linhas trazidas do database               |
| `noticias`     | `[]Noticias` | array de noticias                                |
| `mid`          | `string`     | mensagem da resposta caso o codigo http seja 200 |

| Noticias    | type         | description          |
| ----------- | ------------ | -------------------- |
| `id`        | `string`     | id da noticia        |
| `titulo`    | `string`     | titulo da notícia    |
| `conteudos` | `[]conteudo` | topicos da notícia   |
| `categoria` | `string`     | categoria da notícia |

| conteudo    | type     | description                            |
| ----------- | -------- | -------------------------------------- |
| `subTitulo` | `string` | subtitulo da notícia(titulo do topico) |
| `texto`     | `string` | texto do topico                        |


## 8. [HOST:PORT]/noticia/{titCat}/list

listando noticias por titulo ou categoria

#### - _Request_

| request | type | method |
| ------- | ---- | ------ |
| queries | -    | GET    |


| attribute name | type     | size | is it required? | description                                      |
| -------------- | -------- | ---- | --------------- | ------------------------------------------------ |
| `mid`          | `string` | `-`  | `false`         | mensagem da resposta caso o codigo http seja 200 |


#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |


| attribute name | type         | description                                      |
| -------------- | ------------ | ------------------------------------------------ |
| `count`        | `int`        | numero linhas trazidas do database               |
| `noticias`     | `[]Noticias` | array de noticias                                |
| `mid`          | `string`     | mensagem da resposta caso o codigo http seja 200 |

| Noticias    | type         | description          |
| ----------- | ------------ | -------------------- |
| `id`        | `string`     | id da noticia        |
| `titulo`    | `string`     | titulo da notícia    |
| `conteudos` | `[]conteudo` | topicos da notícia   |
| `categoria` | `string`     | categoria da notícia |

| conteudo    | type     | description                            |
| ----------- | -------- | -------------------------------------- |
| `subTitulo` | `string` | subtitulo da notícia(titulo do topico) |
| `texto`     | `string` | texto do topico                        |


## 9. [HOST:PORT]/noticia/{nid}/update

atualizando uma noticia por id de noticia

#### - _Request_

| request | type   | method |
| ------- | ------ | ------ |
| body    | object | PUT    |


| attribute name | type         | size  | is it required? | description                                      |
| -------------- | ------------ | ----- | --------------- | ------------------------------------------------ |
| `titulo`       | `string`     | `255` | `true`          | titulo da notícia                                |
| `conteudos`    | `[]conteudo` | `-`   | `true`          | topicos da notícia                               |
| `categoria`    | `string`     | `100` | `true`          | categoria da notícia                             |
| `mid`          | `string`     | `-`   | `false`         | mensagem da resposta caso o codigo http seja 200 |

| conteudo    | type     | size   | is it required? | description                            |
| ----------- | -------- | ------ | --------------- | -------------------------------------- |
| `cid`       | `string` | `36`   | `true`          | id do conteudo                         |
| `subTitulo` | `string` | `255`  | `true`          | subtitulo da notícia(titulo do topico) |
| `texto`     | `string` | `5000` | `true`          | texto do topico                        |


#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |


| attribute name | type     | description                                      |
| -------------- | -------- | ------------------------------------------------ |
| `mid`          | `string` | mensagem da resposta caso o codigo http seja 200 |


## 10. [HOST:PORT]/noticia/{nid}/remove

deletando uma noticia por id de noticia

#### - _Request_

| request | type | method |
| ------- | ---- | ------ |
| queries | -    | DELETE |


| attribute name | type     | size | is it required? | description                                      |
| -------------- | -------- | ---- | --------------- | ------------------------------------------------ |
| `mid`          | `string` | `-`  | `false`         | mensagem da resposta caso o codigo http seja 200 |


#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |


| attribute name | type     | description                                      |
| -------------- | -------- | ------------------------------------------------ |
| `mid`          | `string` | mensagem da resposta caso o codigo http seja 200 |

## 11. [HOST:PORT]/noticia/{id}/find

buscando uma noticia pelo seu id

#### - _Request_

| request | type | method |
| ------- | ---- | ------ |
| queries | -    | GET    |


| attribute name | type     | size | is it required? | description                                      |
| -------------- | -------- | ---- | --------------- | ------------------------------------------------ |
| `mid`          | `string` | `-`  | `false`         | mensagem da resposta caso o codigo http seja 200 |


#### - _Response_

| request | type   | status |
| ------- | ------ | ------ |
| body    | object | 200    |


| attribute name | type         | description                                      |
| -------------- | ------------ | ------------------------------------------------ |
| `id`           | `string`     | id da noticia                                    |
| `titulo`       | `string`     | titulo da notícia                                |
| `conteudos`    | `[]conteudo` | topicos da notícia                               |
| `categoria`    | `string`     | categoria da notícia                             |
| `mid`          | `string`     | mensagem da resposta caso o codigo http seja 200 |

| conteudo    | type     | description                            |
| ----------- | -------- | -------------------------------------- |
| `subTitulo` | `string` | subtitulo da notícia(titulo do topico) |
| `texto`     | `string` | texto do topico                        |



the end!
made by Jonatas.