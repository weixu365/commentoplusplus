docker_instance_id=$(docker ps | grep postgres | cut -d" " -f1)

docker exec -it -u postgres $docker_instance_id psql -c "drop database if exists commento_test;"
docker exec -it -u postgres $docker_instance_id psql -c "create database commento_test;"
docker exec -it -u postgres $docker_instance_id psql -c "grant all privileges on database commento_test to postgres;"


