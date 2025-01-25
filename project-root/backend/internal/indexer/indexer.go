package indexer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"project/internal/repository"
)

type Email struct {
	Path      string    `json:"path"`
	Content   string    `json:"content"`
	Sender    string    `json:"sender"`
	Recipient string    `json:"recipient"`
	Subject   string    `json:"subject"`
	Date      time.Time `json:"date"`
	Folder    string    `json:"folder"`
}

type EmailIndexer struct {
	repo *repository.ZincRepository
}

func NewEmailIndexer(repo *repository.ZincRepository) *EmailIndexer {
	return &EmailIndexer{repo: repo}
}

type indexJob struct {
	filePath string
	content  string
}

func (e *EmailIndexer) IndexEmailsFromPath(path string) error {
	// Canal de trabajos y resultados
	jobs := make(chan indexJob, 100)
	results := make(chan indexResult, 100)

	// Número de workers basado en CPUs
	numWorkers := runtime.NumCPU()

	// Contadores con mutex para concurrencia segura
	var (
		indexedCount int
		errorCount   int
		skippedCount int
		mu           sync.Mutex
	)

	// Crear workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				result := e.processEmail(job)
				results <- result
			}
		}()
	}

	// Recolector de archivos
	go func() {
		defer close(jobs)

		filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}

			// Filtrar archivos
			ext := strings.ToLower(filepath.Ext(filePath))
			validExtensions := []string{".txt", "", ".eml"}
			isValidExtension := false
			for _, validExt := range validExtensions {
				if ext == validExt {
					isValidExtension = true
					break
				}
			}

			if isValidExtension {
				content, err := ioutil.ReadFile(filePath)
				if err == nil && len(content) > 0 {
					jobs <- indexJob{
						filePath: filePath,
						content:  string(content),
					}
				}
			}

			return nil
		})
	}()

	// Recolector de resultados
	go func() {
		wg.Wait()
		close(results)
	}()

	// Procesar resultados
	for result := range results {
		mu.Lock()
		if result.err != nil {
			errorCount++
			if strings.Contains(result.err.Error(), "vacío") {
				skippedCount++
			}
		} else {
			indexedCount++

			// Imprimir progreso cada 100 correos
			if indexedCount%100 == 0 {
				fmt.Printf("Indexados %d correos...\n", indexedCount)
			}
		}
		mu.Unlock()
	}

	fmt.Printf("Indexación completada. Total de correos indexados: %d, Errores: %d, Archivos vacíos: %d\n",
		indexedCount, errorCount, skippedCount)

	return nil
}

type indexResult struct {
	email Email
	err   error
}

func (e *EmailIndexer) processEmail(job indexJob) indexResult {
	// Verificar si el archivo está vacío
	if len(job.content) == 0 {
		return indexResult{err: fmt.Errorf("archivo vacío")}
	}

	// Parsear email
	email, err := e.parseEmail(job.filePath, job.content)
	if err != nil {
		return indexResult{err: err}
	}

	// Indexar en ZincSearch
	err = e.repo.Index("emails", email)
	if err != nil {
		return indexResult{err: err}
	}

	return indexResult{email: email}
}

func (e *EmailIndexer) parseEmail(filePath string, contentStr string) (Email, error) {
	// Verificar si el contenido está vacío o no es un correo válido
	if !strings.Contains(contentStr, "From:") || !strings.Contains(contentStr, "To:") {
		return Email{}, fmt.Errorf("contenido no parece ser un correo válido")
	}

	// Parseo mejorado de fecha
	dateStr := extractHeader(contentStr, "Date:")
	var parsedDate time.Time
	if dateStr != "" {
		// Intentar múltiples formatos de fecha
		formats := []string{
			time.RFC1123Z,
			time.RFC822,
			time.RFC850,
			"Mon, 2 Jan 2006 15:04:05 -0700",
			"2 Jan 2006 15:04:05 -0700",
		}

		for _, format := range formats {
			if date, err := time.Parse(format, dateStr); err == nil {
				parsedDate = date
				break
			}
		}
	}

	sender := extractHeader(contentStr, "From:")
	recipient := extractHeader(contentStr, "To:")
	subject := extractHeader(contentStr, "Subject:")

	// Obtener carpeta relativa
	relativeFolder, err := filepath.Rel(filepath.Dir(filePath), filepath.Dir(filePath))
	if err != nil {
		relativeFolder = filepath.Base(filepath.Dir(filePath))
	}

	return Email{
		Path:      filePath,
		Content:   contentStr,
		Sender:    sender,
		Recipient: recipient,
		Subject:   subject,
		Date:      parsedDate,
		Folder:    relativeFolder,
	}, nil
}

// Función de utilidad para extraer encabezados con más robustez
func extractHeader(content, header string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, header) {
			return strings.TrimSpace(strings.TrimPrefix(line, header))
		}
	}
	return ""
}
