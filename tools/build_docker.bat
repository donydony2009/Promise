
SET SERVICE=%1
rmdir /S /Q docker
mkdir docker
xcopy src\%SERVICE% docker /E
xcopy src\mysql docker /E
xcopy src\rest docker /E