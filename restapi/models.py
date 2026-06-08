from sqlalchemy import Column, Integer, String
from database import Base

class FruitModel(Base):
    __tablename__ = "fruits"

    id = Column(Integer, primary_key=True, index=True)
    name = Column(String, nullable=False)
    description = Column(String, nullable=True)
