.PHONY: help
help:
	@echo

.PHONY: clean
clean:
	rm -rf *.pem *.fpx

.PHONY: fpx
fpx:
	openssl req -x509 \
		-sha256 -days 365 \
		-nodes \
		-newkey rsa:4096 \
		-subj "/C=US/ST=CA/O=Azure/CN=myapp" \
		-keyout private.pem -out certificate.pem

	openssl pkcs12 -export -in certificate.pem -inkey private.pem -out bundle.fpx

.env:
	cp sample.env .env