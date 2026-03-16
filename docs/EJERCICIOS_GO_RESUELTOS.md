# Taller Practico de Go: Ejercicios Resueltos (desde cero)

## Como usar este documento

Este taller esta pensado para aprender Go **practicando**. Cada ejercicio trae:

1. Objetivo.
2. Concepto que entrenas.
3. Codigo resuelto.
4. Explicacion simple.
5. Resultado esperado.

Tip: no copies directo. Primero intenta resolver, luego compara.

---

## Ejercicio 1: Hola Go

### Objetivo

Imprimir un mensaje en consola.

### Concepto

- `package main`
- `func main()`
- Uso de paquete `fmt`

### Solucion

```go
package main

import "fmt"

func main() {
	fmt.Println("Hola Go")
}
```

### Explicacion

- `fmt.Println` imprime texto y salto de linea.
- Todo programa ejecutable empieza en `main()`.

### Resultado esperado

```text
Hola Go
```

---

## Ejercicio 2: Variables y tipos

### Objetivo

Declarar variables con `:=` y `var`.

### Concepto

- Inferencia de tipos.
- Tipos basicos (`string`, `int`, `bool`).

### Solucion

```go
package main

import "fmt"

func main() {
	nombre := "Diego"
	edad := 22
	activo := true

	var ciudad string = "Armenia"

	fmt.Println(nombre, edad, activo, ciudad)
}
```

### Explicacion

- `:=` declara e inicializa automaticamente.
- `var` permite declarar con tipo explicito.

### Resultado esperado

```text
Diego 22 true Armenia
```

---

## Ejercicio 3: If + comparaciones

### Objetivo

Clasificar una nota.

### Concepto

- `if`, `else if`, `else`
- Operadores relacionales

### Solucion

```go
package main

import "fmt"

func main() {
	nota := 4.3

	if nota >= 4.5 {
		fmt.Println("Excelente")
	} else if nota >= 3.0 {
		fmt.Println("Aprobado")
	} else {
		fmt.Println("Reprobado")
	}
}
```

### Resultado esperado

```text
Aprobado
```

---

## Ejercicio 4: For y sumatoria

### Objetivo

Sumar los numeros del 1 al 5.

### Concepto

- Bucle `for`
- Acumulador

### Solucion

```go
package main

import "fmt"

func main() {
	suma := 0
	for i := 1; i <= 5; i++ {
		suma += i
	}
	fmt.Println("Suma:", suma)
}
```

### Resultado esperado

```text
Suma: 15
```

---

## Ejercicio 5: Funciones con retorno

### Objetivo

Crear funcion para calcular area de rectangulo.

### Concepto

- Definicion de funciones.
- Parametros y retorno.

### Solucion

```go
package main

import "fmt"

func areaRectangulo(base float64, altura float64) float64 {
	return base * altura
}

func main() {
	area := areaRectangulo(4, 2.5)
	fmt.Println("Area:", area)
}
```

### Resultado esperado

```text
Area: 10
```

---

## Ejercicio 6: Funciones con error (estilo Go)

### Objetivo

Dividir dos numeros y manejar division por cero.

### Concepto

- Retorno multiple.
- `error`.
- Patron `if err != nil`.

### Solucion

```go
package main

import (
	"errors"
	"fmt"
)

func dividir(a float64, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("no se puede dividir por cero")
	}
	return a / b, nil
}

func main() {
	resultado, err := dividir(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Resultado:", resultado)
}
```

### Resultado esperado

```text
Resultado: 5
```

---

## Ejercicio 7: Structs (como en tu proyecto)

### Objetivo

Modelar una imagen con nombre y ruta.

### Concepto

- `type ... struct`
- Instanciacion de structs

### Solucion

```go
package main

import "fmt"

type ImageData struct {
	Name string
	Path string
}

func main() {
	img := ImageData{
		Name: "foto1.jpg",
		Path: "imagenes/foto1.jpg",
	}
	fmt.Println("Nombre:", img.Name)
	fmt.Println("Ruta:", img.Path)
}
```

### Resultado esperado

```text
Nombre: foto1.jpg
Ruta: imagenes/foto1.jpg
```

---

## Ejercicio 8: Slices y append

### Objetivo

Construir lista dinamica de extensiones validas.

### Concepto

- `[]string`
- `append`
- Recorrido con `range`

### Solucion

