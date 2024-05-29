include .env

test:
	ASSETS_DIR=$$(pwd)/assets go test -tags=test ./...

clean:
	rm service/*.png

deploy:
	gcloud functions deploy short-url \
	--gen2 \
	--runtime go122 \
	--trigger-http \
	--allow-unauthenticated \
	--entry-point HandleRequest \
	--set-env-vars ASSETS_DIR=serverless_function_source_code/assets \
	--set-env-vars SHORT_PREFIX=${SHORT_PREFIX} \
	--set-env-vars GCP_PROJECT_ID=${GCP_PROJECT_ID} \
	--project ${GCP_PROJECT_ID} \
	--region ${GCP_REGION}
