prepare-for-upload: zip copy-for-upload

.PHONY: zip
zip:
	zip goilluminating.zip go.mod hello.go

.PHONY: copy-for-upload
copy-for-upload:
	cp goilluminating.zip deploy/terraform

.PHONY: build
build:
	go build .
