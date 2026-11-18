
GOKARGS = -i gokrazy --parent_dir .

assets:
	cd frontend && \
	npm install && \
	npm run build

binary:
	GOOS=linux GOARCH=arm64 go build -o "gokrazy/inshrimp" .


update:
	gok $(GOKARGS) update


.PHONY: update