FROM golang:1.16

WORKDIR /app

COPY . /app

RUN cd server; go build -o bin_app .
RUN mv server/bin_app bin_app

CMD ./bin_app

