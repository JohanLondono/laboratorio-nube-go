# LAB Go + Nube Paso a Paso

## Objetivo del laboratorio

En este laboratorio vas a:

1. Ejecutar tu app Go localmente.
2. Contenerizarla con Docker.
3. Verificar que corre bien en contenedor.
4. Publicar la imagen en un registro.
5. Desplegarla en la nube.
6. Validar health checks y URL publica.

Este lab esta pensado para copiar/pegar comandos en orden.

---

## 0) Requisitos previos

Instala y verifica:

1. Go (1.21+ recomendado)
2. Docker Desktop
3. Git (opcional, recomendado)
4. Cuenta cloud (GCP recomendado para este lab)
5. CLI de GCP: gcloud

Verifica versiones:

```powershell
go version
docker --version
gcloud --version
```

Si alguno falla, instala antes de continuar.

---

## 1) Estructura esperada del proyecto

Debes estar en la carpeta del proyecto que contiene:

- main.go
- go.mod
- Dockerfile
- templates/
- static/
- imagenes/

En PowerShell, ubicate en el proyecto:

```powershell
Set-Location "g:\Mi unidad\Universidad_del_Quindio\9no_Semestre\Cloud_Computing\Laboratorios\laboratorio-nube-go"
Get-ChildItem
```

---

## 2) Ejecutar local en Go (sin Docker)

### 2.1 Ejecutar servidor

```powershell
go run . -port=8000 -dir=imagenes
```

### 2.2 Probar endpoints

En otra terminal:

```powershell
curl http://localhost:8000/hola
curl http://localhost:8000/
```

Si ves respuesta de hola y HTML en /, esta OK.

Deten el proceso con Ctrl+C.

---

## 3) Build local del binario (opcional, recomendado)

```powershell
go build -o server-local .
./server-local -port=8000 -dir=imagenes
```

Esto valida que compila antes de ir a Docker.

Deten con Ctrl+C.

---

## 4) Construir imagen Docker

### 4.1 Build

```powershell
docker build -t laboratorio-go:lab1 .
```

### 4.2 Ver imagen

```powershell
docker images | findstr laboratorio-go
```

---

## 5) Ejecutar contenedor local

### 5.1 Run simple

```powershell
docker run --rm -p 8000:8000 laboratorio-go:lab1
```

### 5.2 Probar en otra terminal

```powershell
curl http://localhost:8000/hola
curl http://localhost:8000/
```

Si todo responde, tu imagen funciona.

Deten con Ctrl+C.

### 5.3 Run montando carpeta imagenes local (opcional)

```powershell
docker run --rm -p 8000:8000 -v "${PWD}\imagenes:/app/imagenes" laboratorio-go:lab1
```

Nota en Windows: si Docker no reconoce la ruta, prueba con ruta absoluta y slash normal:

```powershell
docker run --rm -p 8000:8000 -v "g:/Mi unidad/Universidad_del_Quindio/9no_Semestre/Cloud_Computing/Laboratorios/laboratorio-nube-go/imagenes:/app/imagenes" laboratorio-go:lab1
```

---

## 6) Preparar despliegue en GCP Cloud Run

Este es el camino mas facil para tu caso.

### 6.1 Login en GCP

```powershell
gcloud auth login
```

### 6.2 Definir proyecto

Reemplaza TU_PROJECT_ID:

```powershell
gcloud config set project TU_PROJECT_ID
```

### 6.3 Habilitar APIs necesarias

```powershell
gcloud services enable run.googleapis.com

gcloud services enable artifactregistry.googleapis.com

gcloud services enable cloudbuild.googleapis.com
```

### 6.4 Crear repositorio de Artifact Registry

```powershell
gcloud artifacts repositories create go-lab-repo --repository-format=docker --location=us-central1 --description="Repo Docker para lab Go"
```

Si ya existe, puedes seguir sin problema.

---

## 7) Etiquetar y subir imagen al registro

### 7.1 Definir variables en PowerShell

Reemplaza TU_PROJECT_ID:

```powershell
$PROJECT_ID="TU_PROJECT_ID"
$REGION="us-central1"
$REPO="go-lab-repo"
$IMAGE="laboratorio-go"
$TAG="v1"
$FULL_IMAGE="$REGION-docker.pkg.dev/$PROJECT_ID/$REPO/$IMAGE:$TAG"
```

### 7.2 Configurar autenticacion Docker con gcloud

```powershell
gcloud auth configure-docker "$REGION-docker.pkg.dev"
```

### 7.3 Tag + Push

```powershell
docker tag laboratorio-go:lab1 $FULL_IMAGE
docker push $FULL_IMAGE
```

---

## 8) Desplegar en Cloud Run

Importante: Cloud Run inyecta variable PORT. Tu app debe usar PORT del entorno para ser 100% cloud-native.

