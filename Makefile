assets:
	cd frontend && \
	npm install && \
	npm run build

# binary:
# 	GOOS=linux GOARCH=arm64 go build -o "gokrazy/inshrimp" .

update:
	./bin/gok.sh update


.PHONY: update assets
