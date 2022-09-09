# Estrutura de erros

Quando algum endpoint estiver com erros, a API devolvera essa estrutura

STATUS = 401, 500 => {
  "status: interger,
  "code": integer,
  "message": string,
  "mid": string,
}

| STATUS | Desciptions            |
| ------ | ---------------------- |
| 401    | usuario não autorizado |
| 500    | error interno na API   |

| Atribute Name | Type Value | Description               |
| ------------- | ---------- | ------------------------- |
| `status`      | `int`      | `codigo do status`        |
| `code`        | `int`      | `codigo da API`           |
| `message`     | `string`   | `messagem de erro`        |
| `mid`         | `string`   | `messagem de verificação` |