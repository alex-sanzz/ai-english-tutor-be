@echo off 

set DATABASE_URL=postgres://postgres:default@localhost:5432/ai_tutor?sslmode=disable

set COMMAND=%1

migrate -path ./migrations -database %DATABASE_URL% %*