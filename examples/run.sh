#!/bin/sh

echo "\n--- cat some files --------------------------------------------\n"
(cd cat; echo -e "/etc/issue\n/unknown/file" | xargs go run cat.go)

echo "\n--- run a short timer -----------------------------------------\n"
(cd timer; go run timer.go 3 1)

echo "\n--- count words of some files ---------------------------------\n"
(cd wc; go run wc.go /etc/group /unknown/file)
