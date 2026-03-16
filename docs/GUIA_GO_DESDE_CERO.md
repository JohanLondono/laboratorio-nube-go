# Curso Practico: Go desde cero explicando tu proyecto

## 0. Objetivo de esta guia

Esta guia esta escrita para alguien que **no sabe Go** y quiere entender:

- Conceptos base del lenguaje.
- Sintaxis esencial.
- Como pensar la logica en Go.
- Como funciona tu archivo `main.go` linea por linea.
- Que hace cada `import`.
- Como funcionan las rutas HTTP (`/` y `/hola`).
- Como se leen imagenes, se convierten a Base64 y se renderizan en HTML.

La idea es que este documento sea como un mini-curso aplicado al codigo real de tu proyecto.

---

## 1. Que es Go y por que se usa aqui

Go (Golang) es un lenguaje compilado creado por Google. Tiene estas ventajas para tu caso:

1. Es simple de leer y mantener.
2. Trae librerias estandar muy buenas para web (`net/http`).
3. Compila a un ejecutable unico (facil para Docker y Linux).
4. Tiene manejo de errores explicito, ideal para codigo robusto.

Tu proyecto usa Go para montar un servidor HTTP que:

- Recibe peticiones web.
- Lee archivos de imagen del disco.
- Selecciona imagenes aleatorias.
- Las codifica a Base64.
- Renderiza una plantilla HTML con esos datos.

---

## 2. Anatomia general de un programa Go

Todo ejecutable en Go tiene:

1. `package main`
2. `import (...)`
3. Tipos/funciones auxiliares
4. `func main()` como punto de entrada

Ejemplo minimo:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hola Go")
}
```

Cuando ejecutas `go run .`, Go compila y ejecuta el paquete `main`.

---

## 3. Sintaxis base de Go (lo que debes dominar primero)

> **Nota Mental 🤔:** Go está diseñado para ser explícito y fácil de leer. No intenta ocultar la complejidad detrás de "magia"; prefiere que el código diga claramente lo que hace.

### 3.1 Variables y Declaración

En Go, hay dos formas principales de declarar variables: **corta** y **explícita**.

**Declaración Corta (`:=`)**
Es la más usada, el compilador infiere el tipo de dato automáticamente. Sólo puede usarse dentro de funciones.

```go
nombre := "Johan" // Go sabe que es un string
edad := 22        // Go sabe que es un int
```

**Declaración Explícita (`var`)**
Se usa cuando quieres declarar una variable sin inicializarla de inmediato (toma un valor por defecto: `""` para strings, `0` para ints), o a nivel de paquete (fuera de funciones).

```go
var puerto string = "8000"
var activo bool   // Toma por defecto 'false'
var contador int  // Toma por defecto '0'
```

### 3.2 Tipos Comunes

Go es _fuertemente tipado_. Todo debe tener un tipo y los tipos no se mezclan automáticamente.

- `string`: texto (`"Hola"`)
- `int`: enteros positivos o negativos (`-5`, `42`)
- `float64`: números con decimales (`3.1416`)
- `bool`: verdadero/falso (`true`, `false`)

**Tipos compuestos:**

- `[]T` (Slice): Una lista dinámica. Ej: `[]string` es una lista de textos.
- `struct`: Tu propio tipo de dato agrupando otros. Igual a cómo creas "modelos" o "clases" para estructurar datos.

### 3.3 Funciones (El corazón de Go)

Las funciones son bloques de código reutilizables. Comienzan con la palabra `func`.

**Función simple:**

```go
func saludar(nombre string) {
    fmt.Println("Hola", nombre)
}
```

**Función con retorno:**
Se debe especificar el tipo de dato que va a retornar al final de la definición.

```go
func sumar(a int, b int) int {
    return a + b
}
```

> **🔥 El Superpoder de Go: Múltiples Retornos**
> Es sumamente común que las funciones en Go devuelvan dos cosas: el resultado esperado y un posible error.

```go
func dividir(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("división por cero no permitida") // Retorna resultado vacío y un Error
    }
    return a / b, nil // nil significa "no hay error"
}
```

### 3.4 Condicionales (`if`) y Bucles (`for`)

**If / Else:** No necesita paréntesis alrededor de la condición, pero sí exige las llaves `{}`.

```go
if edad >= 18 {
    fmt.Println("Mayor de edad")
} else {
    fmt.Println("Menor de edad")
}
```

**El único bucle: `for`**
Go decidió eliminar `while` o `do-while`. Todo se hace con `for`.

_Ciclo clásico:_

```go
for i := 0; i < 5; i++ {
    fmt.Println(i)
}
```

_Como un "While":_

```go
x := 0
for x < 5 {
    // ...
    x++
}
```

> **🌟 Iterando listas (`range`)**
> Para recorrer slices (listas) se usa la magia del `range`. Te entrega dos variables: el `índice` (posición) y el `valor`.

```go
nombres := []string{"Diego", "Johan"}
for indice, valor := range nombres {
    fmt.Println("En la posición", indice, "está:", valor)
}
// Si no quieres usar el índice, lo ignoras con un guion bajo `_`
for _, valor := range nombres { ... }
```

### 3.5 El famoso Manejo de Errores en Go

En otros lenguajes (Java, JS) usas `try/catch`. En Go, **los errores son valores**, se retornan y se chequean manualmente de inmediato. Esto hace el código muy robusto porque te obliga a pensar qué hacer si algo falla.

```go
contenido, err := os.ReadFile("foto.png")
if err != nil { // "Si el error NO es nulo" (hubo un fallo)
    log.Println("Ups, falló al leer:", err)
    return // Sale de la función porque falló
}
// Si llega aquí, es porque err es nil y todo salió bien.
fmt.Println("Archivo leído exitosamente")
```

Este patrón `if err != nil` es lo que más vas a escribir y ver en tu proyecto.

### 3.6 Punteros (El caso de `flag`)

Un puntero almacena la **dirección en memoria** de una variable, no su valor directo.
Imagina que te doy la _dirección de una casa_ (el puntero, `&`) vs darte la _casa entera_ (el valor).

En tu proyecto usas `flag`:

```go
// flag.String no devuelve un texto, devuelve LA DIRECCIÓN de memoria donde se guarda ese texto
port := flag.String("port", "8000", "Puerto HTTP")

