#!/bin/bash

# Beeflang development helper script

case "$1" in
  run)
    if [ -z "$2" ]; then
      echo "Usage: ./dev.sh run <file.beef>"
      exit 1
    fi
    go run main.go "$2"
    ;;
  test)
    go test ./... -v
    ;;
  test-pkg)
    if [ -z "$2" ]; then
      echo "Usage: ./dev.sh test-pkg <package>"
      echo "Example: ./dev.sh test-pkg lexer"
      exit 1
    fi
    go test "./internal/$2" -v
    ;;
  build)
    go build -o beeflang main.go
    echo "Built: ./beeflang"
    ;;
  clean)
    rm -f beeflang
    echo "Cleaned build artifacts"
    ;;
  *)
    echo "Beeflang Dev Commands:"
    echo "  ./dev.sh run <file.beef>  - Run a Beeflang program"
    echo "  ./dev.sh test             - Run all tests"
    echo "  ./dev.sh test-pkg <pkg>   - Run tests for specific package"
    echo "  ./dev.sh build            - Build beeflang binary"
    echo "  ./dev.sh clean            - Remove build artifacts"
    ;;
esac
