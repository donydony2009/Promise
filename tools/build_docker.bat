SET SCRIPTS_PATH=..\scripts
SET SERVICE=%1
docker build %SCRIPTS_PATH%\%SERVICE% --no-cache -t %SERVICE%:1.0