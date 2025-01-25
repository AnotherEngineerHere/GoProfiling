@echo off
setlocal

if exist indexer.exe del indexer.exe
if exist server.exe del server.exe

go build -o indexer.exe ./cmd/indexer
go build -o server.exe ./cmd/server

:: Stop and remove existing zincsearch container if it exists
docker stop zincsearch 2>nul
docker rm zincsearch 2>nul

:: Start new zincsearch container
docker run -d --name zincsearch -p 4080:4080 ^
    -e ZINC_FIRST_ADMIN_USER=admin ^
    -e ZINC_FIRST_ADMIN_PASSWORD=admin ^
    public.ecr.aws/zinclabs/zinc:latest

:: Wait a few seconds for the container to start
timeout /t 5

indexer.exe "D:\Descargas\Programacion\TruoraSWE\enron_mail_20110403\maildir"

server.exe -port 3000