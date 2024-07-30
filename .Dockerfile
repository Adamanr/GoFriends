version: '3.9'

services:
  db:
    image: postgres
    restart: always
    ports:
        - 5431:5432
    environment:
      POSTGRES_PASSWORD: admin21
      POSTGRES_USER: postgres
      POSTGRES_DB: go_friends
  adminer:
    image: adminer
    restart: always
    ports:
      - 8091:8080