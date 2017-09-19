package main

import (
	"log"
	"os"
	"path/filepath"
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
	for _, obj := range objects {
		err = obj.Remove()
		if err != nil {
			log.Printf("Can't remove file: %v", err)
			continue
		}
	}

	var wg sync.WaitGroup
	upload := make(chan string, 1048576)

	go func() {
		wg.Add(1)
		for file := range upload {
			wg.Add(1)
			go func(filename string) {
				defer wg.Done()

				file, err := os.Open(filename)
				if err != nil {
					log.Printf("Can't open file %s: %v", filename, err)
					return
				}
				defer file.Close()

				err = api.Upload(file, options.Container, filepath.ToSlash(filename), "")
				if err != nil {
					log.Printf("Can't upload file %s: %v", filename, err)
					return
				}

				log.Printf("File %s uploaded", filename)

			}(file)

		}
		wg.Done()
	}()

	err = filepath.Walk(options.Dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && path[0] != '.' {
			upload <- path
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
