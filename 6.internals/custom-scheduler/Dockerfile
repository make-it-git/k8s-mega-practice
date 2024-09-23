FROM golang:1.22 as builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o simple-scheduler .

FROM scratch
COPY --from=builder /app/simple-scheduler /simple-scheduler
ENTRYPOINT ["/simple-scheduler"]

