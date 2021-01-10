#!/usr/bin/env bash

RED="\e[31m"
GREEN="\e[32m"
YELLOW="\e[33m"
RESET="\e[0m"
CHECK="✓"
CROSS="✗"
WARN="⚠"

describe() {
    printf "%s" "$1"
    dots=${2:-3}
    for _ in $(seq 1 "${dots}"); do sleep 0.035; printf "."; done
    sleep 0.035
}

log_warn() {
    message=${1:-"Warning"}
    log="${YELLOW}${WARN} ${message}${RESET}"
    echo -e "${log}"
    [ -f "$3" ] && echo "$2 ${log}" >> "$3"
}

log_failed() {
    message=${1:-"Failed"}
    log="${RED}${CROSS} ${message}${RESET}"
    echo -e "${log}"
    [ -f "$3" ] && echo "$2 ${log}" >> "$3"
}

log_success() {
    message=${1:-"Success"}
    log="${GREEN}${CHECK} ${message}${RESET}"
    echo -e "${log}"
    [ -f "$3" ] && echo "$2 ${log}" >> "$3"
}
