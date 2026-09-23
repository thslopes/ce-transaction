# Contratos de API de Acesso e Administracao

## Objetivo

Este documento define os contratos minimos das rotas de cadastro, autenticacao e administracao de usuarios para o MVP. O foco e preparar a implementacao incremental e os primeiros testes de contrato.

## Convencoes Gerais

- Formato: JSON sobre HTTPS.
- Autenticacao: sessao ou token de acesso para rotas protegidas.
- Autorizacao: validacao de perfil e status do usuario em toda rota protegida.
- Status validos de usuario: `inactive`, `active`.
- Perfis validos no MVP: `admin`, `requester`, `approver`, `executor`.

## Transporte de Tokens

- Tokens de autenticacao de usuario devem trafegar no header `Authorization: Bearer <token>`.
- A resposta de autenticacao tambem retorna o token no header `Authorization`.
- O token do link de cadastro pode aparecer na URL apenas na entrada do fluxo web, por exemplo `GET /signup?token=...`.
- A chamada de API que conclui o cadastro deve enviar o token no corpo da requisicao, nao depender da query string.

## Regras Transversais

- Concluir cadastro nao concede acesso automaticamente.
- Usuario autenticado com `status = inactive` nao acessa rotas protegidas.
- Usuario sem perfis nao acessa rotas protegidas.
- Toda resposta de erro deve ser previsivel e estruturada.

## Estrutura Padrao de Erro

```json
{
  "error": {
    "code": "forbidden",
    "message": "User does not have permission to perform this action",
    "details": {}
  }
}
```

Codigos esperados no MVP:

- `invalid_request`
- `invalid_token`
- `expired_token`
- `invalid_credentials`
- `forbidden`
- `not_found`
- `conflict`
- `validation_error`

## 1. Gerar Link de Cadastro

### `POST /users/registration-link`

### Acesso

- Requer usuario `active`
- Requer perfil `admin`

### Request

```json
{}
```

Corpo vazio e suficiente no MVP.

### Response `201 Created`

```json
{
  "registration_url": "https://app.example.com/signup?token=eyJ...",
  "expires_at": "2026-09-23T12:10:00Z"
}
```

### Erros esperados

- `401 Unauthorized`
- `403 Forbidden`

## 2. Realizar Cadastro

### `POST /signup`

### Acesso

- Publico
- Exige token valido no corpo da requisicao

### Entrada no fluxo web

O frontend pode abrir a pagina de cadastro por meio de uma URL como:

`GET /signup?token=eyJ...`

Depois de carregar a pagina, a aplicacao cliente deve usar o token recebido para chamar `POST /signup`.

### Request

```json
{
  "token": "eyJ...",
  "name": "Thiago Lopes",
  "email": "thiago@example.com",
  "phone": "+5511999999999",
  "password": "strong-password",
  "password_confirmation": "strong-password"
}
```

### Response `201 Created`

```json
{
  "id": "usr_123",
  "name": "Thiago Lopes",
  "email": "thiago@example.com",
  "phone": "+5511999999999",
  "status": "inactive",
  "profiles": [],
  "created_at": "2026-09-23T12:05:00Z"
}
```

### Regras adicionais

- `email` deve ser unico.
- `phone` e obrigatorio.
- `password` e `password_confirmation` devem coincidir.

### Erros esperados

- `400 Bad Request` com `validation_error`
- `401 Unauthorized` com `invalid_token`
- `401 Unauthorized` com `expired_token`
- `409 Conflict` se o usuario ja existir

## 3. Autenticar Usuario

### `POST /login`

### Acesso

- Publico

### Request

```json
{
  "email": "thiago@example.com",
  "password": "strong-password"
}
```

### Response `200 OK`

Header:

```http
Authorization: Bearer <access_token>
```

```json
{
  "user": {
    "id": "usr_123",
    "name": "Thiago Lopes",
    "email": "thiago@example.com",
    "phone": "+5511999999999",
    "status": "active",
    "profiles": ["requester"]
  }
}
```

### Regras adicionais

- Usuario `inactive` nao autentica com sucesso.
- Usuario sem perfis nao autentica com sucesso.

### Erros esperados

- `401 Unauthorized` com `invalid_credentials`
- `403 Forbidden` com `forbidden` quando o usuario estiver `inactive`

### Header esperado nas rotas autenticadas

```http
Authorization: Bearer <access_token>
```

## 4. Consultar Proprio Perfil

### `GET /me`

### Acesso

- Requer usuario `active`

### Response `200 OK`

```json
{
  "id": "usr_123",
  "name": "Thiago Lopes",
  "email": "thiago@example.com",
  "phone": "+5511999999999",
  "status": "active",
  "profiles": ["requester"],
  "validated_at": "2026-09-23T12:15:00Z",
  "validated_by": "usr_admin_1"
}
```

## 5. Listar Usuarios

### `GET /users`

### Acesso

- Requer usuario `active`
- Requer perfil `admin`

### Query params opcionais

- `status`
- `profile`
- `search`
- `page`
- `page_size`

### Response `200 OK`

```json
{
  "items": [
    {
      "id": "usr_123",
      "name": "Thiago Lopes",
      "email": "thiago@example.com",
      "phone": "+5511999999999",
      "status": "inactive",
      "profiles": [],
      "created_at": "2026-09-23T12:05:00Z",
      "validated_at": null,
      "validated_by": null
    }
  ],
  "page": 1,
  "page_size": 20,
  "total": 1
}
```

## 6. Validar Usuario

### `POST /users/{id}/validate`

### Acesso

- Requer usuario `active`
- Requer perfil `admin`

### Request

```json
{
  "profiles": ["requester"]
}
```

### Response `200 OK`

```json
{
  "id": "usr_123",
  "status": "active",
  "profiles": ["requester"],
  "validated_at": "2026-09-23T12:15:00Z",
  "validated_by": "usr_admin_1"
}
```

### Regras adicionais

- A lista `profiles` nao pode ser vazia.
- Perfis invalidos devem ser rejeitados.
- No MVP, a validacao ativa o usuario na mesma operacao.

### Erros esperados

- `400 Bad Request` com `validation_error`
- `404 Not Found`
- `409 Conflict` se a regra escolhida proibir validar usuario ja ativo

## 7. Desativar Usuario

### `POST /users/{id}/deactivate`

### Acesso

- Requer usuario `active`
- Requer perfil `admin`

### Request

```json
{}
```

### Response `200 OK`

```json
{
  "id": "usr_123",
  "status": "inactive"
}
```

### Regras adicionais

- O usuario nao e removido fisicamente.
- O historico e preservado.

## 8. Reativar Usuario

### `POST /users/{id}/reactivate`

### Acesso

- Requer usuario `active`
- Requer perfil `admin`

### Request

```json
{}
```

### Response `200 OK`

```json
{
  "id": "usr_123",
  "status": "active",
  "profiles": ["requester"]
}
```

### Regras adicionais

- Reativacao exige que o usuario ja possua ao menos um perfil associado.

## Ordem Recomendada de Testes de Contrato

1. `POST /users/registration-link`
2. `POST /signup`
3. `GET /users`
4. `POST /users/{id}/validate`
5. `POST /login`
6. `GET /me`
7. `POST /users/{id}/deactivate`
8. `POST /users/{id}/reactivate`