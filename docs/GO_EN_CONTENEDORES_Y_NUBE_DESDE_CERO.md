# Curso Completo: Go en Contenedores y Nube (desde cero)

## 0) A quien va dirigido

Este documento es para ti si:

- No sabes Go o estas comenzando.
- Quieres ejecutar aplicaciones Go dentro de contenedores.
- Quieres entender como llevarlas a la nube de forma profesional.
- Quieres aprender la logica, sintaxis y arquitectura real para proyectos cloud.

La meta es que, al final, puedas:

1. Construir una API o servidor web en Go.
2. Empaquetarlo en Docker.
3. Ejecutarlo localmente y en cloud.
4. Configurarlo con variables de entorno.
5. Monitorearlo, asegurarlo y escalarlo.

---

## 1) Fundamentos minimos de Go para nube

### 1.1 Que es Go y por que se usa tanto en cloud

Go es un lenguaje compilado orientado a simplicidad, performance y mantenibilidad.
Es muy usado en cloud por estas razones:

1. Compila a binarios unicos, faciles de desplegar.
2. Arranca rapido y consume poca memoria.
3. Concurrencia simple con goroutines.
4. Libreria estandar excelente para red y HTTP.
5. Facil de contenerizar en imagenes pequenas.

Herramientas cloud famosas hechas en Go: Docker, Kubernetes, Terraform, Prometheus.

### 1.2 Sintaxis minima que debes dominar

- Paquete principal:

```go
package main
```

- Punto de entrada:

```go
func main() {
    // inicia aqui
}
```

- Variables:

```go
nombre := "api-go"
var puerto string = "8080"
```

- Manejo de errores (estilo Go):

```go
valor, err := algunaFuncion()
if err != nil {
    // manejar error
    return
}
```

- Struct para agrupar datos:

```go
type Config struct {
    Port string
    Env  string
}
```

Este patron es clave en apps cloud: agrupar configuracion en structs y validarla al iniciar.

---

## 2) De app local a app cloud: el mapa mental

Piensa en 5 etapas:

1. Escribir app Go.
2. Construir binario.
3. Empaquetar en imagen Docker.
4. Ejecutar contenedor local.
5. Subir imagen a registro y desplegar en nube.

Flujo resumido:

Codigo Go -> go build -> Docker image -> Registry -> Cloud runtime

---

## 3) Proyecto base de servidor HTTP en Go

### 3.1 Ejemplo minimo de servidor

```go
package main

import (
    "log"
    "net/http"
    "os"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ok"))
    })

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        _, _ = w.Write([]byte("Hola desde Go en contenedor"))
    })

    addr := ":" + port
    log.Printf("server listening on %s", addr)
    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatal(err)
    }
}
```

### 3.2 Que conceptos cloud ya aparecen aqui

1. Uso de variable PORT desde entorno.
2. Endpoint de health check.
3. Logs de inicio.
4. Puerto configurable (no fijo).

Esto es obligatorio para muchas plataformas gestionadas (Cloud Run, Heroku, etc.).

---

## 4) Go modules: dependencia y versionado

### 4.1 Que es go.mod

El archivo go.mod define:

- Nombre del modulo.
- Version de Go.
- Dependencias y versiones.

Ejemplo:

```go
module laboratorio-nube-go

go 1.22
```

### 4.2 Comandos basicos

- Inicializar modulo:

```bash
go mod init mi-modulo
```

- Agregar dependencias segun imports:

```bash
go mod tidy
```

- Compilar:

```bash
go build .
```

En contenedores, un go.mod limpio mejora cache y velocidad de build.

---

## 5) Docker para Go: fundamentos esenciales

### 5.1 Conceptos clave

1. Imagen: plantilla inmutable.
2. Contenedor: instancia en ejecucion de una imagen.
3. Dockerfile: receta para construir imagen.
4. Registry: repositorio remoto de imagenes.

### 5.2 Por que usar multi-stage build

Porque separa:

- Etapa builder (compilar)
- Etapa runtime (solo ejecutar)

Beneficios:

1. Imagen final mas pequena.
2. Menor superficie de ataque.
3. Deploy mas rapido.

### 5.3 Dockerfile recomendado para Go

```dockerfile
# Stage 1: build
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Copiar mod files primero para aprovechar cache
COPY go.mod go.sum* ./
RUN go mod download

# Copiar el resto del codigo
COPY . .

# Compilar binario estatico para linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

# Stage 2: runtime minima
FROM alpine:3.20
WORKDIR /app

# Opcional: certificados para llamadas HTTPS salientes
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server /app/server
COPY templates /app/templates
COPY static /app/static
COPY imagenes /app/imagenes

EXPOSE 8000
CMD ["/app/server"]
```

