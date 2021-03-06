// Copyright 2017 Baptist-Publication Information Technology Services Co.,Ltd.
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
// limitations under the License.

package expr

import (
	"strings"
)

func Compile(exprStr string) ([]byte, error) {
	exprObj, err := ParseReader("", strings.NewReader(exprStr))
	if err != nil {
		return nil, err
	}
	bz, err := exprObj.(Byteful).Bytes()
	if err != nil {
		return nil, err
	}
	return bz, err
}
