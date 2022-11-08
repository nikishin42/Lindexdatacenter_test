build:
	docker build -f Dockerfile.Multistage -t lindex_test .

run_json:
	docker run --rm lindex_test db.json

run_csv:
	docker run --rm lindex_test db.csv