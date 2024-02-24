# goph-keeper

# Запуск клиента
Для начала клиент необходимо собрать командой `make build-client`.
Бинарник для запуска клиента будет находиться по пути `cmd/client/gclient`.

## Функции клиента

- login - функция авторизации на сервере. Необходима для получения токена.
- register - функция регистрации нового пользователя.
- logout - очистка пользовательского кэша и аутентификационных данных.
- records put [record_type] [path|data] [name] - отправка данных на сервер.
- records get [id] - получение данных с сервера, сохранение в кэш.
- records list - получение списка файлов с сервера.
- records sync - синхронизация данных между клиентом и сервером.

### Регистрация клиента
```
./cmd/client/bin/gclient register
```
Вводим Login
Вводим Password
- создается локальный json файл с конфигурацией, содержащий в себе логин, токен аутентификации, адрес API gophkeeper
- создается локальная папка для синхронизации записей с сервером

### Авторизация клиента для получения токена
```
./cmd/client/bin/gclient login
```
Вводим Login
Вводим Password

### Добавление данных
```
./cmd/client/bin/gclient records put pass login:password secretpassword
```

### Получение данных
```
./cmd/client/bin/gclient records list
```

### Получение записи
```
./cmd/client/bin/gclient records get secretpassword
```

### Синхронизация записей с сервера
При удалении локальной папки с записями можно воспользоваться командой `./cmd/client/bin/gclient records sync`
для восстановления записей с сервера.
```
./cmd/client/bin/gclient records sync
```

# Запуск сервера
Для начала сервер необходимо собрать командой `make build`.
Бинарник для запуска сервера будет находиться по пути `cmd/gophkeeper/gophkeeper$`.
Запустить клиент можно с помощью `make run`.