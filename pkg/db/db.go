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

package db

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"openpitrix.io/logger"

	"kubesphere.io/im/pkg/config"
)

type Database struct {
	cfg *config.Config
	*gorm.DB
}

type Options struct {
	SqlInitDB    []string
	SqlInitTable []string
	SqlInitData  []string
}

func OpenDatabase(cfg *config.Config) (*Database, error) {
	cfg = cfg.Clone()

	logger.Infof(nil, "DB config: begin")
	logger.Infof(nil, "\tType: %s", cfg.DB.Type)
	logger.Infof(nil, "\tHost: %s", cfg.DB.Host)
	logger.Infof(nil, "\tPort: %d", cfg.DB.Port)
	logger.Infof(nil, "\tUser: %s", cfg.DB.User)
	logger.Infof(nil, "\tDatabase: %s", cfg.DB.Database)
	logger.Infof(nil, "DB config: end")

	var p = &Database{cfg: cfg}
	var err error

	p.DB, err = gorm.Open(cfg.DB.Type, cfg.DB.GetUrl())
	if err != nil {
		return nil, err
	}

	// Enable Logger, show detailed log
	p.DB.LogMode(true)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	p.DB.DB().SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	p.DB.DB().SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	p.DB.DB().SetConnMaxLifetime(time.Hour)

	return p, nil
}

func (p *Database) Close() error {
	return p.DB.Close()
}
