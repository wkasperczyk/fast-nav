#!/bin/bash

# Test coverage script for fn application

set -e

echo "Running tests with coverage..."

# Run all tests with coverage
go test ./... -coverprofile=coverage.out -covermode=atomic

echo ""
echo "Coverage summary:"
go tool cover -func=coverage.out

echo ""
echo "Generating HTML coverage report..."
go tool cover -html=coverage.out -o coverage.html

echo ""
echo "Coverage report saved to coverage.html"
echo "Open coverage.html in a browser to view detailed coverage"

# Optional: Open the report automatically (uncomment if desired)
# if command -v open >/dev/null 2>&1; then
#     open coverage.html
# elif command -v xdg-open >/dev/null 2>&1; then
#     xdg-open coverage.html
# fi