# Image Processing API

## Setup

```powershell
Copy-Item .env.example .env
docker compose up -d
.\.venv\Scripts\python.exe -m pip install -r requirements.txt
.\.venv\Scripts\alembic.exe upgrade head
.\.venv\Scripts\uvicorn.exe main:app --reload
```

API documentation: `http://127.0.0.1:8000/docs`

## Migrations

After changing SQLAlchemy models:

```powershell
.\.venv\Scripts\alembic.exe revision --autogenerate -m "describe change"
.\.venv\Scripts\alembic.exe upgrade head
```
