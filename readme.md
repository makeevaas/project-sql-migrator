- Утилита sql-мигратор для управления миграциями баз данных


Доступные команды:
- create - создание нового файла для описания миграции
- up - накатывает миграцию
- down - откатывает миграцию
- redo - повторяет последнюю миграцию
- status - выводит статус всех миграций в консоль

Настройка среды для запуска:

- Необходимое ПО для работы приложения
    - Docker
    - Go 1.24
- Запустить бд с postgres 
    - скачивание образа: 
        - docker pull postgres:17.0
    - запуск контейнера: 
        - make postgres
    - просмотр данных через консоль
        - docker exec -it ($containerID)  psql -U postgres -w main_db
- Сборка утилиты
    - make build
- Необходимые переменные среды для работы приложения:
    - DB_CONNECTION_PATH - dsn для подключения к БД
        - export DB_CONNECTION_PATH="postgres://postgres:pwd@localhost:5432/main_db?sslmode=disable"
    - MIGRATIONS_PATH - директория хранения файлов миграций
        - export MIGRATIONS_PATH=./migrations
- Запуск команд утилиты
    - make migrate_create
    - make migrate_up     
    - make migrate_down
    - make migrate_redo
    - make migrate_status            

- Запуск тестов
    - make tests    