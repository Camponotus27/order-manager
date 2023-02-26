package file

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"order-manager/app/domain/model"
	"order-manager/app/domain/repository"
)

const (
	initialPostfix = 1
)

type Config struct {
	PathStoreNote string
	PathTempNote  string
}

type fileRepository struct {
	config *Config
}

func NewFileRepository(config *Config) repository.FileRepository {
	return &fileRepository{config: config}
}

func (r fileRepository) RenameFileWithSubContext(subContext string) error {
	pathRenameFile := r.config.PathTempNote
	files, err := os.ReadDir(pathRenameFile)
	if err != nil {
		return err
	}

	postFixNumber := initialPostfix

	msgErr := ""
	for _, f := range files {
		fileName := f.Name()
		splitPoint := strings.Split(fileName, ".")

		lenSplitPoint := len(splitPoint)
		if lenSplitPoint < 1 {
			continue
		}
		ext := splitPoint[lenSplitPoint-1]
		if ext != "png" {
			continue
		}

		newFileName := fmt.Sprintf("%s_%d.%s", subContext, postFixNumber, ext)

		errRename := os.Rename(filepath.Join(pathRenameFile, fileName), filepath.Join(pathRenameFile, newFileName))
		if errRename != nil {
			msgErr = fmt.Sprintf("%s %s", msgErr, errRename.Error())
			continue
		}

		postFixNumber++
	}

	if msgErr != "" {
		return errors.New(msgErr)
	}
	return nil
}

func (r fileRepository) CreateFolderFromTask(path string, task *model.Note) ([]string, error) {
	pathAbsolute := filepath.Join(r.config.PathStoreNote, path)
	errMkDirAll := os.MkdirAll(pathAbsolute, os.ModePerm)
	if errMkDirAll != nil {
		return nil, errMkDirAll
	}

	postFixNumber := initialPostfix
	var pathsFile []string
	for _, file := range task.Files {
		f := filepath.Join(pathAbsolute, fmt.Sprintf("%s_%d.%s", file.Name, postFixNumber, "png"))
		err := downloadFile(file.Url, f)
		if err != nil {
			return nil, err
		}
		pathsFile = append(pathsFile, f)
		postFixNumber++
	}

	defer func() {
		cmd := exec.Command("open", pathAbsolute)
		_, _ = cmd.Output()
	}()

	return pathsFile, nil
}

func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
