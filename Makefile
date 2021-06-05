build: 
	go build -o pm2systray *.go
	chmod +x pm2systray

dev:
	fresh
