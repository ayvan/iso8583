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

func lbcd(data []byte) []byte {
	if len(data)%2 != 0 {
		return bcd(append(data, "\x00"...))
	}
	return bcd(data)
}

func rbcd(data []byte) []byte {
	if len(data)%2 != 0 {
		return bcd(append([]byte("\x00"), data...))
	}
	return bcd(data)
}

// Encode numeric in ascii into bsd (be sure len(data) % 2 == 0)
func bcd(data []byte) []byte {
	if len(data)%2 != 0 {
		panic("length of raw data must be even")
	}
	out := make([]byte, len(data)/2)
	for i, j := 0, 0; i < len(out); i++ {
		out[i] = ((data[j] & 0x0f) << 4) | (data[j+1] & 0x0f)
		j += 2
	}
	return out
}

func bcdl2Ascii(data []byte, length int) []byte {
	return bcd2Ascii(data)[:length]
}

func bcdr2Ascii(data []byte, length int) []byte {
	out := bcd2Ascii(data)
	return out[len(out)-length:]
}

func bcd2Ascii(data []byte) []byte {
	outLen := len(data) * 2
	out := make([]byte, outLen)
	for i := 0; i < outLen; i++ {
		bcdIndex := i / 2
		if i%2 == 0 {
			// higher order bits to ascii:
			out[i] = (data[bcdIndex] >> 4) + '0'
		} else {
			// lower order bits to ascii:
			out[i] = (data[bcdIndex] & 0x0f) + '0'
		}
	}
	return out
}
