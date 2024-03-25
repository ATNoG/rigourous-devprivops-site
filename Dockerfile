# Use image with golang
FROM golang

COPY . /src/
WORKDIR /src/ 

# Build the application
RUN ls
RUN go mod tidy
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go build

# Cleanup
RUN mv devprivops-dashboard /bin/devprivops-dashboard
RUN rm -rf /src/
