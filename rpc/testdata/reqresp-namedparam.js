/*
 * // Copyright 2018 The go-ethereum Authors
 * // Copyright 2019 The go-aigar Authors
 * // This file is part of the go-aigar library.
 * //
 * // The go-aigar library is free software: you can redistribute it and/or modify
 * // it under the terms of the GNU Lesser General Public License as published by
 * // the Free Software Foundation, either version 3 of the License, or
 * // (at your option) any later version.
 * //
 * // The go-aigar library is distributed in the hope that it will be useful,
 * // but WITHOUT ANY WARRANTY; without even the implied warranty of
 * // MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * // GNU Lesser General Public License for more details.
 * //
 * // You should have received a copy of the GNU Lesser General Public License
 * // along with the go-aigar library. If not, see <http://www.gnu.org/licenses/>.
 */

// This test checks that an error response is sent for calls
// with named parameters.

--> {"jsonrpc":"2.0","method":"test_echo","params":{"int":23},"id":3}
<-- {"jsonrpc":"2.0","id":3,"error":{"code":-32602,"message":"non-array args"}}
