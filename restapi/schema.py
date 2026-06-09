from pydantic import BaseModel
from typing import Optional

class FruitBase(BaseModel):
    name: str
    description: Optional[str] = None

class FruitCreate(FruitBase):
    """Schema for creating or updating a fruit (no id — the DB decides it)."""
    pass

class Fruit(FruitBase):
    """Schema for reading a fruit (includes the auto-generated id)."""
    id: int

    class Config:
        from_attributes = True


class UserBase(BaseModel):
    username: str
    email: str
    is_admin: bool = False

class UserCreate(UserBase):
    password: str

class User(UserBase):
    id: int
    is_active: bool

    class Config:
        from_attributes = True


class Token(BaseModel):
    access_token: str
    token_type: str

class TokenData(BaseModel):
    username: str | None = None