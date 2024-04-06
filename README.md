Команды для миграций:

    Создать файлы миграции: migrate create -ext sql -dir ./migrations -seq init
    Применить миграции: migrate -path ./migrations -database 'postgres://pguser:pgpassword@localhost:4445/postgres?sslmode=disable' up

POSTGRES_USER = pguser POSTGRES_DB = postgres POSTGRES_PASSWORD = pgpassword POSTGRES_HOST = localhost POSTGRES_PORT = 4445

<!-- 
POSTGRES_USER = pguser
POSTGRES_DB = postgres
POSTGRES_PASSWORD = pgpassword
POSTGRES_HOST = localhost
POSTGRES_PORT = 4445 -->