# Use image with golang
FROM golang

COPY . /src/
WORKDIR /src/ 

# Build the application
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
RUN tailwindcss -i static/css/source.css -o static/css/style.css --minify
RUN go mod tidy
RUN go build

# Cleanup
RUN mv devprivops-dashboard /bin/devprivops-dashboard
RUN rm -rf /src/
