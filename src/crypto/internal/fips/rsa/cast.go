// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rsa

import (
	"bytes"
	"crypto/internal/fips"
	"crypto/internal/fips/bigmod"
	_ "crypto/internal/fips/check"
	"errors"
	"sync"
)

func testPrivateKey() *PrivateKey {
	// https://www.rfc-editor.org/rfc/rfc9500.html#section-2.1
	N, _ := bigmod.NewModulus([]byte{
		0xB0, 0xF9, 0xE8, 0x19, 0x43, 0xA7, 0xAE, 0x98,
		0x92, 0xAA, 0xDE, 0x17, 0xCA, 0x7C, 0x40, 0xF8,
		0x74, 0x4F, 0xED, 0x2F, 0x81, 0x48, 0xE6, 0xC8,
		0xEA, 0xA2, 0x7B, 0x7D, 0x00, 0x15, 0x48, 0xFB,
		0x51, 0x92, 0xAB, 0x28, 0xB5, 0x6C, 0x50, 0x60,
		0xB1, 0x18, 0xCC, 0xD1, 0x31, 0xE5, 0x94, 0x87,
		0x4C, 0x6C, 0xA9, 0x89, 0xB5, 0x6C, 0x27, 0x29,
		0x6F, 0x09, 0xFB, 0x93, 0xA0, 0x34, 0xDF, 0x32,
		0xE9, 0x7C, 0x6F, 0xF0, 0x99, 0x8C, 0xFD, 0x8E,
		0x6F, 0x42, 0xDD, 0xA5, 0x8A, 0xCD, 0x1F, 0xA9,
		0x79, 0x86, 0xF1, 0x44, 0xF3, 0xD1, 0x54, 0xD6,
		0x76, 0x50, 0x17, 0x5E, 0x68, 0x54, 0xB3, 0xA9,
		0x52, 0x00, 0x3B, 0xC0, 0x68, 0x87, 0xB8, 0x45,
		0x5A, 0xC2, 0xB1, 0x9F, 0x7B, 0x2F, 0x76, 0x50,
		0x4E, 0xBC, 0x98, 0xEC, 0x94, 0x55, 0x71, 0xB0,
		0x78, 0x92, 0x15, 0x0D, 0xDC, 0x6A, 0x74, 0xCA,
		0x0F, 0xBC, 0xD3, 0x54, 0x97, 0xCE, 0x81, 0x53,
		0x4D, 0xAF, 0x94, 0x18, 0x84, 0x4B, 0x13, 0xAE,
		0xA3, 0x1F, 0x9D, 0x5A, 0x6B, 0x95, 0x57, 0xBB,
		0xDF, 0x61, 0x9E, 0xFD, 0x4E, 0x88, 0x7F, 0x2D,
		0x42, 0xB8, 0xDD, 0x8B, 0xC9, 0x87, 0xEA, 0xE1,
		0xBF, 0x89, 0xCA, 0xB8, 0x5E, 0xE2, 0x1E, 0x35,
		0x63, 0x05, 0xDF, 0x6C, 0x07, 0xA8, 0x83, 0x8E,
		0x3E, 0xF4, 0x1C, 0x59, 0x5D, 0xCC, 0xE4, 0x3D,
		0xAF, 0xC4, 0x91, 0x23, 0xEF, 0x4D, 0x8A, 0xBB,
		0xA9, 0x3D, 0x39, 0x05, 0xE4, 0x02, 0x8D, 0x7B,
		0xA9, 0x14, 0x84, 0xA2, 0x75, 0x96, 0xE0, 0x7B,
		0x4B, 0x6E, 0xD9, 0x92, 0xF0, 0x77, 0xB5, 0x24,
		0xD3, 0xDC, 0xFE, 0x7D, 0xDD, 0x55, 0x49, 0xBE,
		0x7C, 0xCE, 0x8D, 0xA0, 0x35, 0xCF, 0xA0, 0xB3,
		0xFB, 0x8F, 0x9E, 0x46, 0xF7, 0x32, 0xB2, 0xA8,
		0x6B, 0x46, 0x01, 0x65, 0xC0, 0x8F, 0x53, 0x13})
	d, _ := bigmod.NewNat().SetBytes([]byte{
		0x41, 0x18, 0x8B, 0x20, 0xCF, 0xDB, 0xDB, 0xC2,
		0xCF, 0x1F, 0xFE, 0x75, 0x2D, 0xCB, 0xAA, 0x72,
		0x39, 0x06, 0x35, 0x2E, 0x26, 0x15, 0xD4, 0x9D,
		0xCE, 0x80, 0x59, 0x7F, 0xCF, 0x0A, 0x05, 0x40,
		0x3B, 0xEF, 0x00, 0xFA, 0x06, 0x51, 0x82, 0xF7,
		0x2D, 0xEC, 0xFB, 0x59, 0x6F, 0x4B, 0x0C, 0xE8,
		0xFF, 0x59, 0x70, 0xBA, 0xF0, 0x7A, 0x89, 0xA5,
		0x19, 0xEC, 0xC8, 0x16, 0xB2, 0xF4, 0xFF, 0xAC,
		0x50, 0x69, 0xAF, 0x1B, 0x06, 0xBF, 0xEF, 0x7B,
		0xF6, 0xBC, 0xD7, 0x9E, 0x4E, 0x81, 0xC8, 0xC5,
		0xA3, 0xA7, 0xD9, 0x13, 0x0D, 0xC3, 0xCF, 0xBA,
		0xDA, 0xE5, 0xF6, 0xD2, 0x88, 0xF9, 0xAE, 0xE3,
		0xF6, 0xFF, 0x92, 0xFA, 0xE0, 0xF8, 0x1A, 0xF5,
		0x97, 0xBE, 0xC9, 0x6A, 0xE9, 0xFA, 0xB9, 0x40,
		0x2C, 0xD5, 0xFE, 0x41, 0xF7, 0x05, 0xBE, 0xBD,
		0xB4, 0x7B, 0xB7, 0x36, 0xD3, 0xFE, 0x6C, 0x5A,
		0x51, 0xE0, 0xE2, 0x07, 0x32, 0xA9, 0x7B, 0x5E,
		0x46, 0xC1, 0xCB, 0xDB, 0x26, 0xD7, 0x48, 0x54,
		0xC6, 0xB6, 0x60, 0x4A, 0xED, 0x46, 0x37, 0x35,
		0xFF, 0x90, 0x76, 0x04, 0x65, 0x57, 0xCA, 0xF9,
		0x49, 0xBF, 0x44, 0x88, 0x95, 0xC2, 0x04, 0x32,
		0xC1, 0xE0, 0x9C, 0x01, 0x4E, 0xA7, 0x56, 0x60,
		0x43, 0x4F, 0x1A, 0x0F, 0x3B, 0xE2, 0x94, 0xBA,
		0xBC, 0x5D, 0x53, 0x0E, 0x6A, 0x10, 0x21, 0x3F,
		0x53, 0xB6, 0x03, 0x75, 0xFC, 0x84, 0xA7, 0x57,
		0x3F, 0x2A, 0xF1, 0x21, 0x55, 0x84, 0xF5, 0xB4,
		0xBD, 0xA6, 0xD4, 0xE8, 0xF9, 0xE1, 0x7A, 0x78,
		0xD9, 0x7E, 0x77, 0xB8, 0x6D, 0xA4, 0xA1, 0x84,
		0x64, 0x75, 0x31, 0x8A, 0x7A, 0x10, 0xA5, 0x61,
		0x01, 0x4E, 0xFF, 0xA2, 0x3A, 0x81, 0xEC, 0x56,
		0xE9, 0xE4, 0x10, 0x9D, 0xEF, 0x8C, 0xB3, 0xF7,
		0x97, 0x22, 0x3F, 0x7D, 0x8D, 0x0D, 0x43, 0x51}, N)
	p, _ := bigmod.NewModulus([]byte{
		0xDD, 0x10, 0x57, 0x02, 0x38, 0x2F, 0x23, 0x2B,
		0x36, 0x81, 0xF5, 0x37, 0x91, 0xE2, 0x26, 0x17,
		0xC7, 0xBF, 0x4E, 0x9A, 0xCB, 0x81, 0xED, 0x48,
		0xDA, 0xF6, 0xD6, 0x99, 0x5D, 0xA3, 0xEA, 0xB6,
		0x42, 0x83, 0x9A, 0xFF, 0x01, 0x2D, 0x2E, 0xA6,
		0x28, 0xB9, 0x0A, 0xF2, 0x79, 0xFD, 0x3E, 0x6F,
		0x7C, 0x93, 0xCD, 0x80, 0xF0, 0x72, 0xF0, 0x1F,
		0xF2, 0x44, 0x3B, 0x3E, 0xE8, 0xF2, 0x4E, 0xD4,
		0x69, 0xA7, 0x96, 0x13, 0xA4, 0x1B, 0xD2, 0x40,
		0x20, 0xF9, 0x2F, 0xD1, 0x10, 0x59, 0xBD, 0x1D,
		0x0F, 0x30, 0x1B, 0x5B, 0xA7, 0xA9, 0xD3, 0x63,
		0x7C, 0xA8, 0xD6, 0x5C, 0x1A, 0x98, 0x15, 0x41,
		0x7D, 0x8E, 0xAB, 0x73, 0x4B, 0x0B, 0x4F, 0x3A,
		0x2C, 0x66, 0x1D, 0x9A, 0x1A, 0x82, 0xF3, 0xAC,
		0x73, 0x4C, 0x40, 0x53, 0x06, 0x69, 0xAB, 0x8E,
		0x47, 0x30, 0x45, 0xA5, 0x8E, 0x65, 0x53, 0x9D})
	q, _ := bigmod.NewModulus([]byte{
		0xCC, 0xF1, 0xE5, 0xBB, 0x90, 0xC8, 0xE9, 0x78,
		0x1E, 0xA7, 0x5B, 0xEB, 0xF1, 0x0B, 0xC2, 0x52,
		0xE1, 0x1E, 0xB0, 0x23, 0xA0, 0x26, 0x0F, 0x18,
		0x87, 0x55, 0x2A, 0x56, 0x86, 0x3F, 0x4A, 0x64,
		0x21, 0xE8, 0xC6, 0x00, 0xBF, 0x52, 0x3D, 0x6C,
		0xB1, 0xB0, 0xAD, 0xBD, 0xD6, 0x5B, 0xFE, 0xE4,
		0xA8, 0x8A, 0x03, 0x7E, 0x3D, 0x1A, 0x41, 0x5E,
		0x5B, 0xB9, 0x56, 0x48, 0xDA, 0x5A, 0x0C, 0xA2,
		0x6B, 0x54, 0xF4, 0xA6, 0x39, 0x48, 0x52, 0x2C,
		0x3D, 0x5F, 0x89, 0xB9, 0x4A, 0x72, 0xEF, 0xFF,
		0x95, 0x13, 0x4D, 0x59, 0x40, 0xCE, 0x45, 0x75,
		0x8F, 0x30, 0x89, 0x80, 0x90, 0x89, 0x56, 0x58,
		0x8E, 0xEF, 0x57, 0x5B, 0x3E, 0x4B, 0xC4, 0xC3,
		0x68, 0xCF, 0xE8, 0x13, 0xEE, 0x9C, 0x25, 0x2C,
		0x2B, 0x02, 0xE0, 0xDF, 0x91, 0xF1, 0xAA, 0x01,
		0x93, 0x8D, 0x38, 0x68, 0x5D, 0x60, 0xBA, 0x6F})
	qInv, _ := bigmod.NewNat().SetBytes([]byte{
		0x0A, 0x81, 0xD8, 0xA6, 0x18, 0x31, 0x4A, 0x80,
		0x3A, 0xF6, 0x1C, 0x06, 0x71, 0x1F, 0x2C, 0x39,
		0xB2, 0x66, 0xFF, 0x41, 0x4D, 0x53, 0x47, 0x6D,
		0x1D, 0xA5, 0x2A, 0x43, 0x18, 0xAA, 0xFE, 0x4B,
		0x96, 0xF0, 0xDA, 0x07, 0x15, 0x5F, 0x8A, 0x51,
		0x34, 0xDA, 0xB8, 0x8E, 0xE2, 0x9E, 0x81, 0x68,
		0x07, 0x6F, 0xCD, 0x78, 0xCA, 0x79, 0x1A, 0xC6,
		0x34, 0x42, 0xA8, 0x1C, 0xD0, 0x69, 0x39, 0x27,
		0xD8, 0x08, 0xE3, 0x35, 0xE8, 0xD8, 0xCB, 0xF2,
		0x12, 0x19, 0x07, 0x50, 0x9A, 0x57, 0x75, 0x9B,
		0x4F, 0x9A, 0x18, 0xFA, 0x3A, 0x7B, 0x33, 0x37,
		0x79, 0xED, 0xDE, 0x7A, 0x45, 0x93, 0x84, 0xF8,
		0x44, 0x4A, 0xDA, 0xEC, 0xFF, 0xEC, 0x95, 0xFD,
		0x55, 0x2B, 0x0C, 0xFC, 0xB6, 0xC7, 0xF6, 0x92,
		0x62, 0x6D, 0xDE, 0x1E, 0xF2, 0x68, 0xA4, 0x0D,
		0x2F, 0x67, 0xB5, 0xC8, 0xAA, 0x38, 0x7F, 0xF7}, p)
	dP := []byte{
		0x09, 0xED, 0x54, 0xEA, 0xED, 0x98, 0xF8, 0x4C,
		0x55, 0x7B, 0x4A, 0x86, 0xBF, 0x4F, 0x57, 0x84,
		0x93, 0xDC, 0xBC, 0x6B, 0xE9, 0x1D, 0xA1, 0x89,
		0x37, 0x04, 0x04, 0xA9, 0x08, 0x72, 0x76, 0xF4,
		0xCE, 0x51, 0xD8, 0xA1, 0x00, 0xED, 0x85, 0x7D,
		0xC2, 0xB0, 0x64, 0x94, 0x74, 0xF3, 0xF1, 0x5C,
		0xD2, 0x4C, 0x54, 0xDB, 0x28, 0x71, 0x10, 0xE5,
		0x6E, 0x5C, 0xB0, 0x08, 0x68, 0x2F, 0x91, 0x68,
		0xAA, 0x81, 0xF3, 0x14, 0x58, 0xB7, 0x43, 0x1E,
		0xCC, 0x1C, 0x44, 0x90, 0x6F, 0xDA, 0x87, 0xCA,
		0x89, 0x47, 0x10, 0xC3, 0x71, 0xE9, 0x07, 0x6C,
		0x1D, 0x49, 0xFB, 0xAE, 0x51, 0x27, 0x69, 0x34,
		0xF2, 0xAD, 0x78, 0x77, 0x89, 0xF4, 0x2D, 0x0F,
		0xA0, 0xB4, 0xC9, 0x39, 0x85, 0x5D, 0x42, 0x12,
		0x09, 0x6F, 0x70, 0x28, 0x0A, 0x4E, 0xAE, 0x7C,
		0x8A, 0x27, 0xD9, 0xC8, 0xD0, 0x77, 0x2E, 0x65}
	dQ := []byte{
		0x8C, 0xB6, 0x85, 0x7A, 0x7B, 0xD5, 0x46, 0x5F,
		0x80, 0x04, 0x7E, 0x9B, 0x87, 0xBC, 0x00, 0x27,
		0x31, 0x84, 0x05, 0x81, 0xE0, 0x62, 0x61, 0x39,
		0x01, 0x2A, 0x5B, 0x50, 0x5F, 0x0A, 0x33, 0x84,
		0x7E, 0xB7, 0xB8, 0xC3, 0x28, 0x99, 0x49, 0xAD,
		0x48, 0x6F, 0x3B, 0x4B, 0x3D, 0x53, 0x9A, 0xB5,
		0xDA, 0x76, 0x30, 0x21, 0xCB, 0xC8, 0x2C, 0x1B,
		0xA2, 0x34, 0xA5, 0x66, 0x8D, 0xED, 0x08, 0x01,
		0xB8, 0x59, 0xF3, 0x43, 0xF1, 0xCE, 0x93, 0x04,
		0xE6, 0xFA, 0xA2, 0xB0, 0x02, 0xCA, 0xD9, 0xB7,
		0x8C, 0xDE, 0x5C, 0xDC, 0x2C, 0x1F, 0xB4, 0x17,
		0x1C, 0x42, 0x42, 0x16, 0x70, 0xA6, 0xAB, 0x0F,
		0x50, 0xCC, 0x4A, 0x19, 0x4E, 0xB3, 0x6D, 0x1C,
		0x91, 0xE9, 0x35, 0xBA, 0x01, 0xB9, 0x59, 0xD8,
		0x72, 0x8B, 0x9E, 0x64, 0x42, 0x6B, 0x3F, 0xC3,
		0xA7, 0x50, 0x6D, 0xEB, 0x52, 0x39, 0xA8, 0xA7}
	return &PrivateKey{
		pub: PublicKey{
			N: N, E: 65537,
		},
		d: d, p: p, q: q, qInv: qInv, dP: dP, dQ: dQ,
	}

}

