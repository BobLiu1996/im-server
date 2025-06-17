# start redis
docker run -d --name redis -v /home/mystic/redis_data:/data -p 6379:6379 redis
# start mysql
docker run -d --name mysql -e MYSQL_ROOT_PASSWORD=mystic -p 3306:3306 -v /home/mystic/mysql_data:/var/lib/mysql mysql:latest