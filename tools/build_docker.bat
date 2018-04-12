call get_service_info %*
docker build %SCRIPTS_PATH%\%SERVICE% --no-cache -t %SERVICE%:1.0