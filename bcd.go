//Copyright 2015 ideazxy
//
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package iso8583

import (
	"encoding/hex"
)

func lbcd(data []byte) []byte {
	if len(data)%2 != 0 {
		return bcd(append(data, "0"...))
	}
	return bcd(data)
}

func rbcd(data []byte) []byte {
	if len(data)%2 != 0 {
		return bcd(append([]byte("0"), data...))
	}
	return bcd(data)
}

// Encode numeric in ascii into bsd (be sure len(data) % 2 == 0)
func bcd(data []byte) []byte {
	out := make([]byte, len(data)/2+1)
	n, err := hex.Decode(out, data)
	if err != nil {
		panic(err.Error())
	}
	return out[:n]
}

func bcdl2Ascii(data []byte, length int) []byte {
	return bcd2Ascii(data)[:length]
}

func bcdr2Ascii(data []byte, length int) []byte {
	out := bcd2Ascii(data)
	return out[len(out)-length:]
}

func bcd2Ascii(data []byte) []byte {
	out := make([]byte, len(data)*2)
	n := hex.Encode(out, data)
	return out[:n]
}
