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

package jsonutil

import (
	"encoding/json"

	simplejson "github.com/bitly/go-simplejson"
	"openpitrix.io/logger"
)

func Encode(o interface{}) ([]byte, error) {
	return json.Marshal(o)
}

func Decode(y []byte, o interface{}) error {
	return json.Unmarshal(y, o)
}

func ToString(o interface{}) string {
	b, err := Encode(o)
	if err != nil {
		logger.Errorf(nil, "Encode [%+v] failed: %+v", o, err)
		return ""
	}
	return string(b)
}

// FIXME: need improve performance
func ToJson(o interface{}) Json {
	var j Json
	j = &fakeJson{simplejson.New()}
	b, err := Encode(o)
	if err != nil {
		logger.Errorf(nil, "Encode [%+v] to []byte failed: %+v", o, err)
		return j
	}
	j, err = NewJson(b)
	if err != nil {
		logger.Errorf(nil, "Decode [%+v] to json failed: %+v", o, err)
	}
	return j
}

type fakeJson struct {
	*simplejson.Json
}

func NewJson(y []byte) (Json, error) {
	j, err := simplejson.NewJson(y)
	return &fakeJson{j}, err
}

func (j *fakeJson) Get(key string) Json {
	return &fakeJson{j.Json.Get(key)}
}

func (j *fakeJson) GetPath(branch ...string) Json {
	return &fakeJson{j.Json.GetPath(branch...)}
}

func (j *fakeJson) CheckGet(key string) (Json, bool) {
	result, ok := j.Json.CheckGet(key)
	return &fakeJson{result}, ok
}

//
//func (j *fakeJson) UnmarshalJSON(p []byte) error {
//	return j.Json.UnmarshalJSON(p)
//}
//
//func (j *fakeJson) MarshalJSON() ([]byte, error) {
//	return j.Json.MarshalJSON()
//}
