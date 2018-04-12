FROM golang:latest 
RUN mkdir /authentication 
ADD . /authentication/ 
WORKDIR /authentication 
RUN go get -u github.com/satori/go.uuid
RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/go-sql-driver/mysql
RUN go get -u github.com/donydony2009/Promise/scripts/mysql
RUN go build -o authentication ./service
RUN go build -o main .
CMD ["/authentication/main"]