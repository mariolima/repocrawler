docker image rm -f repocrawler
docker build . -t repocrawler
docker run -it -e 'GITHUB_ACCESS_TOKEN=TOKEN' repocrawler -h
