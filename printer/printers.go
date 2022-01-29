package printers

import (
	"os"
	"log"
	"fmt"
	"time"
	"strings"
	"net/http"
	"path/filepath"
	"encoding/json"
	"github.com/jean-bonilha/goprint"
)

type ResponseJson struct {
	Status  bool
	Message string
	Printers []string
}

var PrinterName, _ = goprint.GetDefaultPrinterName()

var fileName string = ""

var printerHandle uintptr

func GetPrinters(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/printers" {
		http.Error(rw, "404 not found.", http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	var responjson ResponseJson

	printerNames, err := goprint.GetPrinterNames()

	if err != nil {
		rw.WriteHeader(500)
		responjson = ResponseJson{
			false,
			err.Error(),
			nil,
		}
	} else {
		responjson = ResponseJson{
			true,
			"The first printer on your list is the printer default.",
			printerNames,
		}
	}

	json.NewEncoder(rw).Encode(responjson)
}

func PrintRawText(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/print" || r.Method != "POST" {
		http.Error(rw, "404 not found.", http.StatusNotFound)
		return
	}
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	var responjson ResponseJson

	if _, err := os.Stat("C:/PrinterLabel"); os.IsNotExist(err) {
		responjson = ResponseJson{
			false,
			"Error: File C:/PrinterLabel not exists.",
			nil,
		}
	}	else if err := r.ParseForm(); err != nil {
		responjson = ResponseJson{
			false,
			"Error: ParseForm().",
			nil,
		}
	}

	zpl := r.FormValue("zpl")
	printerName := r.FormValue("printer_name")
	fileName = r.FormValue("file_name")

	if printerName != "" {
		PrinterName = printerName
	}

	currentTime := time.Now()
	dateTime := string(currentTime.Format("2006-01-02.15-04-05.000000"))
	pathFilePartial := "C:/PrinterLabel/Backup/"
	pathFile := pathFilePartial + dateTime + ".txt"
	if fileName != "" {
		pathFile = pathFilePartial + dateTime + "-" + fileName + ".txt"
	}
	f, err := os.Create(pathFile)
	defer f.Close()

	printerHandle, err = goprint.GoOpenPrinter(PrinterName)
	if err != nil {
		responjson = ResponseJson{
			false,
			"Error: Create printerHandle.",
			nil,
		}
	} else {
		defer goprint.GoClosePrinter(printerHandle)

		_, err = f.WriteString(zpl)
		if err != nil {
			responjson = ResponseJson{
				false,
				"Error: Write File.",
				nil,
			}
		} else if err = PrintFile(pathFile, printerHandle, false); err != nil{
			responjson = ResponseJson{
				false,
				"Error: Print File.",
				nil,
			}
		} else {
			err = PrintFile(pathFile, printerHandle, false)
			if err == nil {
				responjson = ResponseJson{
					true,
					"Printing from " + PrinterName,
					nil,
				}
			}
		}
	}

	json.NewEncoder(rw).Encode(responjson)
}

func PrintFile(filePath string, printerHandle uintptr, backup bool) error {

	if fileInfo, err := os.Stat(filePath); err == nil && fileInfo.Size() > 0 {
		// Send to printer:
		err = goprint.GoPrint(printerHandle, filePath)
		if err != nil {
			return err
		}
	}

	if backup == false {
		return nil
	}
	// Set backup path of file
	dir, file := filepath.Split(filePath)
	backupFile := fmt.Sprintf("%sBackup\\%s", dir, file)

	// Rename file
	filePath = strings.Replace(filePath, "\\", "/", -1)
	err := os.Rename(filePath, backupFile)

	return err
}

func RunPrinterService() {
	http.HandleFunc("/printers", GetPrinters)

	http.HandleFunc("/print", PrintRawText)

	err := http.ListenAndServe(":7190", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
