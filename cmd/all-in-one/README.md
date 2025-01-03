# Pacote `cmd/all-in-one`

A ser documentado.

Mas por enquanto, aqui vão alguns exemplos: 

```terminal
$ nats-server -D -js
$ nats -s localhost:4222 stream create payments --subjects "payment.process" --storage memory --replicas 1 --retention=limits --discard=old --max-msgs 1_000_000 --max-msgs-per-subject 100_000 --max-bytes 4GiB --max-age 1d --max-msg-size 10MiB --dupe-window 2m --allow-rollup --no-deny-delete --no-deny-purge
$ go run ./cmd/all-in-one/
$ curl localhost:8080/payments
$ curl -X POST localhost:8080/users -d '{"id": "jpkroehling"}'
$ curl -X POST localhost:8080/subscriptions -d '{"id": "jpkroehling", "user_id":"jpkroehling", "plan_id":"silver"}'
$ curl -X POST localhost:8080/payments -d '{"id": "some-uuid", "subscription_id":"jpkroehling", "amount":99, "status":"FAILED"}'
$ nats -s localhost:4222 stream view payments
```
