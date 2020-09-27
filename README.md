# Avito
[![Actions Status](https://github.com/eprokofyev/testTask/workflows/avitobuild/badge.svg?branch=master)](https://github.com/eprokofyev/avito/actions)


## Запуск

docker-compose up

## Что сделано

* [x] Реализованы все необходимые методы
* [x] Код написан на Go
* [x] Использована база данных PostgreSQL
* [x] Развертывание с помощью Docker
* [x] API возвращает коды ошибок с описанием
* [x] Выполнены дополнительные задания


## HTTP API
#### Post /api/transfer - метод, реализующий переводы между пользователями, списание и получение средств.

recipient_id - id пользователя, получающего деньги

sender_id - id пользователя, у которого списываются деньги

Примеры:

* Запрос: POST http://127.0.0.1:8080/api/transfer
{"recipient_id":10, "amount":1000, "message":"hello"}

    Ответ: {"status":201,"body":{"message":"OK"}}

* Запрос: POST http://127.0.0.1:8080/api/transfer
{"sender_id":10, "amount":100.5, "message":"hello"}

    Ответ: {"status":201,"body":{"message":"OK"}}

* Запрос: POST http://127.0.0.1:8080/api/transfer
{"recipient_id":11, "amount":500, "message":"hello"}
 
    Ответ: {"status":201,"body":{"message":"OK"}}

* Запрос: POST http://127.0.0.1:8080/api/transfer
{"sender_id":11,"recipient_id":5, "amount":1000, "message":"hello"}

    Ответ: {"status":201,"body":{"message":"OK"}}

* Запрос: POST http://127.0.0.1:8080/api/transfer
{"sender_id":11,"recipient_id":10, "amount":2000, "message":"for tea"}

    Ответ: {"status":409,"body":{"message":"Insufficient funds to write off"}}

* Запрос: POST http://127.0.0.1:8080/api/transfer
{"sender_id":35,"recipient_id":10, "amount":2000, "message":"for tea"}

    Ответ: {"status":409,"body":{"message":"Insufficient funds to write off"}}

*Запрос: POST http://127.0.0.1:8080/api/transfer
{"sender_id":11,"recipient_id":11, "amount":20, "message":"for tea"}

    Ответ: {"status":406,"body":{"message":"sender_id and recipient_id can't be zero at the same request"}}

#### Get /api/balance/<id>?currency=<> - метод для получения информации о балансе пользователя.

Примеры:

* Запрос: GET http://127.0.0.1:8080/api/balance/10

    Ответ: {"status":200,"body":{"balance":{"user_id":10,"total":899.5,"currency":"RUB"}}}
 
* Запрос: GET http://127.0.0.1:8080/api/balance/10?currency=USD

    Ответ: {"status":200,"body":{"balance":{"user_id":10,"total":11.57544719135,"currency":"USD"}}}

* Запрос: GET http://127.0.0.1:8080/api/balance/10?currency=EUR

    Ответ: {"status":200,"body":{"balance":{"user_id":10,"total":9.94967092245,"currency":"EUR"}}}
 
* Запрос: GET http://127.0.0.1:8080/api/balance/30
 
    Ответ: {"status":404,"body":{"message":"User's balance is not found"}}

#### Get /api/list/<id>?sort=<date/amount>&order=<desc/asc>&limit=<>&offset=<> - метод для получения списка транзакций пользователя.

Примеры:

* Запрос: GET http://127.0.0.1:8080/api/list/10

    Ответ: {"status":200,"body":{"list":[{"sender_id":10,"amount":100.5,"message":"hello","date":"2020-09-27T19:22:10+03:00"},{"recipient_id":10,"amount":1000,"message":"hello","date":"2020-09-27T19:21:15+03:00"}]}}
 
* Запрос: GET http://127.0.0.1:8080/api/list/10?order=asc&limit=1

    Ответ: {"status":200,"body":{"list":[{"recipient_id":10,"amount":1000,"message":"hello","date":"2020-09-27T19:21:15+03:00"}]}}

* Запрос: GET http://127.0.0.1:8080/api/list/10?sort=amount

    Ответ: {"status":200,"body":{"list":[{"recipient_id":10,"amount":1000,"message":"hello","date":"2020-09-27T19:21:15+03:00"},{"sender_id":10,"amount":100.5,"message":"hello","date":"2020-09-27T19:22:10+03:00"}]}}

* Запрос: GET http://127.0.0.1:8080/api/list/10?sort=amount&limit=1&offset=1

    Ответ: {"status":200,"body":{"list":[{"sender_id":10,"amount":100.5,"message":"hello","date":"2020-09-27T19:22:10+03:00"}]}}

* Запрос: GET http://127.0.0.1:8080/api/list/30

    Ответ: {"status":200,"body":{"list":[]}}