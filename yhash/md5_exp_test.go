package yhash

import (
	"testing"
)

var path = "/Users/vej/Desktop/20191030/hash_test/empty.txt"

func BenchmarkMD5v1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MD5v1(path)
	}
}

func BenchmarkMD5v2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MD5v2(path)
	}
}

func BenchmarkMD5v3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MD5v3(path)
	}
}

func BenchmarkMD5v4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MD5v4(path)
	}
}

func BenchmarkMD5v5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MD5v5(path)
	}
}
