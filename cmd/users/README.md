# Pacote `cmd/users`

Módulo para gerenciar usuários.

### Listar usuários

Para listar todos os usuários, utilize o seguinte comando:

```terminal
$ curl -X GET localhost:8080/users
```

### Criar um novo usuário

Para criar um novo usuário, utilize o seguinte comando:

```terminal
$ curl -X POST localhost:8080/users -d '{"id": "jpkroehling", "name":"Juraci Paixão Kröhling", "email":"juraci@example.com"}'
```

## Inicialização

Para iniciar a aplicação, execute o comando abaixo:

```terminal
$ go run ./cmd/users/
```
