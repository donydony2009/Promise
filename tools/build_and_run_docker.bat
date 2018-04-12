call build_docker %*
SET PORT=%2
docker container run --name promise-%SERVICE% --link promise-mysql -p %PORT%:8080 %SERVICE%:1.0