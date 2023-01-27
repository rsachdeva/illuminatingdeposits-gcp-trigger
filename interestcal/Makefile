.PHONY: upload-apply
upload-apply: prepare-upload apply

.PHONY: prepare-upload
prepare-upload: clean zip move-to-tf

.PHONY: init
init:
	cd deploy/terraform && terraform init

.PHONY: destroy
destroy:
	cd deploy/terraform && terraform destroy -auto-approve

.PHONY: apply
apply:
	cd deploy/terraform && terraform apply -auto-approve

.PHONY: clean
clean:
	rm -f *.zip
	rm -f deploy/terraform/*.zip

.PHONY: zip
zip:
	zip illuminating-gosource.zip go.mod interestcal.go

.PHONY: copy-for-upload
move-to-tf:
	mv illuminating-gosource.zip deploy/terraform

#only for checking

.PHONY: unzip-check
unzip-check:
	unzip -p deploy/terraform/illuminating-gosource.zip hello.go | less

.PHONY: build-check
build-check:
	go build .
