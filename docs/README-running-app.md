# Запуск приложения локально

## 1. Поднять базу данных

```bash
docker-compose up -d postgres
```

## 2. Запустить backend

```bash
cd backend
go run ./cmd/server -v
```

## 3. Запустить frontend

```bash
cd frontend
npm install
npm run dev
```

Открыть приложение: http://localhost:3000

## 4. Полная проверка перед коммитом

```bash
cd frontend
npm run build
npm run typecheck
npm test
```

## 5. Полезные заметки

- База данных доступна на порту 5432.
- Backend слушает API на порту 8080.
- Frontend запускается на порту 3000.
- Для локальной разработки удобно сначала поднять postgres, затем backend, затем frontend.