flag.Parse()

// Para leer el valor real guardado en esa dirección, usamos el asterisco (*)
fmt.Println("El puerto es:", *port)
```

---

## 4. Explicacion de los imports de tu proyecto

En tu `main.go` usas:

```go
import (
    "encoding/base64"
    "flag"
    "html/template"
    "log"
    "math/rand"
    "mime"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"
)
```

Que hace cada uno:

1. `encoding/base64`: convierte bytes de imagen a texto Base64.
2. `flag`: lee argumentos de consola como `-port=8000`.
3. `html/template`: renderiza HTML con datos de forma segura.
4. `log`: imprime logs (info/error) con formato.
5. `math/rand`: aleatoriedad para seleccionar imagenes.
6. `mime`: obtiene tipo MIME (`image/png`, `image/jpeg`) por extension.
7. `net/http`: servidor web, rutas y respuestas HTTP.
8. `os`: hostname, lectura de archivos/directorios.
9. `path/filepath`: manejo portable de rutas y extensiones.
10. `strings`: utilidades de texto (`ToLower`, `TrimSpace`).
11. `time`: fecha/hora y semilla aleatoria.

---

## 5. Estructuras (struct) usadas en tu codigo

## 5.1 `ImageData`

```go
type ImageData struct {
    Name string
    Data template.URL
}
```

Representa una imagen lista para la plantilla:

- `Name`: nombre del archivo.
- `Data`: Data URI (ejemplo: `data:image/png;base64,...`).

## 5.2 `PageData`

```go
type PageData struct {
    Title         string
    HostName      string
    GeneratedAt   string
    TotalFiles    int
    SelectedCount int
    Images        []ImageData
    Message       string
}
```

Es el paquete completo de datos que se envia al HTML.

---

## 6. Flujo completo del `main()` explicado como clase

## 6.1 Semilla aleatoria

```go
rand.Seed(time.Now().UnixNano())
```

Si no haces esto, la aleatoriedad puede repetir patrones al reiniciar.

## 6.2 Parametros por CLI

```go
port := flag.String("port", "8000", "Puerto...")
imagesDir := flag.String("dir", "imagenes", "Carpeta...")
flag.Parse()
```

- `-port`: puerto del servidor.
- `-dir`: carpeta donde estan las imagenes.

## 6.3 Hostname de la maquina

```go
hostName, err := os.Hostname()
if err != nil {
    hostName = "desconocido"
}
```

Esto muestra en la pagina en que host se esta ejecutando.

## 6.4 Cargar plantilla HTML

```go
tpl, tplErr := template.ParseFiles("templates/index.html")
if tplErr != nil {
    log.Fatalf("no fue posible cargar la plantilla HTML: %v", tplErr)
}
```

- Si la plantilla no existe o tiene error, el programa termina.
- `Fatalf` escribe log y finaliza.

## 6.5 Servir archivos estaticos

```go
http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
```

Significa:

1. Si llega una URL que comienza con `/static/`...
2. Se busca en carpeta local `static/`...
3. Se entrega el archivo (CSS en tu caso).

Ejemplo:

- URL: `/static/css/styles.css`
- Archivo real: `static/css/styles.css`

## 6.6 Ruta `/hola`

```go
http.HandleFunc("/hola", func(w http.ResponseWriter, _ *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    _, _ = w.Write([]byte("Hola mundo"))
})
```

- Define endpoint de prueba.
- Devuelve texto plano.
- `_ *http.Request` ignora el parametro porque no se usa.

## 6.7 Ruta principal `/`

```go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    // ... resto del flujo
})
```

La validacion evita que esta ruta capture accidentalmente otras rutas.

Dentro de esta ruta ocurre toda la logica:

1. Leer imagenes validas (`loadImageFiles`).
2. Seleccionar subset aleatorio sin repetir (`selectRandomWithoutRepeats`).
3. Codificar cada imagen (`encodeImageToDataURI`).
4. Construir `PageData`.
5. Renderizar HTML con `tpl.Execute(w, data)`.

## 6.8 Arrancar servidor

```go
addr := ":" + strings.TrimSpace(*port)
log.Printf("Servidor iniciado en http://localhost%s", addr)
if err := http.ListenAndServe(addr, nil); err != nil {
    log.Fatalf("no fue posible iniciar el servidor: %v", err)
}
```

- Si `port = 8000`, `addr` queda `:8000`.
- `ListenAndServe` bloquea el hilo principal y queda escuchando peticiones.

---

## 7. Funciones auxiliares explicadas

## 7.1 `loadImageFiles(dir string) ([]string, error)`

Responsabilidad: listar archivos validos en la carpeta.

Pasos:

1. Lee entradas con `os.ReadDir(dir)`.
2. Ignora subdirectorios.
3. Obtiene extension con `filepath.Ext`.
4. Convierte extension a minuscula (`strings.ToLower`).
5. Acepta solo `.png`, `.jpg`, `.jpeg`.
6. Devuelve slice de rutas.

Conceptos aqui:

- Filtrado por condicion.
- Uso de `append` para slices.
- Retorno `(resultado, error)`.

## 7.2 `selectRandomWithoutRepeats(paths []string) []string`

Responsabilidad: elegir una cantidad aleatoria de imagenes sin repetir.

Pasos:

1. Si no hay entradas, retorna `nil`.
2. Copia slice original para no mutarlo.
3. Lo mezcla con `rand.Shuffle`.
4. Elige `count` entre 1 y `len(shuffled)`.
5. Retorna `shuffled[:count]`.

Importante:

- No repite porque baraja una sola vez y corta.
- Siempre devuelve minimo 1 si hay imagenes.

## 7.3 `encodeImageToDataURI(path string) (ImageData, error)`

Responsabilidad: convertir archivo de imagen en algo embebible en HTML.

Pasos:

1. Lee bytes con `os.ReadFile(path)`.
2. Detecta MIME (`detectMimeType`).
3. Convierte bytes a Base64.
4. Crea string tipo `data:image/png;base64,AAAA...`.
5. Retorna `ImageData`.

## 7.4 `detectMimeType(path string) string`

Responsabilidad: deducir MIME por extension.

- Intenta con `mime.TypeByExtension(ext)`.
- Si falla, usa `switch` manual.
- Si no coincide nada, `application/octet-stream`.

---

## 8. Rutas HTTP del proyecto y su logica

Rutas activas:

1. `/hola`
2. `/`
3. `/static/...` (archivos estaticos)

## 8.1 Diferencia entre `Handle` y `HandleFunc`

- `http.Handle(pattern, handler)` recibe un objeto que implementa `ServeHTTP`.
- `http.HandleFunc(pattern, func)` permite pasar una funcion directamente.

En tu codigo:

- `Handle` para `FileServer` de estaticos.
- `HandleFunc` para endpoints con funciones anonimas.

## 8.2 Objeto Request y ResponseWriter

En handlers:

```go
func(w http.ResponseWriter, r *http.Request)
```

- `r` contiene datos entrantes (path, headers, query, metodo).
- `w` se usa para escribir status, headers y body de respuesta.

---

## 9. Plantillas Go (`html/template`) en contexto

Tu HTML usa sintaxis tipo:

- `{{.Title}}`
- `{{if .Message}} ... {{end}}`
- `{{range .Images}} ... {{end}}`

Regla mental:

- El punto `.` representa el objeto actual.
- Al iniciar `Execute(w, data)`, `.` es `PageData`.
- Dentro de `range .Images`, `.` pasa a ser cada `ImageData`.

Eso explica por que dentro del `range` usas `{{.Name}}` y `{{.Data}}`.

---

## 10. Conceptos de PENSAMIENTO Go (La Filosofía, sin notarlo)

Hasta ahora has programado en Go como lo hacen en Google gracias a estos conceptos inconscientes que tiene tu proyecto:

1. **Retorno Temprano (Early Return):** No hay indentaciones gigantes. Se valida el error, se maneja, y se usa `return`. El "camino feliz" de la función queda limpio, pegado al margen izquierdo (sin meter todo dentro de un `else`).
2. **Pequeñas Piezas Especializadas:** Tu código no hace todo en el `main`. Tienes una función que sólo lee, otra que sólo mezcla, otra que sólo codifica. Eso se llama _Cohesión Fuerte_ e _Interfaces Simples_.
3. **Manejo explícito sobre la magia:** Go te exige leer explícitamente el directorio, validar la extensión de frente e interpolar errores. Esto previene que el código falle con cosas raras imposibles de debugear.
4. **Agrupar Datos para el Contexto:** Crear ese gran paquete `PageData` en lugar de pasar 10 variables sueltas a la vista. Esa es la manera **Orientada a Datos** que usa Go: creas una estructura clara del estado que necesitas y lo procesas.

---

## 11. Buenas practicas recomendadas para seguir aprendiendo

## 11.1 Evitar estado global de `math/rand`

Tu doc original ya lo menciona. Podrias migrar a:

```go
r := rand.New(rand.NewSource(time.Now().UnixNano()))
```

Y luego usar `r.Shuffle`, `r.Intn`.

## 11.2 Validar tamano de archivo

Antes de codificar, puedes revisar tamano para evitar respuestas HTML gigantes.

## 11.3 Escribir pruebas unitarias

Funciones ideales para test:

- `loadImageFiles`
- `selectRandomWithoutRepeats`
- `detectMimeType`

## 11.4 Separar handlers en archivos

Cuando el proyecto crezca:

- `cmd/server/main.go`
- `internal/http/handlers.go`
- `internal/service/image_service.go`

---

## 12. Glosario rapido (terminos clave)

1. Compilar: convertir codigo Go a ejecutable.
2. Handler: funcion que responde una ruta HTTP.
3. Slice: lista dinamica.
4. Struct: tipo de datos personalizado con campos.
5. MIME: tipo de contenido de archivo (`image/png`, etc.).
6. Data URI: contenido embebido en una URL `data:...`.
7. Template: plantilla HTML con marcadores dinamicos.
8. Pointer: variable que guarda direccion de memoria.

---

## 13. Mapa mental de tu request `GET /`

1. Cliente abre `/`.
2. Handler valida path exacto.
3. Lee archivos de imagen en `-dir`.
4. Filtra extensiones validas.
5. Mezcla y toma N aleatorias.
6. Convierte cada imagen a Base64.
7. Arma `PageData` con metadata.
8. Ejecuta plantilla HTML.
9. Navegador muestra cards con imagenes.

---

## 14. Comandos utiles para practicar

Ejecutar con defaults:

```bash
go run .
```

Ejecutar cambiando puerto y carpeta:

```bash
go run . -port=8080 -dir=imagenes
```

Compilar binario:

```bash
go build -o servidor-imagenes .
```

Ejecutar binario:

```bash
./servidor-imagenes
```

Probar endpoint rapido:

```bash
curl http://localhost:8000/hola
```

---

## 15. Ruta de aprendizaje sugerida (tipo curso)

Semana 1:

1. Sintaxis minima: variables, `if`, `for`, funciones.
2. Manejo de errores con `err != nil`.
3. Structs y slices.

Semana 2:

1. `net/http`: handlers, metodos, status codes.
2. `html/template`: `if`, `range`, datos dinamicos.
3. Proyecto mini CRUD en memoria.

Semana 3:

1. Archivos y directorios (`os`, `filepath`).
2. JSON (`encoding/json`).
3. Tests unitarios con `testing`.

Semana 4:

1. Concurrencia basica con goroutines y channels.
2. Contextos (`context.Context`) para timeouts.
3. Estructura profesional de proyectos Go.

---

## 16. Resumen final

Tu codigo ya implementa un backend web funcional y bien organizado para nivel inicial:

- Lee entradas por CLI.
- Levanta servidor HTTP.
- Maneja rutas y estaticos.
- Procesa archivos del sistema.
- Convierte binario a Base64.
- Renderiza HTML dinamico con plantilla.

Si dominas este archivo, ya tienes una base muy solida de Go aplicado a desarrollo web.

Siguiente paso recomendado: crear pruebas unitarias y luego separar codigo en paquetes para escalar el proyecto.
