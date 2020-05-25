// Copyright 2020 Luke Reed <luke@lreed.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/lucasreed/smol/pkg/app"
	"github.com/lucasreed/smol/pkg/data"
	"github.com/lucasreed/smol/pkg/storage/boltdb"
	"github.com/lucasreed/smol/pkg/storage/rediscache"
)

var (
	boltdbPath  string
	listen      string
	listenPort  string
	redisHost   string
	redisPort   string
	storageType string
	version     = "development"
	commit      = "n/a"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.Flags().StringVarP(&listen, "listen-ip", "i", "0.0.0.0", "IP to listen on")
	rootCmd.Flags().StringVarP(&listenPort, "listen-port", "p", "8080", "port to listen on")
	rootCmd.Flags().StringVar(&storageType, "storage", "boltdb", "What storage backend to use. Valid options: redis, boltdb")
	rootCmd.Flags().StringVar(&boltdbPath, "boltdb-path", "./boltdb", "location of boltdb file")
	rootCmd.Flags().StringVar(&redisHost, "redis-host", "localhost", "hostname/IP of redis")
	rootCmd.Flags().StringVar(&redisPort, "redis-port", "6379", "port redis is listening on")
}

var rootCmd = &cobra.Command{
	Use:   "smolserv",
	Short: "smolserv makes urls smol",
	Long:  `A simple url shortener API written in go`,
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := setupStorage()
		if err != nil {
			log.Fatal("error setting up storage - ", err)
		}
		app := app.NewServer(storage, listen+":"+listenPort)
		app.Run()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the current version of the tool.",
	Long:  `Prints the current version.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:" + version + " Commit:" + commit)
	},
}

func setupStorage() (data.StorageReadWrite, error) {
	var store data.StorageReadWrite
	switch storageType {
	case "boltdb":
		bolt := boltdb.NewStore(boltdbPath)
		err := bolt.Open()
		if err != nil {
			return nil, err
		}
		store = bolt
	case "redis":
		redisStore := rediscache.NewStore(redisHost, redisPort)
		err := redisStore.Open()
		if err != nil {
			return nil, err
		}
		store = redisStore
	default:
		return nil, fmt.Errorf("not a valid storage backend: %v", storageType)
	}
	return store, nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("error starting smolserv - ", err)
	}
}
