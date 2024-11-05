test:
	go list -m -f '{{.Dir}}' | xargs -I {} go test -v {}

vet:
	go list -m -f '{{.Dir}}' | xargs -I {} go vet {}