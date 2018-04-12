FROM golang:latest 
RUN mkdir /promise 
ADD . /promise/ 
WORKDIR /promise 
RUN go get -u github.com/satori/go.uuid
RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/go-sql-driver/mysql
RUN go get -u github.com/donydony2009/Promise/scripts/mysql
RUN go get -u github.com/donydony2009/Promise/scripts/authentication
RUN go get -u github.com/donydony2009/Promise/scripts/rest
RUN go build -o promise ./service
RUN go build -o main .
CMD ["/promise/main"]