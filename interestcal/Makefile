prepare-upload: clean zip move-to-tf

.PHONY: clean
clean:
	rm -f *.zip
	rm -f ../deploy/terraform/*.zip

.PHONY: zip
zip:
	zip illuminating-gosource.zip go.mod hello.go

.PHONY: copy-for-upload
move-to-tf:
	mv illuminating-gosource.zip ../deploy/terraform

#only for checking zip contents
.PHONY: unzip-check
unzip-check:
	unzip -p ../deploy/terraform/illuminating-gosource.zip hello.go | less
