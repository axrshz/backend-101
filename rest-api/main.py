from fastapi import FastAPI, HTTPException, Depends, status
from fastapi.security import OAuth2PasswordRequestForm
from datetime import timedelta
from sqlalchemy.orm import Session
from schema import Fruit, FruitCreate, User, UserCreate, Token
from models import FruitModel, UserModel
from database import get_db, Base, engine
from auth import (
    ACCESS_TOKEN_EXPIRE_MINUTES,
    get_password_hash,
    verify_password,
    create_access_token,
    get_current_user,
    require_admin,
)

app = FastAPI()

Base.metadata.create_all(bind=engine)

@app.get("/")
def read_root():
    return {"Hello": "World"}

#Authentication endpoints
@app.post("/auth/register", response_model=User)
def register_user(user: UserCreate, db: Session = Depends(get_db)):
    db_user = db.query(UserModel).filter(UserModel.username == user.username).first()
    if db_user:
        raise HTTPException(status_code=400, detail="Username already registered.")
    hashed_password = get_password_hash(user.password)
    new_user = UserModel(username=user.username, email=user.email, hashed_password=hashed_password)
    db.add(new_user)
    db.commit()
    db.refresh(new_user)
    return new_user

@app.post("/auth/login", response_model=Token)
def login_for_access_token(
    form_data: OAuth2PasswordRequestForm = Depends(),
    db: Session = Depends(get_db),
):
    user = db.query(UserModel).filter(UserModel.username == form_data.username).first()
    if not user or not verify_password(form_data.password, user.hashed_password):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect username or password.",
            headers={"WWW-Authenticate": "Bearer"},
        )
    access_token_expires = timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
    access_token = create_access_token(
        data={"sub": user.username}, expires_delta=access_token_expires
    )
    return {"access_token": access_token, "token_type": "bearer"}

@app.get("/fruits/", response_model=list[Fruit])
def read_fruits(
    db: Session = Depends(get_db),
    current_user: UserModel = Depends(get_current_user),
):
    return db.query(FruitModel).all()


@app.get("/fruits/{fruit_id}", response_model=Fruit)
def read_fruit(
    fruit_id: int,
    db: Session = Depends(get_db),
    current_user: UserModel = Depends(get_current_user),
):
    fruit = db.query(FruitModel).filter(FruitModel.id == fruit_id).first()
    if not fruit:
        raise HTTPException(status_code=404, detail="Fruit not found.")
    return fruit

@app.post("/fruits/", response_model=Fruit)
def create_fruit(fruit: FruitCreate, db: Session = Depends(get_db), current_user: UserModel = Depends(require_admin)):
    db_fruit = FruitModel(name=fruit.name, description=fruit.description)
    db.add(db_fruit)
    db.commit()
    db.refresh(db_fruit)
    return db_fruit

@app.put("/fruits/{fruit_id}", response_model=Fruit)
def update_fruit(fruit_id: int, fruit: FruitCreate, db: Session = Depends(get_db), current_user: UserModel = Depends(require_admin)):
    db_fruit = db.query(FruitModel).filter(FruitModel.id == fruit_id).first()
    if not db_fruit:
        raise HTTPException(status_code=404, detail="Fruit not found.")
    db_fruit.name = fruit.name
    db_fruit.description = fruit.description
    db.commit()
    db.refresh(db_fruit)
    return db_fruit

@app.delete("/fruits/{fruit_id}")
def delete_fruit(fruit_id: int, db: Session = Depends(get_db), current_user: UserModel = Depends(require_admin)):
    db_fruit = db.query(FruitModel).filter(FruitModel.id == fruit_id).first()
    if not db_fruit:
        raise HTTPException(status_code=404, detail="Fruit not found.")
    db.delete(db_fruit)
    db.commit()
    return {"message": f"Fruit {fruit_id} ({db_fruit.name}) deleted successfully.", "fruit": Fruit.model_validate(db_fruit)}
