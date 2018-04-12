call get_service_info %*
call build_docker %*
SET PORT=%SERVICE_ID%000
docker rm -f promise-%SERVICE%
docker container run --name promise-%SERVICE% --net promise-net --ip 172.18.0.%SERVICE_ID% --link promise-mysql -p %PORT%:%PORT% %SERVICE%:1.0