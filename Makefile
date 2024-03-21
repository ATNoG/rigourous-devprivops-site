.Phony: templ
templ:
	templ generate

.Phony: css
css:
	tailwindcss -i static/css/style.css -o static/css/style.css --minify

.Phony: clean
clean:
	rm static/css/style.css  
	rm templates/*_templ.go
