#!/bin/bash

set -e

OUT_DIR="./bin"

if [[ "${1:-}" == "clean" ]]; then
    rm -rf "$OUT_DIR"
    exit 0
fi

case "${1:-build}" in
    "build")
        mkdir -p "$OUT_DIR"
        CMD="go build -o $OUT_DIR"
        ;;
    "run")
        CMD="go run"
        ;;
    # "test")
    #     CMD="go test -v"
    #     ;;
    *)
        echo "Invalid argument: $1"; exit 1 ;;
esac

case "${2:-all}" in
    "all")
        FILES="./..."
        ;;
    "gui")
        FILES="./cmd/distribute-gui/gui.go"
        ;;
    "cli")
        FILES="./cmd/distribute-cli/cli.go"
        ;;
    *)
        echo "Invalid argument: $SECOND_ARG"; exit 1 ;;
esac

if [[ "$CMD" == "go run" && "$FILES" == "./..." ]]; then
    echo "Please specify if the gui or the cli should run!"
    exit 1
fi

# shellcheck disable=SC2086
$CMD "$FILES"
