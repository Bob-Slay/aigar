//  Copyright 2018 The go-ethereum Authors
//  Copyright 2019 The go-aigar Authors
//  This file is part of the go-aigar library.
//
//  The go-aigar library is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Lesser General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  The go-aigar library is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
//  GNU Lesser General Public License for more details.
//
//  You should have received a copy of the GNU Lesser General Public License
//  along with the go-aigar library. If not, see <http://www.gnu.org/licenses/>.

// +build gofuzz

package blake2b

import (
	"encoding/binary"
)

func Fuzz(data []byte) int {
	// Make sure the data confirms to the input model
	if len(data) != 211 {
		return 0
	}
	// Parse everything and call all the implementations
	var (
		rounds = binary.BigEndian.Uint16(data[0:2])

		h [8]uint64
		m [16]uint64
		t [2]uint64
		f uint64
	)
	for i := 0; i < 8; i++ {
		offset := 2 + i*8
		h[i] = binary.LittleEndian.Uint64(data[offset : offset+8])
	}
	for i := 0; i < 16; i++ {
		offset := 66 + i*8
		m[i] = binary.LittleEndian.Uint64(data[offset : offset+8])
	}
	t[0] = binary.LittleEndian.Uint64(data[194:202])
	t[1] = binary.LittleEndian.Uint64(data[202:210])

	if data[210]%2 == 1 { // Avoid spinning the fuzzer to hit 0/1
		f = 0xFFFFFFFFFFFFFFFF
	}
	// Run the blake2b compression on all instruction sets and cross reference
	want := h
	fGeneric(&want, &m, t[0], t[1], f, uint64(rounds))

	have := h
	fSSE4(&have, &m, t[0], t[1], f, uint64(rounds))
	if have != want {
		panic("SSE4 mismatches generic algo")
	}
	have = h
	fAVX(&have, &m, t[0], t[1], f, uint64(rounds))
	if have != want {
		panic("AVX mismatches generic algo")
	}
	have = h
	fAVX2(&have, &m, t[0], t[1], f, uint64(rounds))
	if have != want {
		panic("AVX2 mismatches generic algo")
	}
	return 1
}
