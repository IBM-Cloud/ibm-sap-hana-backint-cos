// Copyright 2026 IBM Corp.
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

package logging

import (
	"math"
	"strconv"
)

/*
Generating the printable size
*/
func printableSize(value int64) string {
	if value == 0 {
		return "0B"
	}
	sizeUnits := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	base := math.Log(float64(value)) / math.Log(1024)
	size := round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	unit := sizeUnits[int(math.Floor(base))]
	return strconv.FormatFloat(size, 'f', -1, 64) + " " + string(unit)
}

/*
Rounding for printable size
*/
func round(value float64, roundOn float64, places int) float64 {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * value
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal := round / pow
	return newVal
}
