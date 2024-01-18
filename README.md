# Persons

## Запуск

1. Запустим субд и приложение

`docker-compose up`

2. Для накатывания миграций должен быть установлен goose:

```
go install github.com/pressly/goose/v3/cmd/goose@latest
```

3. Затем, чтобы накатить миграции выполни
```make up```

или 

```
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING=postgresql://postgres:example@127.0.0.1:5432/postgres?sslmode=disable

goose -dir db/migrations up
```

## Примеры запросов

### Запрос на создание

```bash
curl --location 'http://localhost:8080/person' \
--header 'Content-Type: application/json' \
--data '{
    "name": "Dmitriy",
    "surname": "Petrov"
}'
```

ответ 
```json
{
    "id": 14,
    "name": "Dmitriy",
    "surname": "Petrov",
    "nationality": "UA",
    "age": 43,
    "gender": "male"
}
```

### Запрос с фильтрами 

```bash
curl --location --globoff 'http://localhost:8080/person?page=1&filters[nationality]=RU'
```

```json
[
    {
        "id": 15,
        "name": "Dmitriy",
        "surname": "Petrov",
        "nationality": "UA",
        "age": 43,
        "gender": "male"
    }
]
```

### Удаление по идентификатору
```bash
curl --location --request DELETE 'http://localhost:8080/person/14'
```