Si tu app aun usa solo flag -port fijo, igual puede correr en algunos escenarios, pero lo correcto es soportar `os.Getenv("PORT")` como prioridad.

### 8.1 Deploy

```powershell
gcloud run deploy go-lab-service --image $FULL_IMAGE --platform managed --region $REGION --allow-unauthenticated
```

Cuando termine, mostrara una URL publica.

---

## 9) Validacion final en nube

### 9.1 Probar endpoint de prueba

Reemplaza TU_URL:

```powershell
curl TU_URL/hola
```

### 9.2 Probar pagina principal

```powershell
curl TU_URL/
```

### 9.3 Ver logs en tiempo real

```powershell
gcloud run services logs read go-lab-service --region $REGION
```

---

## 10) Actualizar version (iteracion de despliegue)

Cada cambio en codigo sigue este ciclo:

1. Build nueva imagen
2. Tag nuevo (v2, v3...)
3. Push
4. Deploy

Ejemplo v2:

```powershell
$TAG="v2"
$FULL_IMAGE="$REGION-docker.pkg.dev/$PROJECT_ID/$REPO/$IMAGE:$TAG"

docker build -t laboratorio-go:lab2 .
docker tag laboratorio-go:lab2 $FULL_IMAGE
docker push $FULL_IMAGE

gcloud run deploy go-lab-service --image $FULL_IMAGE --platform managed --region $REGION --allow-unauthenticated
```

---

## 11) Limpieza (opcional)

Eliminar servicio Cloud Run:

```powershell
gcloud run services delete go-lab-service --region $REGION
```

Eliminar imagen local:

```powershell
docker rmi laboratorio-go:lab1
```

---

## 12) Checklist de aprobacion del laboratorio

Marca cada punto:

- [ ] Go local responde /hola
- [ ] Docker build exitoso
- [ ] Contenedor local responde /hola y /
- [ ] Imagen subida a Artifact Registry
- [ ] Servicio desplegado en Cloud Run
- [ ] URL publica funcional
- [ ] Logs visibles desde gcloud

Si completas todo, ya hiciste un flujo real Dev -> Container -> Cloud.

---

## 13) Troubleshooting rapido

### Error: port no disponible localmente

Causa: ya hay algo en 8000.
Solucion:

```powershell
# usa otro puerto local
docker run --rm -p 8080:8000 laboratorio-go:lab1
curl http://localhost:8080/hola
```

### Error: Permission denied o ruta con espacios al montar volumen

Causa: path Windows con espacios o permisos de Docker Desktop.
Solucion:

1. Habilita sharing de disco en Docker Desktop.
2. Usa ruta absoluta con slash normal (g:/...).
3. Prueba sin volumen para aislar el problema.

### Error: Cloud Run no inicia el contenedor

Causa comun: app no escucha en puerto esperado por PORT.
Solucion recomendada en Go:

```go
port := os.Getenv("PORT")
if port == "" {
    port = "8000"
}
addr := ":" + port
```

### Error: push denied al registry

Causa: auth o proyecto incorrecto.
Solucion:

1. Revisar `gcloud config get-value project`
2. Repetir `gcloud auth configure-docker`
3. Verificar nombre completo de imagen

---

## 14) Extension opcional: AWS ECS (ruta corta)

Si quieres AWS en lugar de GCP, el flujo conceptual es igual:

1. Crear repositorio ECR.
2. Hacer login Docker contra ECR.
3. Tag + push de imagen.
4. Crear task definition ECS.
5. Crear servicio ECS en Fargate.
6. Exponer con ALB.

Comandos guia (resumen):

```bash
aws ecr create-repository --repository-name laboratorio-go
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account>.dkr.ecr.us-east-1.amazonaws.com
docker tag laboratorio-go:lab1 <account>.dkr.ecr.us-east-1.amazonaws.com/laboratorio-go:v1
docker push <account>.dkr.ecr.us-east-1.amazonaws.com/laboratorio-go:v1
```

Luego despliegas en ECS con consola o IaC.

---

## 15) Mejora recomendada para tu codigo antes de nube (muy importante)

En apps cloud modernas conviene que el puerto venga del entorno y, si no existe, use fallback.

Patron sugerido:

```go
portValue := os.Getenv("PORT")
if strings.TrimSpace(portValue) == "" {
    portValue = *port
}
addr := ":" + strings.TrimSpace(portValue)
```

Con esto tu app funciona igual en local y en plataformas cloud.

---

## 16) Resultado esperado del laboratorio

Si hiciste todo bien, al final tendras:

1. Una app Go funcional en contenedor.
2. Una imagen versionada en registry.
3. Un servicio publico en nube accesible por URL.
4. Flujo base para hacer CI/CD despues.

Ya estaras ejecutando un pipeline real de backend cloud.
