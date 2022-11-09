build:
	docker build -f Dockerfile.Multistage -t linxdc_test .

run_json:
	docker run --rm linxdc_test db.json

run_csv:
	docker run --rm linxdc_test db.csv