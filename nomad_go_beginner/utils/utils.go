package utils

import (
	"compress/flate"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mholt/archiver"
)

func ExitIfError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func ExitIfStatusCodeError(resp *http.Response) {
	if resp.StatusCode != 200 {
		log.Fatalln(fmt.Sprintf("StatusCode : %d", resp.StatusCode))
	}
}

func CompressFiles(fileNames []string, zipFileName string) error {
	zip := archiver.Zip{
		CompressionLevel:       flate.BestCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      true,
		ImplicitTopLevelFolder: false,
	}
	return zip.Archive(fileNames, zipFileName)
}

func RemoveFiles(fileNames []string) error {
	for _, fileName := range fileNames {
		err := os.Remove(fileName)
		if err != nil {
			return err
		}
	}
	return nil
}
