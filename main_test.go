package main

import "testing"

func TestBuildSlowAndFastMatch(t *testing.T) {
	for _, n := range []int{1, 2, 10, 100, 1000} {
		if got, want := buildFast(n), buildSlow(n); got != want {
			t.Fatalf("implementations differ for n=%d", n)
		}
	}
}

func BenchmarkBuildSlow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = buildSlow(1000)
	}
}

func BenchmarkBuildFast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = buildFast(1000)
	}
}
