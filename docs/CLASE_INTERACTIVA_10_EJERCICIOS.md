# Clase interactiva: 10 ejercicios sin resolver (Go)

## Instrucciones

1. Responde cada ejercicio debajo de su seccion.
2. Puedes responder con explicacion o con codigo.
3. No mires los ejercicios resueltos mientras practicas.
4. Cuando termines, me envias el archivo y yo te corrijo como profesor.

---

## Ejercicio 1: Variables y tipos basicos

Crea un programa que declare:

- un nombre (string)
- una edad (int)
- un valor booleano de estudiante (bool)

Luego imprime todo en una sola linea.

### Tu respuesta

---

## Ejercicio 2: Funcion con retorno

Crea una funcion `sumar` que reciba dos enteros y retorne el resultado.
En `main`, llama la funcion con 7 y 5 e imprime el resultado.

### Tu respuesta

---

## Ejercicio 3: Condicionales

Crea una variable `nota` de tipo `float64`.

- Si `nota >= 4.5`, imprime: Excelente
- Si `nota >= 3.0`, imprime: Aprobado
- Si no, imprime: Reprobado

### Tu respuesta

---

## Ejercicio 4: Bucle for

Usa un `for` para sumar los numeros del 1 al 10 e imprime el total final.

### Tu respuesta

---

## Ejercicio 5: Slice y range

Crea un slice con estos valores:
`.png`, `.jpg`, `.jpeg`, `.txt`

Recorre el slice y muestra solo las extensiones validas de imagen (`.png`, `.jpg`, `.jpeg`).

### Tu respuesta

---

## Ejercicio 6: Struct basico

Define un struct llamado `Persona` con campos:

- `Nombre`
- `Edad`

Crea una instancia con tus datos y muestrala en consola.

### Tu respuesta

---

## Ejercicio 7: Manejo de errores estilo Go

Crea una funcion:
`dividir(a, b float64) (float64, error)`

Reglas:

- Si `b` es 0, retorna error.
- En `main`, prueba `dividir(10, 0)`.
- Maneja el error con `if err != nil`.

### Tu respuesta

---

## Ejercicio 8: Lectura de parametros por consola

Usa `flag` para leer:

- `-port` (por defecto `8000`)
- `-dir` (por defecto `imagenes`)

Luego imprime ambos valores ya parseados.

### Tu respuesta

---

## Ejercicio 9: Ruta HTTP simple

Crea un servidor con `net/http` que tenga:

- Ruta `/hola` que responda `Hola mundo`
- Header `Content-Type` en `text/plain; charset=utf-8`

Dejalo escuchando en puerto `8000`.

### Tu respuesta

---

## Ejercicio 10: Logica tipo proyecto (sin Base64)

Crea una funcion que reciba un directorio y retorne solo archivos con extension:

- `.png`
- `.jpg`
- `.jpeg`

Pistas:

- Lee directorio
- Ignora carpetas
- Usa `filepath.Ext` y `strings.ToLower`

En `main`, llama la funcion con `imagenes` e imprime el listado resultante.

### Tu respuesta

---

## Entrega

Cuando completes todos los ejercicios, me dices:

- Ya termine la clase interactiva

Y te hago correccion completa por ejercicio:

- Que esta bien
- Que mejorar
- Nota por ejercicio
- Version recomendada
