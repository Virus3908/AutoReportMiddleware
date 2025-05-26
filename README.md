Конечно! Вот базовый README.md для Go backend части твоего проекта. Он оформлен в стиле продакшн-ориентированного репозитория без использования Docker, с поддержкой make, логгированием, Kafka, S3, PostgreSQL и Wire.

⸻


# 🚀 AI-Report Middleware (Go)

Это серверная часть проекта на Go, отвечающая за маршрутизацию задач, работу с Kafka, PostgreSQL и хранилищем файлов (S3), а также за обработку задач через сервисный слой.

---

## 📦 Используемые технологии

- **Go 1.21+**
- **Kafka (consumer + producer)**
- **PostgreSQL (через pgx)**
- **S3 (Minio / AWS)**
- **Zap logger**
- **Wire (dependency injection)**
- **mux (router)**
- **Protobuf (gRPC совместимость)**

---

## ⚙️ Быстрый старт

### 1. Установи зависимости

```bash
go mod tidy
```

2. Подготовь переменные окружения

Создай файл .env и укажи:

DB_USER=postgres
DB_PASSWORD=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=ai_report

AWS_ACCESS_KEY_ID=your_key
AWS_SECRET_ACCESS_KEY=your_secret
S3_BUCKET_NAME=my-bucket
S3_ENDPOINT_URL=http://localhost:9000
S3_PUBLIC_BASE_URL=http://localhost:9000

KAFKA_BROKERS=localhost:9092
KAFKA_GROUP_ID=ai-go-group

3. Запусти приложение

make run


⸻

🛠 Makefile команды

make run          # запуск приложения
make wire         # пересобрать зависимости через wire
make build        # сборка бинарника
make lint         # запуск linters (если подключено)


⸻

🧾 Структура проекта

cmd/
└── app/           # Точка входа в приложение

internal/
├── config/        # Загрузка конфигураций
├── handlers/      # HTTP-обработчики
├── kafka/         # Producer / Consumer
├── logger/        # Zap логгер
├── postgres/      # БД клиент
├── repositories/  # Работа с pgx и таблицами
├── services/      # Сервисный слой (TaskDispatcher, и т.д.)
└── storage/       # Работа с S3

proto/
└── messages.proto # Протокол обмена


⸻

✅ Особенности
	•	Используется транзакционный менеджер для операций с PostgreSQL.
	•	Логгер (zap) инкапсулирован, может быть заменён на любой другой через интерфейс.
	•	Обработка задач Convert, Diarize, Transcribe, SemiReport, Report реализована в TaskDispatcher.
	•	Kafka consumer читает сообщения и обрабатывает их асинхронно, offset коммитится только после успеха.

⸻

📄 Генерация wire

make wire

Wire используется для сборки зависимостей (Logger, Services, Handlers, и т.д.).

⸻

📌 Заметки
	•	Компоненты с контекстом (context.Context) пробрасываются строго.
	•	Kafka используется как основная очередь, без блокировок.
	•	PostgreSQL использует pgx + connection pooling.

⸻

📄 Лицензия

MIT