Nota practica:

1. Este `CMD` asume que tu app toma `PORT` desde variable de entorno.
2. Si tu app aun usa solo flags (`-port`, `-dir`), puedes mantener temporalmente:

```dockerfile
CMD ["/app/server", "-port=8000", "-dir=imagenes"]
```

3. Para cloud moderna, conviene soportar ambos: primero `PORT` por entorno y luego fallback a flag.

### 5.4 Comandos Docker basicos

Construir:

```bash
docker build -t laboratorio-go:latest .
```

Ejecutar:

```bash
docker run --rm -p 8000:8000 laboratorio-go:latest
```

Montar carpeta local de imagenes:

```bash
docker run --rm -p 8000:8000 -v "${PWD}/imagenes:/app/imagenes" laboratorio-go:latest
```

---

## 6) Configuracion cloud-friendly en Go

### 6.1 Nunca hardcodear secretos o endpoints

Mal:

```go
apiKey := "123456"
```

Bien:

```go
apiKey := os.Getenv("API_KEY")
if apiKey == "" {
    log.Fatal("API_KEY is required")
}
```

### 6.2 Struct de configuracion

```go
type AppConfig struct {
    Port    string
    Env     string
    LogJSON bool
}

func loadConfig() AppConfig {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    env := os.Getenv("APP_ENV")
    if env == "" {
        env = "dev"
    }

    return AppConfig{
        Port:    port,
        Env:     env,
        LogJSON: os.Getenv("LOG_JSON") == "true",
    }
}
```

### 6.3 Validacion de config al arrancar

La app debe fallar rapido si falta configuracion critica.

---

## 7) Rutas y arquitectura para servicios cloud

### 7.1 Rutas minimas recomendadas

1. /health para saber si proceso esta vivo.
2. /ready para saber si esta listo para trafico.
3. /metrics para monitoreo (si usas Prometheus).
4. / endpoint principal.

### 7.2 Diferencia health vs ready

- health: proceso esta corriendo.
- ready: dependencias listas (db, cache, etc.).

### 7.3 Ejemplo basico

```go
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte("ok"))
})

http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
    // aqui validarias db u otros servicios
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte("ready"))
})
```

---

## 8) Logging y observabilidad

### 8.1 Logging util en nube

En cloud, los logs son tu ventana de debugging. Reglas:

1. Logs estructurados (ideal JSON).
2. Incluir nivel (info, warn, error).
3. Incluir request id si existe.
4. Nunca loggear secretos.

Ejemplo simple:

```go
log.Printf("level=info msg=server_started port=%s", cfg.Port)
```

### 8.2 Metricas

Mide al menos:

1. Numero de requests.
2. Latencia por endpoint.
3. Errores por status code.
4. Uso de memoria y CPU.

### 8.3 Trazas

Con OpenTelemetry puedes seguir una request entre servicios.

---

## 9) Manejo de apagado elegante (graceful shutdown)

En contenedores, cuando la plataforma quiere detener tu app, envia una senal. Debes cerrar bien conexiones para no cortar requests en curso.

Ejemplo:

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        _, _ = w.Write([]byte("ok"))
    })

    srv := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    go func() {
        log.Println("server started")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
    <-stop

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("shutdown error: %v", err)
    }

    log.Println("server stopped gracefully")
}
```

---

## 10) Concurrencia en Go aplicada a nube

### 10.1 Goroutines

Una goroutine es una funcion que corre concurrentemente:

```go
go func() {
    // tarea en paralelo
}()
```

### 10.2 Channels

Permiten comunicar goroutines de forma segura.

```go
ch := make(chan string)

go func() {
    ch <- "hecho"
}()

