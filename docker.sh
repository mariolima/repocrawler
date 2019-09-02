docker image rm -f repocrawl
docker build . -t repocrawl
docker run -it -e 'GITHUB_ACCESS_TOKEN=TOKEN' repocrawl -h
