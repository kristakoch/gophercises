package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

func main() {

	// Initial page.
	http.Handle("/", &templateHandler{"upload.gohtml"})

	// Handles uploads.
	http.HandleFunc("/upload", uploadHandler)

	// Option pages.
	http.Handle("/mode", &optionHandler{
		"mode",
		[]string{"5", "6", "7", "8", "9"},
		"num-shapes",
	})

	http.Handle("/num-shapes", &optionHandler{
		"num-shapes",
		[]string{"50", "150", "250", "350"},
		"result",
	})

	// Final page.
	http.HandleFunc("/result", resultHandler)

	// File server.
	http.Handle("/images/", http.FileServer(http.Dir(".")))

	fmt.Println("server is running at http://localhost:9090")
	log.Fatal(http.ListenAndServe(":9090", nil))

}

type optionTemplateData struct {
	ChoiceName       string
	OriginalImageURL string
	ModifiedImages   []modifiedImage
	NextURL          string
}

type resultData struct {
	FinalImageURL string
	FullResultURL string
}

type modifiedImage struct {
	ImageSrc     string
	ChoicesQuery string
}

type templateHandler struct {
	FileName string
}

type optionHandler struct {
	ChoiceName string
	Options    []string
	NextURL    string
}

func (th *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(th.FileName))

	if err := tmpl.Execute(w, nil); err != nil {
		log.Fatalf("error during template execution, err: %s", err)
	}
}

// Handle pages that offer image options.
func (oh *optionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	// Get the query vals from the URL.
	fileName, fileExt := r.FormValue("file"), r.FormValue("ext")
	fileMode, numShapes := r.FormValue("mode"), r.FormValue("num-shapes")

	log.Printf("ensuring directory exists at %s \n", "./images/"+fileName+"/"+oh.ChoiceName+"/")

	if err := os.MkdirAll("./images/"+fileName+"/"+oh.ChoiceName+"/", 0777); err != nil {
		log.Fatal(err)
	}

	// Generate default values for each slice of options.
	modeVals, numShapesVals := generateOptionValues(
		oh.ChoiceName,
		oh.Options,
		fileMode,
		numShapes,
	)

	// Generate the modified images.
	var images []modifiedImage
	if images, err = generateModifiedImages(
		fileName,
		fileExt,
		modeVals,
		numShapesVals,
		oh.ChoiceName,
		false,
	); err != nil {
		log.Fatal(err)
	}

	originalImageURL := fmt.Sprintf("./images/%s/%s_original.%s", fileName, fileName, fileExt)

	// Build and execute the template.
	tmplData := optionTemplateData{
		ChoiceName:       oh.ChoiceName,
		OriginalImageURL: originalImageURL,
		ModifiedImages:   images,
		NextURL:          "/" + oh.NextURL,
	}

	tmpl := template.Must(template.ParseFiles("./options.gohtml"))

	if err := tmpl.Execute(w, tmplData); err != nil {
		log.Fatalf("error during template execution, err: %s", err)
	}

	http.Redirect(w, r, r.Referer(), http.StatusPermanentRedirect)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fb, fh, err := r.FormFile("fileUpload")
	if err != nil {
		log.Fatalf("error getting original image %s from form, err: %s", fh.Filename, err)
	}

	parts := strings.Split(fh.Filename, ".")
	fileName, fileExt := parts[0], parts[1]
	folderPath := fmt.Sprintf("./images/%s/", fileName)

	if err := storeOriginalImage(fileName, fileExt, folderPath, fb); err != nil {
		log.Fatalf("error storing original image file %s, err: %s", fh.Filename, err)
	}

	// After the URL values have been set, we can safely
	// send the user to the page with mode choices.
	v := url.Values{}
	v.Add("file", fileName)
	v.Add("ext", fileExt)
	v.Add("mode", "0")
	v.Add("num-shapes", "50")
	nextURL := "/mode?" + v.Encode()

	http.Redirect(w, r, nextURL, http.StatusPermanentRedirect)
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	// Regenerate the image if requested.
	var force = false
	refresh := r.FormValue("regenerateImage")
	if refresh == "Regenerate Image" {
		force = true
	}

	// Get the query vals from the URL.
	fileName, fileExt := r.FormValue("file"), r.FormValue("ext")
	fileMode, numShapes := r.FormValue("mode"), r.FormValue("num-shapes")

	if err := os.MkdirAll("./images/"+fileName+"/result/", 0777); err != nil {
		log.Fatal(err)
	}

	// Generate the final image.
	var images []modifiedImage
	if images, err = generateModifiedImages(
		fileName,
		fileExt,
		[]string{fileMode},
		[]string{numShapes},
		"result",
		force,
	); err != nil {
		log.Fatal(err)
	}

	// Build and execute the template.
	tmplData := resultData{
		FinalImageURL: images[0].ImageSrc,
		FullResultURL: r.URL.String(),
	}

	tmpl := template.Must(template.ParseFiles("./result.gohtml"))

	if err := tmpl.Execute(w, tmplData); err != nil {
		log.Fatalf("error during template execution, err: %s", err)
	}
}

