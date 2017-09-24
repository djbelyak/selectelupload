package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ernado/selectel/storage"
	config "github.com/tj/go-config"
)

type Options struct {
	User      string `help:"selectel user name"`
	Password  string `help:"selectel password"`
	Container string `help:"name of bucket"`
	Dir       string `help:"directory which will be uploaded"`
}

type UploadTask struct {
	File     string
	Attempts uint
}

func main() {
	var options Options
	config.MustResolve(&options)

	api, err := storage.New(options.User, options.Password)
	if err != nil {
		log.Fatal(err)
	}
	objects, err := api.Container(options.Container).Objects()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for _, obj := range objects {
		wg.Add(1)
		go func(obj storage.ObjectAPI) {
			defer wg.Done()
			err := obj.Remove()
			if err != nil {
				log.Printf("Can't remove file: %v", err)
				return
			}
		}(obj)
	}

	wg.Wait()

	upload := make(chan UploadTask, 1048576)

	go func() {
		wg.Add(1)
		for task := range upload {
			wg.Add(1)
			go func(task UploadTask) {
				defer wg.Done()

				if task.Attempts == 0 {
					return
				}
				filename := task.File

				file, err := os.Open(filename)
				if err != nil {
					log.Printf("Can't open file %s: %v", filename, err)
					task.Attempts--
					upload <- task
					return
				}
				defer file.Close()

				err = api.Upload(file, options.Container, strings.TrimPrefix(filepath.ToSlash(filename), options.Dir), "")
				if err != nil {
					log.Printf("Can't upload file %s: %v", filename, err)
					task.Attempts--
					upload <- task
					return
				}

				log.Printf("File %s uploaded", filename)

			}(task)

		}
		wg.Done()
	}()

	err = filepath.Walk(options.Dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && path[0] != '.' {
			upload <- UploadTask{
				File:     path,
				Attempts: 5,
			}
		}

		return nil
	})

	//upload <- filepath.Join(options.Dir, "README.md")
	close(upload)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	wg.Wait()

	log.Println("Well done!")

}
