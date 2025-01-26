@echo off
setlocal

:: Crear directorio de perfiles si no existe
if not exist profiles mkdir profiles

:: Configurar variables de entorno para profiling
set CPU_PROFILE=%CD%\profiles\cpu.prof
set MEM_PROFILE=%CD%\profiles\mem.prof
set BLOCK_PROFILE=%CD%\profiles\block.prof
set TRACE_PROFILE=%CD%\profiles\trace.out

:: Compilar
go build -o indexer.exe ./cmd/indexer
go build -o server.exe ./cmd/server

:: Iniciar ZincSearch
docker stop zincsearch 2>nul
docker rm zincsearch 2>nul

docker run -d --name zincsearch -p 4080:4080 ^
    -e ZINC_FIRST_ADMIN_USER=admin ^
    -e ZINC_FIRST_ADMIN_PASSWORD=admin ^
    public.ecr.aws/zinclabs/zinc:latest

:: Esperar a que ZincSearch inicie
timeout /t 5

:: Ejecutar indexer con profiling
indexer.exe "D:\Descargas\Programacion\TruoraSWE\enron_mail_20110403\maildir"

:: Generar visualizaciones
start "" cmd /c "go tool pprof -http=:8080 %CPU_PROFILE%"
start "" cmd /c "go tool pprof -http=:8081 %MEM_PROFILE%"
start "" cmd /c "go tool pprof -http=:8082 %BLOCK_PROFILE%"
start "" cmd /c "go tool trace -http=:8083 %TRACE_PROFILE%"

:: Generar SVGs
go tool pprof -svg %CPU_PROFILE% > profiles/cpu.svg
go tool pprof -svg %MEM_PROFILE% > profiles/mem.svg
go tool pprof -svg %BLOCK_PROFILE% > profiles/block.svg

echo "Iniciando servidor"
:: Iniciar servidor en segundo plano
start "" server.exe -port 3000

pause