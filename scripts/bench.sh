#!/bin/bash

# Benchmark script for fn application

set -e

echo "Running performance benchmarks for fn..."

# Run different benchmark categories with shorter times to avoid timeouts
echo ""
echo "=== Storage Creation Benchmarks ==="
go test ./internal/storage -bench="BenchmarkNewStore" -benchtime=500ms

echo ""
echo "=== Read Operation Benchmarks ==="
go test ./internal/storage -bench="BenchmarkGetBookmark$|BenchmarkGetBookmarkMiss" -benchtime=500ms

echo ""
echo "=== Light Write Operation Benchmarks ==="
go test ./internal/storage -bench="BenchmarkUpdateUsage" -benchtime=500ms

echo ""
echo "=== Concurrent Operation Benchmarks ==="
go test ./internal/storage -bench="BenchmarkConcurrentRead" -benchtime=500ms

echo ""
echo "=== Mixed Operation Benchmarks ==="
go test ./internal/storage -bench="BenchmarkMixedOperations" -benchtime=500ms

echo ""
echo "Benchmarking complete!"
echo ""
echo "Note: Some heavy I/O benchmarks (SaveBookmark, DeleteBookmark) are excluded"
echo "from this quick run to avoid timeouts. Run them individually if needed:"
echo "  go test ./internal/storage -bench=BenchmarkSaveBookmark -count=1"