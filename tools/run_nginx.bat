docker rm -f nginx-reverse-proxy
docker run --name nginx-reverse-proxy --link promise-promise --link promise-authentication -p 80:80 -d nginx
docker cp nginx.conf nginx-reverse-proxy:/etc/nginx/nginx.conf 
docker stop nginx-reverse-proxy
docker start nginx-reverse-proxy -a