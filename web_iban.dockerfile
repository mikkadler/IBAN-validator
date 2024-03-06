FROM golang:alpine as builder

# Create a new directory for the app and copy the code into it
RUN mkdir /app
ADD . /app
WORKDIR /app

# Build the Go application
RUN go build -o main .

# Use a minimal base image for the final container
FROM alpine

# Copy the built executable from the builder stage
COPY --from=builder /app/main /app/main
COPY --from=builder /app/iban/data /app/iban/data

# Expose port 8080 for the application
EXPOSE 8080

# Run the main binary when the container starts
CMD ["/app/main"]