func generateOptionValues(
	choiceName string,
	options []string,
	fileMode string,
	numShapes string,
) ([]string, []string) {

	fileModeVals := make([]string, len(options))
	numShapesVals := make([]string, len(options))

	// Create the slices of default values.
	for i := 0; i < len(options); i++ {
		fileModeVals[i] = fileMode
		numShapesVals[i] = numShapes
	}

	// Set the slice with unique values based on the
	// name of the choice being made.
	switch choiceName {
	case "mode":
		fileModeVals = options
	case "num-shapes":
		numShapesVals = options
	}

	return fileModeVals, numShapesVals
}

// storeOriginalImage creates the necessary folders and file
// to store the uploaded image.
func storeOriginalImage(fileName, fileExt, folderPath string, fb multipart.File) error {
	var err error

	log.Printf("ensuring directory and file exist for uploaded image %s", fileName)

	// Create the folder structure and file.
	if !(fileExt == "png" || fileExt == "jpg" || fileExt == "jpeg") {
		return fmt.Errorf("uploaded file %s.%s must be of type png or jpg/jpeg", fileName, fileExt)
	}

	if err = os.MkdirAll(folderPath, 0777); err != nil {
		return err
	}

	originalImagePath := fmt.Sprintf("./images/%s/%s_original.%s",
		fileName, fileName, fileExt)

	var newImage *os.File
	if newImage, err = os.Create(originalImagePath); err != nil {
		return err
	}

	// Store the file contents into the created file.
	if _, err = io.Copy(newImage, fb); err != nil {
		return err
	}

	return nil
}

// generateModifiedImages executes the 'primitive' command in a few
// different ways to generate artistically modified images.
func generateModifiedImages(
	fileName string,
	fileExt string,
	modes []string,
	numShapes []string,
	choiceName string,
	force bool,
) ([]modifiedImage, error) {
	var images []modifiedImage

	for i := 0; i < len(modes); i++ {
		moddedFileName := fmt.Sprintf("%s_%s-%s.%s", fileName, modes[i], numShapes[i], fileExt)
		moddedFilePath := fmt.Sprintf("./images/%s/%s/%s", fileName, choiceName, moddedFileName)
		originalFilePath := fmt.Sprintf("./images/%s/%s_original.%s", fileName, fileName, fileExt)

		// Create the URL for the next page, with the choice built in.
		v := url.Values{}
		v.Add("file", fileName)
		v.Add("ext", fileExt)
		v.Add("mode", modes[i])
		v.Add("num-shapes", numShapes[i])
		nextURL := v.Encode()

		// Don't do the work of generating the file
		// if it already exists.
		if fileExists(moddedFilePath) && !force {
			log.Printf("file %s already exists, not regenerating", moddedFileName)
			images = append(images, modifiedImage{moddedFilePath, nextURL})

			continue
		}

		log.Printf(`
		generating image %s 
		with path %s\n"
		`, moddedFileName, moddedFilePath)

		cmd := exec.Command(
			"primitive",
			"-i", originalFilePath,
			"-o", moddedFilePath,
			"-n", numShapes[i],
			"-m", modes[i],
		)

		if err := cmd.Start(); err != nil {
			return images, fmt.Errorf("error running command, err: %s", err)
		}

		if err := cmd.Wait(); err != nil {
			return images, fmt.Errorf("error waiting for command to finish, err: %s", err)
		}

		images = append(images, modifiedImage{
			moddedFilePath,
			nextURL,
		})
	}

	return images, nil
}

// fileExists checks if a file exists and is not a directory
// before it's used to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