```go
package main

import "fmt"

func main() {
	extensiones := []string{}
	extensiones = append(extensiones, ".png")
	extensiones = append(extensiones, ".jpg")
	extensiones = append(extensiones, ".jpeg")

	for i, ext := range extensiones {
		fmt.Println(i, ext)
	}
}
```

### Resultado esperado

```text
0 .png
1 .jpg
2 .jpeg
```

---

## Ejercicio 9: Filtrar archivos por extension

### Objetivo

Filtrar nombres de archivos para quedarte solo con imagenes validas.

### Concepto

- `strings.ToLower`
- `filepath.Ext`
- Condicional OR

### Solucion

```go
package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func filtrarImagenes(nombres []string) []string {
	validos := make([]string, 0)
	for _, nombre := range nombres {
		ext := strings.ToLower(filepath.Ext(nombre))
		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
			validos = append(validos, nombre)
		}
	}
	return validos
}

func main() {
	archivos := []string{"a.png", "b.txt", "c.JPG", "d.jpeg", "e.pdf"}
	fmt.Println(filtrarImagenes(archivos))
}
```

### Resultado esperado

```text
[a.png c.JPG d.jpeg]
```

Nota: conservas el nombre original, solo normalizas extension para comparar.

---

## Ejercicio 10: Aleatoriedad sin repetir

### Objetivo

Mezclar un slice y tomar un subconjunto.

### Concepto

- `rand.Seed`
- `rand.Shuffle`
- Slicing `[:n]`

### Solucion

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func seleccionarAleatorioSinRepetir(items []string) []string {
	if len(items) == 0 {
		return nil
	}

	copia := make([]string, len(items))
	copy(copia, items)

	rand.Shuffle(len(copia), func(i, j int) {
		copia[i], copia[j] = copia[j], copia[i]
	})

	n := rand.Intn(len(copia)) + 1
	return copia[:n]
}

func main() {
	rand.Seed(time.Now().UnixNano())

	imgs := []string{"1.png", "2.png", "3.png", "4.png"}
	fmt.Println(seleccionarAleatorioSinRepetir(imgs))
}
```

### Resultado esperado

Cada ejecucion cambia, pero:

- No repite elementos.
- Siempre devuelve entre 1 y 4 elementos.

---

## Ejercicio 11: Leer archivo y convertir a Base64

### Objetivo

Transformar bytes en string Base64.

### Concepto

- `os.ReadFile`
- `encoding/base64`
- Manejo de errores

### Solucion

```go
package main

import (
	"encoding/base64"
	"fmt"
	"os"
)

func leerYCodificar(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func main() {
	texto, err := leerYCodificar("go.mod")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Base64 (primeros 40 chars):", texto[:40])
}
```

### Resultado esperado

Muestra texto Base64 (variara segun archivo).

---

## Ejercicio 12: Detectar MIME por extension

### Objetivo

Retornar tipo MIME para imagen.

### Concepto

- `switch`
- Fallback

### Solucion

```go
package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func detectarMime(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	default:
		return "application/octet-stream"
	}
}

func main() {
	fmt.Println(detectarMime("foto.jpg"))
	fmt.Println(detectarMime("icono.png"))
	fmt.Println(detectarMime("archivo.bin"))
}
```

### Resultado esperado

```text
image/jpeg
image/png
application/octet-stream
```

---

## Ejercicio 13: Primer servidor HTTP en Go

### Objetivo

Crear endpoint `/hola`.

### Concepto

- `net/http`
- `HandleFunc`
- `ListenAndServe`

### Solucion

```go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/hola", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("Hola mundo desde Go"))
	})

	fmt.Println("Servidor en http://localhost:8000")
	_ = http.ListenAndServe(":8000", nil)
}
```

### Prueba

En navegador o curl:

```bash
curl http://localhost:8000/hola
```

---

## Ejercicio 14: Servir archivos estaticos

### Objetivo

Exponer carpeta `static/` por HTTP.

### Concepto

- `http.FileServer`
- `http.StripPrefix`

### Solucion

```go
package main

import (
	"net/http"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	_ = http.ListenAndServe(":8000", nil)
}
```

### Explicacion

- URL `http://localhost:8000/static/css/styles.css`
- Apunta al archivo local `static/css/styles.css`

---

## Ejercicio 15: Renderizar plantilla HTML con datos

