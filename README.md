# nats-groups

run NATS server: docker run -d -p4222:4222 -p6222:6222 -p8222:8222 -p6543:6543 --name=nats-groups nats:2.4.0-alpine3.14 nats-server -js --port=4222

# Разные durable

Создаем 2 подписки с разными durable именами, но одинаковыми queue и subject. По идее они должны входить в одну группу получателей.

go run .\subscribe\subscribe.go -diff

Создается 2 подписки, прикрепленные к разным консюмерамю

Публикуем сообщение в subject:

go run .\publish\publish.go

В output видим, что каждая подписка получила сообщение, не смотря на то, что они подписаны на одну queue.

```shell
{"L":"INFO","T":"2021-10-05T15:09:17.513+0300","M":"subject alerts.subject, group queue: msg received","pid":13520,"message":{"Subject":"alerts.subject"
,"Reply":"$JS.ACK.alerts.subscription_1.1.1.1.1633435756647117100.0","Header":null,"Data":"dGVzdCBtZXNzYWdl","Sub":{"Subject":"_INBOX.wCy32uQLSh0wBF8Iwt
OHm6","Queue":"queue"}}}
{"L":"INFO","T":"2021-10-05T15:09:17.513+0300","M":"subject alerts.subject, group queue: msg received","pid":13520,"message":{"Subject":"alerts.subject"
,"Reply":"$JS.ACK.alerts.subscription_2.1.1.1.1633435756647117100.0","Header":null,"Data":"dGVzdCBtZXNzYWdl","Sub":{"Subject":"_INBOX.wCy32uQLSh0wBF8Iwt
OHn3","Queue":"queue"}}}
```

Подписки с разным именами durable не объединяются в группы, даже если подписаны на одну queue?

Для того, чтобы работали queue subscriptions нужн одинаковые durable имена (чтобы создавался один консюмер)?

# Одинаковые durable

Создаем 2 подписки с одинаковыми durable именами, queue и subject.

go run .\subscribe\subscribe.go

Создается 2 подписки, прикрепленные к одному и тому же консьюмеру.

Публикуем сообщение в subject:

go run .\publish\publish.go

В output видим, что только одна подписка получила сообщение.

Всё сработало ожидаемо.