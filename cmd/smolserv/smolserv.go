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

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"

	"github.com/lucasreed/smol/pkg/app"
	"github.com/lucasreed/smol/pkg/data"
	"github.com/lucasreed/smol/pkg/storage/boltdb"
	"github.com/lucasreed/smol/pkg/storage/mysql"
	"github.com/lucasreed/smol/pkg/storage/rediscache"
)

type Config struct {
	ListenIp      string `default:"0.0.0.0"`
	ListenPort    string `default:"8080"`
	BoltdbPath    string `default:"./boltdb"`
	MysqlDatabase string `default:"smol"`
	MysqlHost     string `default:"localhost"`
	MysqlPort     string `default:"3306"`
	MysqlUser     string `default:"smol"`
	MysqlPassword string `default:"smol"`
	RedisHost     string `default:"localhost"`
	RedisPort     string `default:"6379"`
}

var (
	conf        Config
	storageType string
	version     = "development"
	commit      = "n/a"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.Flags().StringVar(&storageType, "storage", "boltdb", "What storage backend to use. Valid options: redis, boltdb")
}

var rootCmd = &cobra.Command{
	Use:   "smolserv",
	Short: "smolserv makes urls smol",
	Long:  `A simple url shortener API written in go`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return envconfig.Process("smol", &conf)
	},
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := setupStorage()
		if err != nil {
			log.Fatal("error setting up storage - ", err)
		}
		app := app.NewServer(storage, conf.ListenIp+":"+conf.ListenPort)
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
		bolt := boltdb.NewStore(conf.BoltdbPath)
		err := bolt.Open()
		if err != nil {
			return nil, err
		}
		store = bolt
	case "redis":
		redisStore := rediscache.NewStore(conf.RedisHost, conf.RedisPort)
		err := redisStore.Open()
		if err != nil {
			return nil, err
		}
		store = redisStore
	case "mysql":
		mysqlStore := mysql.NewStore(conf.MysqlHost, conf.MysqlPort, conf.MysqlUser, conf.MysqlPassword, conf.MysqlDatabase)
		err := mysqlStore.Open()
		if err != nil {
			return nil, err
		}
		store = mysqlStore
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
