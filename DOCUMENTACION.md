# Documentacion del Proyecto - Servidor de Imagenes en Go

## 1. Objetivo del proyecto

Este proyecto implementa un servidor HTTP en Go que:

- Muestra una pagina web responsive con Bootstrap.
- Toma imagenes desde una carpeta del sistema de archivos.
- Filtra solo archivos `.png`, `.jpg` y `.jpeg`.
- Selecciona una cantidad aleatoria de imagenes sin repetir.
- Convierte cada imagen a Base64 y la incrusta en el HTML con Data URI.
- Muestra una ruta de prueba `Hola mundo` en `/hola`.

## 2. Estructura del proyecto

```text
go.mod
main.go
templates/
  index.html
static/
  css/
    styles.css
imagenes/   (carpeta esperada por defecto)
```

## 3. Como ejecutar

Desde la raiz del proyecto:

```bash
go run . -port=8000 -dir=imagenes
```

Opciones:

- `-port`: puerto HTTP del servidor (default `8000`).
- `-dir`: carpeta con imagenes (default `imagenes`).

Rutas principales:

- `http://localhost:8000/hola` -> pagina de prueba.
- `http://localhost:8000/` -> pagina principal con imagenes.

## 4. Flujo de funcionamiento

1. `main()` lee argumentos de linea de comandos con `flag`.
2. Carga el hostname de la maquina (`os.Hostname()`).
3. Carga la plantilla HTML (`templates/index.html`).
4. Expone recursos estaticos en `/static/` para CSS.
5. Ruta `/hola`: responde texto plano "Hola mundo".
6. Ruta `/`:
   - Lee archivos de la carpeta de imagenes.
   - Filtra por extension valida.
   - Baraja y selecciona N imagenes (N aleatorio, sin repetidas).
   - Convierte cada imagen a Base64.
   - Renderiza la plantilla con los datos (`PageData`).

## 5. Explicacion del codigo Go (sintaxis aplicada)

### 5.1 Paquete e imports

En Go, todo programa ejecutable inicia con:

```go
package main
```

`main` indica que el binario se ejecuta desde la funcion `main()`.

Los `import` agregan librerias:

- `net/http`: servidor web y manejo de rutas.
- `html/template`: renderizado seguro de plantillas.
- `flag`: argumentos CLI.
- `os`, `path/filepath`: acceso al sistema de archivos.
- `encoding/base64`: codificacion de bytes a Base64.
- `math/rand`, `time`: aleatoriedad y timestamps.

### 5.2 Structs

Se usan `struct` para agrupar datos:

```go
type ImageData struct {
    Name string
    Data template.URL
}
```

- `Name`: nombre del archivo.
- `Data`: Data URI en Base64.
- `template.URL`: tipo seguro para usar en `src="..."` dentro de plantilla.

`PageData` agrupa toda la informacion que se envia al HTML.

### 5.3 Declaracion de variables con punteros (flag)

```go
port := flag.String("port", "8000", "...")
```

- `flag.String` devuelve `*string` (puntero).
- Se usa `*port` para leer el valor real.

### 5.4 Manejo de errores idiomatico

Go usa retorno explicito de errores:

```go
hostName, err := os.Hostname()
if err != nil {
    hostName = "desconocido"
}
```

Patron comun:

- Llamada a funcion.
- Verificacion inmediata de `err`.
- Accion de fallback o retorno.

### 5.5 Registro de rutas HTTP

```go
http.HandleFunc("/hola", func(w http.ResponseWriter, _ *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    _, _ = w.Write([]byte("Hola mundo"))
})
```

Sintaxis clave:

- `func(...) { ... }` define una funcion anonima.
- `w http.ResponseWriter` permite escribir respuesta.
- `_ *http.Request` ignora parametro no usado.

### 5.6 Rutas estaticas

```go
http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
```

Esto permite servir archivos fisicos de `static/` bajo URL `/static/...`.

### 5.7 Ciclos y filtros

Filtrado de extensiones:

```go
for _, entry := range entries {
    if entry.IsDir() {
        continue
    }

    ext := strings.ToLower(filepath.Ext(entry.Name()))
    if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
        valid = append(valid, filepath.Join(dir, entry.Name()))
    }
}
```

Sintaxis clave:

- `for _, x := range lista` recorre colecciones.
- `continue` salta a la siguiente iteracion.
- `append(slice, item)` agrega elementos a un slice.

### 5.8 Slices y aleatoriedad

```go
shuffled := make([]string, len(paths))
copy(shuffled, paths)
rand.Shuffle(len(shuffled), func(i, j int) {
    shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
})
count := rand.Intn(len(shuffled)) + 1
return shuffled[:count]
```

Puntos importantes:

- `make([]T, n)` crea slices.
- `copy(dst, src)` copia contenido.
- `rand.Shuffle` mezcla sin perder elementos.
- `shuffled[:count]` hace slicing del rango inicial.

### 5.9 Base64 y Data URI

```go
encoded := base64.StdEncoding.EncodeToString(bytes)
Data: template.URL("data:" + mimeType + ";base64," + encoded)
```

- Convierte bytes de imagen a texto Base64.
- Construye `data:image/...;base64,...` para incrustar en HTML.

### 5.10 switch en Go

```go
switch ext {
case ".jpg", ".jpeg":
    return "image/jpeg"
case ".png":
    return "image/png"
default:
    return "application/octet-stream"
}
```

`switch` evita multiples `if` y mejora legibilidad.

## 6. Plantilla HTML y CSS

- `templates/index.html` usa sintaxis de Go templates (`{{.Campo}}`, `{{range ...}}`, `{{if ...}}`).
- `static/css/styles.css` contiene estilos personalizados para header, cards, hostname y footer.
- Bootstrap 5.3 aporta grid responsive (`col-12 col-sm-6 col-lg-4`).

## 7. Notas importantes

- Si no aparecen imagenes, revisa que el directorio indicado en `-dir` exista y contenga archivos validos.
- El proyecto evita imagenes repetidas en la misma carga por el algoritmo de mezcla + slicing.
- El numero de imagenes mostradas cambia en cada request por seleccion aleatoria.

## 8. Posibles mejoras

- Validar tamano maximo de imagen para evitar respuestas muy pesadas.
- Recorrer subdirectorios de forma recursiva.
- Agregar pruebas unitarias para `loadImageFiles`, `selectRandomWithoutRepeats` y `detectMimeType`.
- Migrar a `rand.New(rand.NewSource(...))` para desacoplar estado global de aleatoriedad.
