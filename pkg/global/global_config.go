/*
Copyright 2019 The KubeSphere Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package global

import (
	"sync"

	"openpitrix.io/logger"

	"kubesphere.io/im/pkg/config"
	"kubesphere.io/im/pkg/db"
)

var global *Config
var globalMutex sync.RWMutex

func SetGlobal(config *config.Config) {
	globalMutex.Lock()
	global = NewConfig(config)
	globalMutex.Unlock()
}

func Global() *Config {
	globalMutex.RLock()
	defer globalMutex.RUnlock()
	return global
}

type Config struct {
	Config   *config.Config
	Database *db.Database
}

func NewConfig(config *config.Config) *Config {
	c := &Config{Config: config}
	c.openDatabase()

	return c
}

func (c *Config) openDatabase() {
	database, err := db.OpenDatabase(c.Config)
	if err != nil {
		logger.Criticalf(nil, "failed to connect database")
		panic(err)
	}
	c.Database = database
}
