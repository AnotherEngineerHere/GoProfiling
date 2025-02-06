# Enron Email Search System
**Proyecto de Búsqueda de Correos**

## 🚀 Descripción
Sistema de búsqueda de correos electrónicos del dataset de Enron Corp utilizando un stack moderno de tecnologías.

## 🔧 Tecnologías
- **Backend**: Go
- **Frontend**: Vue 3
- **Base de Datos**: ZincSearch
- **Routing**: Chi
- **Estilos**: Tailwind CSS

## 💻 Requisitos Previos
- Go 1.20+
- Node.js 16+
- ZincSearch
- Docker (opcional)
- Base de datos
http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz (423MB)

## 🖥️ Instalación

### Backend
```bash
cd backend
go mod tidy
go build ./cmd/server
go build ./cmd/indexer
```

### Frontend
```bash
cd frontend
npm install
npm run dev
```

## ▶️ Uso

### Indexar Correos
```bash
./indexer /path/to/enron/dataset
```

### Iniciar Servidor
```bash
./server -port 3000
```

## ⚙️ Características
- Indexación de correos electrónicos
- Búsqueda full-text
- Interfaz web responsiva

## 📈 Optimizaciones
- Profiling detallado
- Mejoras de rendimiento documentadas

## ☁️ Despliegue
- Configuración de Terraform incluida
- Soporte para AWS/LocalStack

## 📄 Licencia
MIT License

## Contribuciones
Las contribuciones son bienvenidas. Por favor, lee las guías de contribución antes de enviar un PR.
