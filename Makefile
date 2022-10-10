
docker.build:
	docker build -t backup:latest .

docker.run:
	docker run --rm --name backup \
		-v ~/Documents:/root/Documents \
		-v ~/Documents/fleet/backup/cfg.yaml:/etc/backup/config.yaml \
		-v ~/.aws:/root/.aws backup:latest