msg := <-ch
```

### 10.3 Cuidados

1. Evitar data races con memoria compartida.
2. Usar context para cancelar operaciones largas.
3. Limitar concurrencia para no saturar CPU o APIs externas.

---

## 11) Seguridad basica para Go en contenedores

### 11.1 Principios obligatorios

1. No correr contenedor como root si no hace falta.
2. Mantener imagen minima.
3. Escanear vulnerabilidades de imagen.
4. No exponer secretos en codigo ni logs.
5. Validar entradas HTTP.

### 11.2 Ajuste de usuario no root (Dockerfile)

```dockerfile
FROM alpine:3.20
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /app/server /app/server
USER appuser
CMD ["/app/server"]
```

### 11.3 Headers de seguridad HTTP

```go
func securityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        next.ServeHTTP(w, r)
    })
}
```

---

## 12) Performance para Go cloud

### 12.1 Buenas practicas

1. Reusar conexiones HTTP salientes con cliente compartido.
2. Definir timeouts en servidor y cliente.
3. Evitar allocs innecesarias en loops criticos.
4. Cargar configuracion una sola vez al inicio.

### 12.2 Timeouts recomendados en servidor

```go
srv := &http.Server{
    Addr:              ":8080",
    Handler:           mux,
    ReadHeaderTimeout: 5 * time.Second,
    ReadTimeout:       10 * time.Second,
    WriteTimeout:      15 * time.Second,
    IdleTimeout:       60 * time.Second,
}
```

Esto protege contra conexiones lentas o abusivas.

---

## 13) CI/CD para Go en contenedores

### 13.1 Pipeline ideal

1. go test
2. go vet
3. build Docker image
4. scan de seguridad
5. push a registry
6. deploy a cloud

### 13.2 Comandos tipicos

```bash
go test ./...
go vet ./...
docker build -t mi-registry/mi-app:1.0.0 .
docker push mi-registry/mi-app:1.0.0
```

---

## 14) Deploy en plataformas cloud (vista general)

### 14.1 Cloud Run (GCP)

- Despliegas imagen.
- Debe escuchar PORT definido por entorno.
- Escala automatico por request.

### 14.2 AWS ECS/Fargate

- Definicion de tarea con contenedor.
- Balanceador y auto scaling.
- Variables y secretos en task definition.

### 14.3 Kubernetes

Objetos minimos:

1. Deployment (replicas y version).
2. Service (red interna/externa).
3. ConfigMap (config no secreta).
4. Secret (credenciales).

---

## 15) Kubernetes basico para app Go

### 15.1 Deployment ejemplo

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-app
spec:
  replicas: 2
  selector:
    matchLabels:
      app: go-app
  template:
    metadata:
      labels:
        app: go-app
    spec:
      containers:
        - name: go-app
          image: mi-registry/go-app:1.0.0
          ports:
            - containerPort: 8080
          env:
            - name: PORT
              value: "8080"
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
          readinessProbe:
            httpGet:
              path: /ready
              port: 8080
```

### 15.2 Service ejemplo

```yaml
apiVersion: v1
kind: Service
metadata:
  name: go-app-service
spec:
  selector:
    app: go-app
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: LoadBalancer
```

---

## 16) Almacenamiento y estado

Regla cloud:

- Trata tus contenedores como efimeros.
- No dependas del disco local para datos criticos.

Usa:

1. Base de datos administrada.
2. Object storage para archivos (S3, GCS).
3. Cache distribuida (Redis) para estado temporal.

En tu proyecto actual, la carpeta imagenes local funciona para laboratorio. En produccion conviene usar bucket y URLs firmadas o streaming.

---

## 17) Secretos y configuracion sensible

No guardes secretos en:

- repositorio
- Dockerfile
- logs

Opciones seguras:

1. Variables de entorno inyectadas por plataforma.
2. Secret managers (AWS Secrets Manager, GCP Secret Manager, Vault).
3. Kubernetes Secrets.

Patron recomendado:

```go
dbURL := os.Getenv("DATABASE_URL")
if dbURL == "" {
    log.Fatal("DATABASE_URL is required")
}
```

---

## 18) Errores comunes (y como evitarlos)

1. No usar PORT de entorno.
   Solucion: leer os.Getenv("PORT").
2. Imagen demasiado grande.
   Solucion: multi-stage y runtime minima.
3. No definir health checks.
   Solucion: endpoints /health y /ready.
4. Logs sin contexto.
   Solucion: incluir campos utiles.
5. No manejar SIGTERM.
   Solucion: graceful shutdown.
6. Secretos hardcodeados.
   Solucion: env vars y secret manager.

---

## 19) Ejemplo completo de servidor cloud-friendly

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

type Config struct {
    Port string
}

func loadConfig() Config {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    return Config{Port: port}
}

