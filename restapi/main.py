from fastapi import FastAPI, HTTPException, Depends
from sqlalchemy.orm import Session
from schema import Fruit, FruitCreate
from models import FruitModel
from database import get_db, Base, engine

app = FastAPI()

Base.metadata.create_all(bind=engine)

@app.get("/")
def read_root():
    return {"Hello": "World"}

@app.get("/fruits/", response_model=list[Fruit])
def read_fruits(db: Session = Depends(get_db)):
    return db.query(FruitModel).all()

@app.get("/fruits/{fruit_id}", response_model=Fruit)
def read_fruit(fruit_id: int, db: Session = Depends(get_db)):
    fruit = db.query(FruitModel).filter(FruitModel.id == fruit_id).first()
    if not fruit:
        raise HTTPException(status_code=404, detail="Fruit not found.")
    return fruit

@app.post("/fruits/", response_model=Fruit)
def create_fruit(fruit: FruitCreate, db: Session = Depends(get_db)):
    db_fruit = FruitModel(name=fruit.name, description=fruit.description)
    db.add(db_fruit)
    db.commit()
    db.refresh(db_fruit)
    return db_fruit

@app.put("/fruits/{fruit_id}", response_model=Fruit)
def update_fruit(fruit_id: int, fruit: FruitCreate, db: Session = Depends(get_db)):
    db_fruit = db.query(FruitModel).filter(FruitModel.id == fruit_id).first()
    if not db_fruit:
        raise HTTPException(status_code=404, detail="Fruit not found.")
    db_fruit.name = fruit.name
    db_fruit.description = fruit.description
    db.commit()
    db.refresh(db_fruit)
    return db_fruit

@app.delete("/fruits/{fruit_id}")
def delete_fruit(fruit_id: int, db: Session = Depends(get_db)):
    db_fruit = db.query(FruitModel).filter(FruitModel.id == fruit_id).first()
    if not db_fruit:
        raise HTTPException(status_code=404, detail="Fruit not found.")
    db.delete(db_fruit)
    db.commit()
    return {"message": f"Fruit {fruit_id} ({db_fruit.name}) deleted successfully.", "fruit": Fruit.model_validate(db_fruit)}
