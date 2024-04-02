build:
	templ generate && go build -ldflags="-s -w"
