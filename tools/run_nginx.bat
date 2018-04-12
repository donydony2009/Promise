
FOR /F "tokens=* USEBACKQ" %%F IN (`docker run --name nginx-reverse-proxy -v nginx.conf:/etc/nginx/nginx.conf:ro -p 80:80 -d nginx`) DO (
SET container=%%F
)
echo %container%