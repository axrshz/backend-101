from fastapi import FastAPI

from app.schemas import HealthResponse

app = FastAPI(title="Image Processing API")


@app.get("/", response_model=HealthResponse)
async def read_root() -> HealthResponse:
    return HealthResponse(status="ok")
