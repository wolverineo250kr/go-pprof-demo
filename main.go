package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "net/http/pprof"
)

const defaultN = 12000
const maxN = 50000

func buildSlow(n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += strconv.Itoa(i) + ","
	}
	return result
}

func buildFast(n int) string {
	var b strings.Builder
	b.Grow(n * 7)
	for i := 0; i < n; i++ {
		b.WriteString(strconv.Itoa(i))
		b.WriteByte(',')
	}
	return b.String()
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/compute", func(w http.ResponseWriter, r *http.Request) {
		n := defaultN
		if raw := r.URL.Query().Get("n"); raw != "" {
			parsed, err := strconv.Atoi(raw)
			if err != nil || parsed < 1 || parsed > maxN {
				http.Error(w, fmt.Sprintf("n должен быть между 1 и %d", maxN), http.StatusBadRequest)
				return
			}
			n = parsed
		}

		mode := r.URL.Query().Get("mode")
		start := time.Now()
		var result string
		switch mode {
		case "", "slow":
			mode = "slow"
			result = buildSlow(n)
		case "fast":
			result = buildFast(n)
		default:
			http.Error(w, "режим должен быть медленным или быстрым", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Compute-Mode", mode)
		w.Header().Set("X-Compute-Duration", time.Since(start).String())
		fmt.Fprintf(w, "mode=%s n=%d bytes=%d\n", mode, n, len(result))
	})

	// net/http/pprof registers its handlers on DefaultServeMux.
	mux.Handle("/debug/pprof/", http.DefaultServeMux)
	mux.Handle("/debug/pprof/cmdline", http.DefaultServeMux)
	mux.Handle("/debug/pprof/profile", http.DefaultServeMux)
	mux.Handle("/debug/pprof/symbol", http.DefaultServeMux)
	mux.Handle("/debug/pprof/trace", http.DefaultServeMux)

	addr := ":8080"
	log.Printf("listening on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
