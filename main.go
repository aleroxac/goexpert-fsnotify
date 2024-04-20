package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

type DBConfig struct {
	DB       string `json:"db"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
}

var config DBConfig
var config_file_name = "config.json"

func MarshalConfig(file string) {
	// parece que os editores fazem o save muito rápido e trunca o arquivo
	// por alguns instantes, para contornar este problema, o sleep abaixo serviu bem
	time.Sleep(100 * time.Millisecond)

	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	MarshalConfig(config_file_name)
	fmt.Println(config)

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// quando edito o arquivo via vim, no primeiro o evento é RENAME, os seguintes não printam mais nada
				// tanto via nano e micro e quanto via vscode, são retornados 2 eventos de WRITE
				fmt.Printf("event: %v\n", event)

				if event.Op&fsnotify.Write == fsnotify.Write {
					MarshalConfig(config_file_name)
					fmt.Printf("modified file: %v\n", event.Name)
					fmt.Println(config)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Printf("error: %v\n", err)
			}
		}
	}()
	err = watcher.Add(config_file_name)
	if err != nil {
		panic(err)
	}
	<-done
}