### Objetivo

Enviar datos dinamicos a una vista.

### Concepto

- `html/template`
- `ParseFiles`
- `Execute`

### Solucion (Go)

```go
package main

import (
	"html/template"
	"net/http"
)

type PageData struct {
	Title string
}

func main() {
	tpl := template.Must(template.ParseFiles("templates/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{Title: "Mi pagina en Go"}
		_ = tpl.Execute(w, data)
	})

	_ = http.ListenAndServe(":8000", nil)
}
```

### Solucion (HTML minimo)

```html
<h1>{{.Title}}</h1>
```

---

## Ejercicio 16: Mini version de tu flujo completo

### Objetivo

Integrar lo mas importante en una sola app.

### Concepto

- Rutas
- Lectura de imagenes
- Seleccion aleatoria
- Render HTML

### Solucion resumida

```go
package main

import (
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type PageData struct {
	Images []string
}

func loadImageNames(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	out := make([]string, 0)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
			out = append(out, e.Name())
		}
	}
	return out
}

func main() {
	rand.Seed(time.Now().UnixNano())
	tpl := template.Must(template.ParseFiles("templates/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		imgs := loadImageNames("imagenes")
		rand.Shuffle(len(imgs), func(i, j int) { imgs[i], imgs[j] = imgs[j], imgs[i] })
		if len(imgs) > 3 {
			imgs = imgs[:3]
		}
		_ = tpl.Execute(w, PageData{Images: imgs})
	})

	_ = http.ListenAndServe(":8000", nil)
}
```

### Que aprendiste aqui

1. Unir funciones pequenas para resolver problema real.
2. Mantener la logica clara por pasos.
3. Aplicar sintaxis Go en un flujo web completo.

---

## Ejercicio 17: Challenge guiado (para que lo hagas tu)

### Reto

Modificar tu servidor para que acepte query param `max` y limite el numero de imagenes.

Ejemplo:

- `/` muestra aleatorio normal.
- `/?max=2` muestra maximo 2.

### Pistas de solucion

1. Leer query:

```go
maxStr := r.URL.Query().Get("max")
```

2. Convertir string a int con `strconv.Atoi`.
3. Validar que `max > 0`.
4. Si `max < len(selected)`, recortar slice.

### Solucion propuesta

```go
maxStr := r.URL.Query().Get("max")
if maxStr != "" {
	max, convErr := strconv.Atoi(maxStr)
	if convErr == nil && max > 0 && max < len(encoded) {
		encoded = encoded[:max]
	}
}
```

---

## Ejercicio 18: Test unitario simple

### Objetivo

Probar `detectMimeType`.

### Concepto

- Paquete `testing`
- Test table-driven

### Solucion

```go
package main

import "testing"

func TestDetectMimeType(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"foto.jpg", "image/jpeg"},
		{"foto.jpeg", "image/jpeg"},
		{"logo.png", "image/png"},
		{"archivo.bin", "application/octet-stream"},
	}

	for _, tc := range cases {
		got := detectMimeType(tc.input)
		if got != tc.want {
			t.Fatalf("detectMimeType(%q) = %q; want %q", tc.input, got, tc.want)
		}
	}
}
```

### Ejecutar tests

```bash
go test ./...
```

---

## Mapa de conceptos que ya deberias dominar

Si llegaste hasta aqui, ya manejas:

1. Estructura de un programa Go.
2. Variables, funciones, if, for.
3. Error handling idiomatico de Go.
4. Structs y slices.
5. Lectura de archivos y rutas.
6. Servidor HTTP basico y rutas.
7. Plantillas HTML en Go.
8. Base64 y MIME.
9. Bases de pruebas unitarias.

---

## Plan de practica de 7 dias

Dia 1: ejercicios 1-4.
Dia 2: ejercicios 5-8.
Dia 3: ejercicios 9-12.
Dia 4: ejercicios 13-15.
Dia 5: ejercicio 16 completo.
Dia 6: challenge 17 sin mirar solucion.
Dia 7: escribir tests del ejercicio 18 y crear uno nuevo para `selectRandomWithoutRepeats`.

---

## Cierre

Este taller esta alineado directamente con la logica de tu servidor real. Si repites cada ejercicio y lo modificas (cambia variables, rutas, validaciones), vas a entender no solo la sintaxis de Go, sino la forma de pensar en Go para construir backend.