func testHash() []byte {
	return []byte{
		0x17, 0x1b, 0x1f, 0x5e, 0x9f, 0x8f, 0x8c, 0x5c,
		0x42, 0xe8, 0x06, 0x59, 0x7b, 0x54, 0xc7, 0xb4,
		0x49, 0x05, 0xa1, 0xdb, 0x3a, 0x3c, 0x31, 0xd3,
		0xb7, 0x56, 0x45, 0x8c, 0xc2, 0xd6, 0x88, 0x62,
	}
}

var fipsSelfTest = sync.OnceFunc(func() {
	fips.CAST("RSASSA-PKCS-v1.5 2048-bit sign and verify", func() error {
		k := testPrivateKey()
		hash := []byte{
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
			0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
			0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
			0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
		}
		want := []byte{
			0x16, 0x98, 0x33, 0xc7, 0x30, 0x2c, 0x0a, 0xdc,
			0x0a, 0x8d, 0x02, 0x58, 0xeb, 0xf9, 0x7d, 0xb6,
			0x2a, 0xad, 0xee, 0x63, 0x72, 0xaa, 0x37, 0x2c,
			0xb3, 0x06, 0x04, 0xdf, 0xdb, 0x2b, 0xbc, 0xb1,
			0x76, 0x3e, 0xeb, 0x87, 0xef, 0x91, 0xef, 0x74,
			0x69, 0x62, 0x27, 0xf3, 0x24, 0xf8, 0xe7, 0x0e,
			0xb2, 0x15, 0x3f, 0xa2, 0x4d, 0xe2, 0x0c, 0xd4,
			0xdc, 0x2d, 0xc1, 0x1a, 0x84, 0x7c, 0x88, 0x80,
			0xb9, 0xa9, 0x23, 0x67, 0x39, 0x2e, 0x86, 0xc0,
			0x53, 0x9b, 0xc1, 0x35, 0xb3, 0x17, 0x5e, 0x62,
			0x95, 0xd6, 0xbc, 0x2a, 0xa6, 0xb1, 0xcf, 0x8f,
			0x99, 0x43, 0x1f, 0x3d, 0xd2, 0x70, 0x3f, 0x01,
			0x37, 0x2b, 0xdd, 0x69, 0x1a, 0x5c, 0x2b, 0x04,
			0x70, 0x92, 0xea, 0x2d, 0x86, 0x00, 0xcb, 0x79,
			0xca, 0xaf, 0xa4, 0x1c, 0xd9, 0x61, 0x21, 0x3b,
			0x1e, 0xc5, 0x88, 0xfb, 0xff, 0xbd, 0xc7, 0x3c,
			0x36, 0xa1, 0xc6, 0x85, 0x03, 0xaf, 0x47, 0x4f,
			0x42, 0x9e, 0x23, 0x65, 0x24, 0x69, 0x17, 0xdb,
			0xe7, 0xb7, 0xdc, 0x51, 0xc6, 0x30, 0x40, 0x32,
			0x4f, 0x71, 0xf1, 0x62, 0x2d, 0xaa, 0x98, 0xdb,
			0x11, 0x14, 0xf9, 0x9c, 0x35, 0xc3, 0x16, 0xe1,
			0x1a, 0xd1, 0x8c, 0x4d, 0x8c, 0xad, 0x06, 0x34,
			0xd2, 0x84, 0x97, 0xa4, 0x0b, 0x6e, 0x6d, 0x19,
			0x9f, 0xa7, 0x40, 0x1e, 0xb5, 0xfc, 0x4e, 0x12,
			0x08, 0xec, 0xf4, 0x07, 0x13, 0xdc, 0x5a, 0x8c,
			0xd5, 0x2a, 0xd6, 0x5a, 0x2c, 0xc9, 0x54, 0x84,
			0x78, 0x34, 0x8f, 0x11, 0xfb, 0x6e, 0xd4, 0x27,
			0x45, 0xd9, 0xfa, 0x90, 0x82, 0x83, 0x73, 0x22,
			0x15, 0xab, 0x96, 0x13, 0x0d, 0x52, 0x1c, 0xdc,
			0x17, 0xde, 0x12, 0x6f, 0x84, 0x46, 0xbb, 0xec,
			0xe3, 0xb1, 0xa1, 0x5d, 0x8b, 0xeb, 0xe6, 0xae,
			0x02, 0xb8, 0x76, 0x47, 0x76, 0x11, 0x61, 0x2b,
		}
		sig, err := signPKCS1v15(k, "SHA-256", hash)
		if err != nil {
			return err
		}
		if err := verifyPKCS1v15(k.PublicKey(), "SHA-256", hash, sig); err != nil {
			return err
		}
		if !bytes.Equal(sig, want) {
			return errors.New("unexpected result")
		}
		return nil
	})
})
