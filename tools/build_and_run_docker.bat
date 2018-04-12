call build_docker %*
SET PORT=%2
SET IN_PORT=%3
docker container run --name promise-%SERVICE% --link promise-mysql -p %PORT%:%IN_PORT% %SERVICE%:1.0