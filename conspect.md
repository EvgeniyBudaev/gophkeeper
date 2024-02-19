Инициализация зависимостей
```
go mod init github.com/EvgeniyBudaev/gophkeeper
```

Сборка
```
go build -v ./cmd/
```

Удаление неиспользуемых зависимостей
```
go mod tidy -v
```

Логирование
https://pkg.go.dev/go.uber.org/zap
```
go get -u go.uber.org/zap
```

Конфиг
dario.cat/mergo
```
go get -u dario.cat/mergo
```

ENV
```
go get github.com/caarlos0/env/v6
```

Подключение к БД
Драйвер для Postgres
```
go get -u github.com/lib/pq
```

GIN framework
```
go get -u github.com/gin-gonic/gin
```

JWT
```
go get -u github.com/golang-jwt/jwt/v4
```

Crypto
```
go get -u golang.org/x/crypto
```

pgconn
```
go get -u github.com/jackc/pgconn
```

spf13
```
go get -u github.com/spf13/viper
go get -u github.com/spf13/cobra
```

Миграции
https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md
https://www.appsloveworld.com/go/83/golang-migrate-installation-failing-on-ubuntu-22-04-with-the-following-gpg-error
```
curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
sudo sh -c 'echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list'
sudo apt-get update
sudo apt-get install -y golang-migrate
```

Если ошибка E: Указаны конфликтующие значения параметра Signed-By из источника
https://packagecloud.io/golang-migrate/migrate/ubuntu/
jammy: /etc/apt/keyrings/golang-migrate_migrate-archive-keyring.gpg !=
```
cd /etc/apt/sources.list.d
ls
sudo rm migrate.list
```

Создание миграционного репозитория
```
migrate create -ext sql -dir migrations UsersCreationMigration
migrate create -ext sql -dir migrations DataRecordsCreationMigration
```
