# Golang com WebSockets - Exemplo 3

- Rode o código

```bash
go run main.go
```

## Inscrevendo-se em um Tópico:

```js
ws.send('{"action": "subscribe", "topic": "Teste"}');
```

## Publicando uma mensagem: 

```js
ws.send('{"action":"publish", "topic":"Teste", "message":"Hello world!"}');
```

## Desinscrevendo-se de um Tópico:
```js
ws.send('{"action": "unsubscribe", "topic": "Teste"}');
```
