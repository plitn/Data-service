FROM mirror.gcr.io/golang:latest AS builder

WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/main ./cmd/data-service/main.go

FROM mirror.gcr.io/python:3.12.3

WORKDIR /app

COPY requirements.txt .
RUN pip install --upgrade pip
RUN pip install -r requirements.txt

COPY --from=builder /go/src/app/build/main /app/build/main
COPY script.py /app/build/

RUN ls -l /app/build/

CMD ["/app/build/main"]
