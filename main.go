package main

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

type ImageData struct {
	Name string
	Data template.URL
}

type PageData struct {
	Title         string
	HostName      string
	GeneratedAt   string
	TotalFiles    int
	SelectedCount int
	Images        []ImageData
	Message       string
}

func main() {
	rand.Seed(time.Now().UnixNano())

	port := flag.String("port", "8000", "Puerto en el que iniciara el servidor HTTP")
	imagesDir := flag.String("dir", "imagenes", "Carpeta con imagenes .png, .jpg o .jpeg")
	flag.Parse()

	hostName, err := os.Hostname()
	if err != nil {
		hostName = "desconocido"
	}

	tpl, tplErr := template.ParseFiles("templates/index.html")
	if tplErr != nil {
		log.Fatalf("no fue posible cargar la plantilla HTML: %v", tplErr)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/hola", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("Hola mundo"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		allImages, listErr := loadImageFiles(*imagesDir)
		if listErr != nil {
			log.Printf("error al leer carpeta de imagenes: %v", listErr)
		}

		selected := selectRandomWithoutRepeats(allImages)
		encoded := make([]ImageData, 0, len(selected))
		for _, imagePath := range selected {
			img, encErr := encodeImageToDataURI(imagePath)
			if encErr != nil {
				log.Printf("error codificando %q: %v", imagePath, encErr)
				continue
			}
			encoded = append(encoded, img)
		}

		msg := ""
		if len(allImages) == 0 {
			msg = "No se encontraron imagenes validas (.png, .jpg, .jpeg) en la carpeta configurada."
		}

		data := PageData{
			Title:         "Servidor de imágenes",
			HostName:      hostName,
			GeneratedAt:   time.Now().Format("02/01/2006 15:04:05"),
			TotalFiles:    len(allImages),
			SelectedCount: len(encoded),
			Images:        encoded,
			Message:       msg,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tpl.Execute(w, data); err != nil {
			log.Printf("error renderizando plantilla: %v", err)
			http.Error(w, "Error interno al renderizar pagina", http.StatusInternalServerError)
		}
	})

	addr := ":" + strings.TrimSpace(*port)
	log.Printf("Servidor iniciado en http://localhost%s", addr)
	log.Printf("Pagina de prueba: http://localhost%s/hola", addr)
	log.Printf("Directorio de imagenes: %s", *imagesDir)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("no fue posible iniciar el servidor: %v", err)
	}
}

func loadImageFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	valid := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
			valid = append(valid, filepath.Join(dir, entry.Name()))
		}
	}

	return valid, nil
}

func selectRandomWithoutRepeats(paths []string) []string {
	if len(paths) == 0 {
		return nil
	}

	shuffled := make([]string, len(paths))
	copy(shuffled, paths)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	count := rand.Intn(len(shuffled)) + 1
	return shuffled[:count]
}

func encodeImageToDataURI(path string) (ImageData, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return ImageData{}, err
	}

	mimeType := detectMimeType(path)
	encoded := base64.StdEncoding.EncodeToString(bytes)

	return ImageData{
		Name: filepath.Base(path),
		Data: template.URL("data:" + mimeType + ";base64," + encoded),
	}, nil
}

func detectMimeType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if m := mime.TypeByExtension(ext); m != "" {
		return m
	}

	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	default:
		return "application/octet-stream"
	}
}
