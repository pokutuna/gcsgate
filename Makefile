.PHONY: build run deploy

build:
	go build -o gcsgate .

run:
	go run .

deploy:
ifndef PROJECT
	$(error PROJECT is not set. Usage: make deploy PROJECT=your-gcp-project SERVICE_ACCOUNT=sa@project.iam.gserviceaccount.com)
endif
ifndef SERVICE_ACCOUNT
	$(error SERVICE_ACCOUNT is not set. Usage: make deploy PROJECT=your-gcp-project SERVICE_ACCOUNT=sa@project.iam.gserviceaccount.com)
endif
	gcloud app deploy --project=$(PROJECT) --service-account=$(SERVICE_ACCOUNT)
