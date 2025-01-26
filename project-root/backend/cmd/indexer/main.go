package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"time"

	"project/internal/indexer"
	"project/internal/repository"
)

func ensureDirectoryExists(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0755)
}

func startServer() {
	// Obtener el directorio actual
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("No se pudo obtener el directorio actual: %v", err)
	}

	// Construir la ruta al directorio del servidor
	serverDir := filepath.Join(currentDir, "cmd", "server")

	// Crear un comando para ejecutar el servidor
	cmd := exec.Command("go", "run", filepath.Join(serverDir, "main.go"), "-port", "3000")

	// Configurar directorios de trabajo y ambiente
	cmd.Dir = serverDir
	cmd.Env = append(os.Environ(), "GO111MODULE=on")

	// Configurar pipes para capturar salida
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Error creando pipe de salida: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("Error creando pipe de error: %v", err)
	}

	// Iniciar goroutines para manejar la salida
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			log.Printf("Servidor (stdout): %s", scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			log.Printf("Servidor (stderr): %s", scanner.Text())
		}
	}()

	// Iniciar el servidor
	err = cmd.Start()
	if err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}

	// Esperar que el servidor termine (en segundo plano)
	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Printf("Servidor terminó con error: %v", err)
		}
	}()
}

func main() {
	// Configurar logging detallado
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	fmt.Println("Iniciando indexador con profiling...")

	// Crear directorio de perfiles
	os.MkdirAll("profiles", os.ModePerm)

	// Verificar si se proporcionó la ruta del dataset
	if len(os.Args) < 2 {
		log.Fatal("Por favor proporciona la ruta del dataset")
	}
	datasetPath := os.Args[1]
	fmt.Printf("Ruta del dataset: %s\n", datasetPath)

	// Verificar existencia del directorio
	if _, err := os.Stat(datasetPath); os.IsNotExist(err) {
		log.Fatalf("El directorio no existe: %s", datasetPath)
	}

	// Rutas de perfiles
	cpuProfile := "profiles/cpu.prof"
	memProfile := "profiles/mem.prof"
	blockProfile := "profiles/block.prof"
	traceProfile := "profiles/trace.out"

	// Profiling de CPU
	cpuFile, err := os.Create(cpuProfile)
	if err != nil {
		log.Fatalf("No se pudo crear archivo de perfil de CPU: %v", err)
	}
	defer cpuFile.Close()

	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		log.Fatalf("No se pudo iniciar perfil de CPU: %v", err)
	}
	defer pprof.StopCPUProfile()

	// Traza de ejecución
	traceFile, err := os.Create(traceProfile)
	if err != nil {
		log.Fatalf("No se pudo crear archivo de traza: %v", err)
	}
	defer traceFile.Close()

	if err := trace.Start(traceFile); err != nil {
		log.Fatalf("No se pudo iniciar traza: %v", err)
	}
	defer trace.Stop()

	// Habilitar profiling de bloqueo
	runtime.SetBlockProfileRate(1)

	// Información de sistema
	fmt.Printf("Número de CPUs: %d\n", runtime.NumCPU())
	fmt.Printf("Número de goroutines inicial: %d\n", runtime.NumGoroutine())

	// Tiempo de inicio
	startTime := time.Now()

	// Configurar repositorio de ZincSearch
	fmt.Println("Inicializando repositorio ZincSearch...")
	zincRepo, err := repository.NewZincRepository()
	if err != nil {
		log.Fatalf("Error inicializando ZincSearch: %v", err)
	}

	// Crear indexador
	fmt.Println("Creando indexador de correos...")
	idx := indexer.NewEmailIndexer(zincRepo)

	// Indexar correos
	fmt.Println("Comenzando indexación de correos...")
	err = idx.IndexEmailsFromPath(datasetPath)
	if err != nil {
		log.Fatalf("Error indexando correos: %v", err)
	}

	// Tiempo total de indexación
	indexTime := time.Since(startTime)
	fmt.Printf("Tiempo total de indexación: %v\n", indexTime)

	// Información de goroutines al final
	fmt.Printf("Número de goroutines al final: %d\n", runtime.NumGoroutine())

	// Generar perfil de memoria
	memFile, err := os.Create(memProfile)
	if err != nil {
		log.Fatalf("No se pudo crear archivo de perfil de memoria: %v", err)
	}
	defer memFile.Close()

	runtime.GC() // Forzar recolección de basura
	if err := pprof.WriteHeapProfile(memFile); err != nil {
		log.Fatalf("No se pudo escribir perfil de memoria: %v", err)
	}

	// Generar perfil de bloqueo
	blockFile, err := os.Create(blockProfile)
	if err != nil {
		log.Fatalf("No se pudo crear archivo de perfil de bloqueo: %v", err)
	}
	defer blockFile.Close()

	if err := pprof.Lookup("block").WriteTo(blockFile, 0); err != nil {
		log.Fatalf("No se pudo escribir perfil de bloqueo: %v", err)
	}

	fmt.Println("Proceso de indexación y profiling completado.")

	// Iniciar servidor
	fmt.Println("Iniciando servidor...")
	startServer()

	// Esperar para mantener el proceso principal activo
	select {}
}