func main() {
    cfg := loadConfig()

    mux := http.NewServeMux()
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ok"))
    })
    mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ready"))
    })
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        _, _ = w.Write([]byte("Hola cloud"))
    })

    srv := &http.Server{
        Addr:              ":" + cfg.Port,
        Handler:           mux,
        ReadHeaderTimeout: 5 * time.Second,
        ReadTimeout:       10 * time.Second,
        WriteTimeout:      15 * time.Second,
        IdleTimeout:       60 * time.Second,
    }

    go func() {
        log.Printf("server listening on :%s", cfg.Port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen error: %v", err)
        }
    }()

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
    <-stop

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("graceful shutdown error: %v", err)
    }

    log.Println("server stopped")
}
```

---

## 20) Checklist final para deploy real

Antes de desplegar revisa:

1. go test ./... pasa.
2. go vet ./... sin problemas.
3. Docker multi-stage habilitado.
4. PORT via entorno.
5. Endpoints /health y /ready.
6. Timeouts en servidor HTTP.
7. Graceful shutdown implementado.
8. Secretos fuera del codigo.
9. Logs claros sin datos sensibles.
10. Versionado de imagen por tag (no solo latest).

---

## 21) Ruta de aprendizaje recomendada (4 semanas)

Semana 1:

1. Sintaxis base Go.
2. net/http basico.
3. Variables de entorno y config.

Semana 2:

1. Docker multi-stage.
2. Debug local con contenedores.
3. Manejo de archivos y errores.

Semana 3:

1. Logging y metricas.
2. Graceful shutdown.
3. Seguridad basica de runtime.

Semana 4:

1. Deploy en Cloud Run o ECS.
2. Kubernetes basico (Deployment/Service/Probes).
3. CI/CD con test, build y push.

---

## 22) Cierre

Si entiendes este documento, ya tienes una base muy solida para pasar de Go local a Go en cloud real.

Tu siguiente salto natural es:

1. Desplegar tu app actual en una plataforma cloud.
2. Reemplazar carpeta local de imagenes por object storage.
3. Agregar metricas y tests automatizados.

Con eso ya estarias trabajando como backend cloud junior con buenas practicas reales.

---

## 23) Tema clave omitido: .dockerignore (importante para velocidad y seguridad)

Si no usas `.dockerignore`, Docker envia al build todo el proyecto (incluyendo basura, binarios, cache, git, etc.).

Ejemplo recomendado:

```text
.git
.github
*.log
*.tmp
tmp/
dist/
bin/
node_modules/
coverage/
server-local
server-linux
```

Beneficios:

1. Build mas rapido.
2. Menos contexto enviado al daemon.
3. Menor riesgo de subir archivos sensibles por error.

---

## 24) Tema clave omitido: 12-Factor App aplicado a Go

Para nube, una app robusta sigue principios 12-factor. Traduccion practica:

1. Config en entorno, no hardcodeada.
2. Logs a stdout/stderr, no a archivos locales.
3. Dependencias declaradas en `go.mod`.
4. Procesos stateless (sin depender de disco local para estado critico).
5. Build, release y run separados.

Como se ve en Go:

```go
port := os.Getenv("PORT")
if port == "" {
    port = "8080"
}
log.Printf("level=info msg=starting port=%s", port)
```

---

## 25) Redes y puertos en contenedores (sin confusiones)

Regla fundamental:

1. El contenedor escucha un puerto interno (por ejemplo 8080).
2. El host publica otro puerto (por ejemplo 8000:8080).

Ejemplo:

```bash
docker run --rm -p 8000:8080 mi-app
```

En este caso:

1. Tu app debe escuchar `:8080` dentro del contenedor.
2. Tu navegador accede a `localhost:8000`.

Error comun: app escuchando `:8000` y run usando `-p 8000:8080`.

---

## 26) Cliente HTTP saliente con timeout y reintentos

Muchas apps cloud no solo reciben requests, tambien llaman APIs externas. Si no pones timeout, una dependencia lenta puede congelar tus workers.

Ejemplo minimo seguro:

```go
client := &http.Client{Timeout: 5 * time.Second}

resp, err := client.Get("https://api.example.com/status")
if err != nil {
    // log y fallback
}
defer resp.Body.Close()
```

Patron de retry simple con backoff:

```go
func retry[T any](attempts int, fn func() (T, error)) (T, error) {
    var zero T
    var lastErr error
    for i := 0; i < attempts; i++ {
        out, err := fn()
        if err == nil {
            return out, nil
        }
        lastErr = err
        time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
    }
    return zero, lastErr
}
```

---

## 27) Context propagation: base para cancelacion y trazas

En Go cloud, `context.Context` viaja por toda la cadena de llamadas para:

1. Cancelar operaciones si el cliente se desconecta.
2. Aplicar deadline.
3. Propagar metadata (trace id).

Ejemplo:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    data, err := servicio(ctx)
    if err != nil {
        http.Error(w, "error", http.StatusInternalServerError)
        return
    }
    _, _ = w.Write(data)
}
```

