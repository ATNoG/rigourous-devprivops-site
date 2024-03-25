# Use image with golang
FROM golang

COPY . /src/
WORKDIR /src/ 

# Build the application
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate

RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
RUN chmod +x tailwindcss-linux-x64
RUN mv tailwindcss-linux-x64 tailwindcss
RUN ./tailwindcss -i static/css/source.css -o static/css/style.css --minify

RUN go mod tidy
RUN go build

# Cleanup
RUN mv devprivops-dashboard /bin/devprivops-dashboard
RUN rm -rf /src/
