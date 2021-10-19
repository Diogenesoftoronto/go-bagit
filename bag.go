package go_bagit

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func ValidateBag(bagLocation string, fast bool, complete bool) error {
	errors := []error{}
	storedOxum, err := GetOxum(bagLocation)
	if err != nil {
		log.Printf("- ERROR - %s", err.Error())
		return err
	}

	err = ValidateOxum(bagLocation, storedOxum)
	if err != nil {
		log.Printf("- ERROR - %s", err.Error())
		return err
	}

	if fast == true {
		log.Printf("- INFO - %s valid according to Payload Oxum", bagLocation)
		return nil
	}

	manifest := filepath.Join(bagLocation, "manifest-sha256.txt")
	e := ValidateManifest(manifest, complete)
	if len(e) > 0 {
		errors = append(errors, e...)
	}

	tagmanifest := filepath.Join(bagLocation, "tagmanifest-sha256.txt")
	e = ValidateManifest(tagmanifest, complete)
	if len(e) > 0 {
		errors = append(errors, e...)
	}

	if len(errors) == 0 {
		log.Printf("- INFO - %s is valid", bagLocation)
	} else {
		errorMsgs := fmt.Sprintf("- ERROR - %s is invalid: Bag validation failed: ", bagLocation)
		for i, e := range errors {
			errorMsgs = errorMsgs + e.Error()
			if i < len(errors)-1 {
				errorMsgs = errorMsgs + "; "
			}
		}
		log.Println(errorMsgs)
	}
	return nil
}

func CreateBag(inputDir string, algorithm string, numProcesses int) error {
	//check that input exists and is a directory
	if err := directoryExists(inputDir); err != nil {
		return err
	}

	log.Printf("- INFO - Creating Bag for directory %s", inputDir)

	//create a slice of files
	filesToBag, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return err
	}

	//check there is at least one file to be bagged.
	if len(filesToBag) < 1 {
		errMsg := fmt.Errorf("Could not create a bag, no files present in %s", inputDir)
		log.Println("- ERROR -", errMsg)
		return errMsg
	}

	//create a data directory for payload
	log.Println("- INFO - Creating data directory")
	dataDirName := filepath.Join(inputDir, "data")
	if err := os.Mkdir(dataDirName, 0777); err != nil {
		log.Println("- ERROR -", err)
		return err
	}

	//move the payload files into data dir
	for _, file := range filesToBag {
		originalLocation := filepath.Join(inputDir, file.Name())
		newLocation := filepath.Join(dataDirName, file.Name())
		log.Printf("- INFO - Moving %s to %s", originalLocation, newLocation)
		if err := os.Rename(originalLocation, newLocation); err != nil {
			log.Println("- ERROR -", err.Error())
			return err
		}
	}

	//Generate the manifest
	if err := CreateManifest(dataDirName, algorithm, numProcesses); err != nil {
		return err
	}
	return nil
}

func directoryExists(inputDir string) error {
	if fi, err := os.Stat(inputDir); err == nil {
		if fi.IsDir() == true {
			return nil
		} else {
			errorMsg := fmt.Errorf("Failed to create bag: input directory %s is not a directory", inputDir)
			log.Println("- Error -", errorMsg)
			return errorMsg
		}
	} else if os.IsNotExist(err) {
		errorMsg := fmt.Errorf(" - ERROR - input %s directory does not exist", inputDir)
		log.Println(errorMsg)
		return err
	} else {
		return err
	}
}
