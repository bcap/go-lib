test:
	go list -m -f '{{.Dir}}' | xargs -I {} go test -v {}