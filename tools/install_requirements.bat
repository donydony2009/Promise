@echo off
for /f "tokens=*" %%a in (requirements.txt) do (
  echo Getting %%a
  go get -u %%a
)
pause