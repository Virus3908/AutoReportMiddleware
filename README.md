# 🛠 Middleware-сервис на Go

Этот репозиторий содержит **middleware-сервис**, реализующий обработку задач через Kafka. Сервис принимает сообщения, обрабатывает их и отправляет результаты обратно — он работает в фоновом режиме и является частью распределённой системы.

> ⚠️ **Важно**: Это не самостоятельный проект — требуется связка с backend-сервисом и frontend-сервисом

---

## 📦 Возможности

- Получение задач из Kafka (`task` топик)<br/>
- Обработка задач: преобразование, диаризация, транскрипция, генерация отчётов<br/>
- Отправка результатов в Kafka (`callback` топик)<br/>
- Интеграция с S3 и PostgreSQL<br/>
- Логгирование с использованием `zap`<br/>
- Конфигурация через YAML + `.env`<br/>

---

## 🚀 Быстрый старт

### 1. Установи зависимости

```bash
go mod tidy

2. Создай .env

PG_HOST=127.0.0.1
PG_PORT=5432
PG_USER=virus
PG_PASSWORD=postgres
PG_DATABASE=db

S3_ACCESS_KEY=virus
S3_SECRET_KEY=password
S3_BUCKET=bucket
S3_REGION=ru-central1
S3_ENDPOINT=http://10.66.66.7:9000
S3_PUBLIC_BASE_URL=http://10.66.66.7:9000
S3_MINIO=true

KAFKA_BROKERS=10.66.66.7:9092
KAFKA_GROUP_ID=go_consumer
```

3. Укажи путь к конфигу
```bash
export CONFIG_PATH=./config/config.yaml
```

4. Создай базу данных в Postgres
Создай базу и в `Makefile` укажи адрес
```
DATABASE_URL=postgres://user:password@host:5432/database_name?sslmode=disable
```
И после этого сделай
```bash
make migrate-up
```

⸻

🐳 Сборка и запуск через Docker
```bash
make docker-build
make docker-run
```
Или вручную:
```bash
docker build -t middleware-app .
docker run --rm -it --env-file .env middleware-app
```

⸻

🧾 Makefile команды<br/>
	•	make generate — генерация gRPC и SQL-кода<br/>
	•	make tidy — go mod tidy<br/>
	•	make clean — удаление сгенерированных файлов<br/>
	•	make docker-build — сборка Docker-образа<br/>
	•	make docker-run — запуск контейнера с переменными из .env<br/>
    •	make migrate-up — создание базы данных<br/>

⸻

# 🔗-зависимости

Этот сервис не будет работать сам по себе. Требуется backend-сервис (Python), который:<br/>
	•	Принимает задачи в Kafka<br/>
	•	Отправляет результаты в Kafka<br/>

Так же необходимо поднять Kafka и Minio.
Так же рекомендую создать свою сеть в docker
👉 Рекомендуется использовать с репозиторием backend-сервиса (https://github.com/Virus3908/AutoReporterBack).<br/>
👉 Рекомендуется использовать с репозиторием frontend-сервиса (https://github.com/Virus3908/AutoReportFrontend).<br/>

⸻

⚙️ Структура проекта
```bash
.
├── cmd/                  # Входная точка приложения
├── internal/             # Основная логика приложения
│   ├── config/           # Загрузка конфигурации
│   ├── kafka/            # Консьюмер/Продюсер
│   ├── postgres/         # Подключение к БД
│   ├── repositories/     # SQL-слой
│   ├── services/         # Логика обработки задач
│   └── logger/           # Интерфейс логгера (zap)
├── pkg/messages/         # gRPC-сообщения
├── config/config.yaml    # Конфигурационный файл
├── .env                  # Переменные окружения
└── Makefile              # Сборка проекта
```

⸻

📋 TODO <br/>
	•	Покрыть unit-тестами<br/>
	•	Автоматическая CI/CD-сборка<br/>
	•	Документация по backend API<br/>
	•	Healthcheck эндпоинты<br/>

⸻

🧑‍💻 Авторы

Разработано Virus3908. Вопросы и предложения — через issue tracker или напрямую.