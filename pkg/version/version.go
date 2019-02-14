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

//go:generate go run gen_helper.go
//go:generate go fmt

package version

import "fmt"

var (
	ShortVersion   = "dev"
	GitSha1Version = "git-sha1"
	BuildDate      = "2017-01-01"
)

func PrintVersionInfo(printer func(string, ...interface{})) {
	printer("Release OpVersion: %s", ShortVersion)
	printer("Git Commit Hash: %s", GitSha1Version)
	printer("Build Time: %s", BuildDate)
}

func GetVersionString() string {
	return fmt.Sprintf("%s; git: %s; build time: %s", ShortVersion, GitSha1Version, BuildDate)
}
