package pgntodb

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Process ... process a single file or all the files of a folder
func Process(filepath string) {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		log.Fatal("Cannot access " + filepath)
	}

	if info.IsDir() {
		fileinfos, err := ioutil.ReadDir(filepath)
		if err != nil {
			log.Fatal("Cannot list files in " + filepath)
		}
		for _, info := range fileinfos {
			if !info.IsDir() {
				log.Println(path.Join(filepath, info.Name()))
				processFile(path.Join(filepath, info.Name()))
			}
		}
	} else {
		processFile(filepath)
	}

}

// ProcessFile ... does everything
func processFile(filepath string) {

	// Open file
	file, err := os.Open(filepath)
	defer file.Close()

	if err != nil {
		log.Fatal("Cannot open file " + filepath)
	}

	// Connect to DB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Ping MongoDB
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("Cannot connect to DB")
	}

	// Do the work
	pgnFileToDB(file, client)
}