docker build -t todoserver .
docker run -p 2323:8888 -it --name todoserver todoserver
http://localhost:2323
docker start -i todoserver
docker system prune -af