Regla de oro: no crear context.Background() dentro de handlers para trabajo request-scoped; usa `r.Context()`.

---

## 28) Base de datos en nube: pool, timeout y migraciones

Si conectas DB en cloud debes cubrir tres puntos:

1. Pool de conexiones.
2. Timeouts de consulta.
3. Migraciones versionadas.

Ejemplo pool con `database/sql`:

```go
db, err := sql.Open("postgres", dsn)
if err != nil {
    return err
}
db.SetMaxOpenConns(20)
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(30 * time.Minute)
```

Consulta con timeout:

```go
ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
defer cancel()
row := db.QueryRowContext(ctx, "SELECT NOW()")
```

Migraciones:

1. Usa herramienta (goose, golang-migrate, atlas).
2. Nunca cambies schema manual en produccion.
3. Ejecuta migraciones en pipeline o job dedicado.

---

## 29) Versionado de API y compatibilidad

Buenas practicas:

1. Versiona endpoints: `/v1/...`, `/v2/...`.
2. Cambios breaking solo en nueva version.
3. Mantener compatibilidad mientras migran clientes.

Ejemplo:

```go
http.HandleFunc("/v1/health", healthV1)
http.HandleFunc("/v2/health", healthV2)
```

---

## 30) Autenticacion y autorizacion (intro minima)

No toda API debe ser publica. Es comun usar JWT o API keys.

Middleware de API key simple:

```go
func requireAPIKey(next http.Handler, expected string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("X-API-Key") != expected {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

En produccion:

1. API key desde secret manager.
2. Rotacion periodica de credenciales.
3. Auditoria de accesos.

---

## 31) Costos y escalado: decisiones practicas

Escalar no es solo "mas replicas".

Debes definir:

1. CPU y RAM por instancia.
2. Min/max instancias.
3. Concurrency por instancia.
4. Timeout maximo de request.

Regla simple para empezar:

1. Mide latencia p95.
2. Si p95 sube con carga, revisa bottleneck (DB, red, CPU).
3. Ajusta concurrencia antes de duplicar replicas sin analisis.

---

## 32) Rollback y estrategia de despliegue

Para evitar caidas en produccion:

1. Usa tags inmutables (`v1.0.3`) ademas de `latest`.
2. Conserva historial de revisiones desplegadas.
3. Ten comando de rollback rapido.

Estrategias comunes:

1. Rolling update.
2. Blue/Green.
3. Canary.

Principio: primero despliega a entorno de staging, valida, luego produccion.

---

## 33) Infrastructure as Code (IaC) desde cero

IaC significa definir infraestructura en codigo versionable.

Ventajas:

1. Repetible.
2. Auditable.
3. Menos error manual.

Herramientas comunes:

1. Terraform.
2. Pulumi.
3. CloudFormation (AWS).

Casos que conviene llevar a IaC:

1. Repositorio de imagenes.
2. Servicio cloud runtime.
3. Redes, IAM, secretos, base de datos.

---

## 34) Runbook operativo minimo (que hacer cuando algo falla)

Un runbook es una guia de respuesta a incidentes.

Plantilla minima:

1. Sintoma: aumento de errores 5xx.
2. Chequeo rapido: logs, metricas, despliegue reciente.
3. Accion inmediata: rollback a version previa.
4. Diagnostico: endpoint afectado, dependencia externa, DB.
5. Mitigacion: feature flag o degradacion controlada.
6. Postmortem: causa raiz y accion preventiva.

Sin runbook, cada incidente se resuelve improvisando.

---

## 35) Mapa final de madurez (de estudiante a perfil cloud)

Nivel 1:

1. App Go corre local.
2. Docker build/run funcional.

Nivel 2:

1. Deploy en cloud.
2. Config por entorno.
3. Health checks.

Nivel 3:

1. Logging estructurado.
2. Metricas y alertas.
3. Graceful shutdown.

Nivel 4:

1. Seguridad y secretos bien manejados.
2. CI/CD con pruebas y rollback.
3. IaC y runbook operativo.

Si completas nivel 4, ya trabajas con estandar de equipo profesional backend cloud.
