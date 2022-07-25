build: 
	env GOOS=linux GOARCH=amd64 go build -o pm2systray main.go
	chmod +x pm2systray

dev:
	fresh
