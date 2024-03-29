# pm2systray

This app is a systray tool that let you perform simple controls on [pm2](https://pm2.keymetrics.io/) (process manager) from a systray applet.

It helps switch **quickly and easily** between projects.

I use pm2 to manage my day to day daemon dev apps, this saves me time to startup a coding day and speedup my task switching along the day.

I also wanted a low footprint on this app and make it portable.

Whatsmore, I wanted to code some go things and create some systray for something that can be **quite usefull**. this is done.



https://user-images.githubusercontent.com/5399780/120869083-db00e500-c595-11eb-9701-8319ee1a4d54.mp4

## Requirements

`sudo apt-get install build-essential libgtk-3-dev libappindicator3-dev gir1.2-appindicator3-0.1 libayatana-appindicator3-dev`

## Features

- **0 configuration**, start the app, it manage pm2 tasks from current setup
- lightweight and **one binary app**
- one click control pm2 processes `start | stop`
- one click `save` pm2 configuration
- work on `linux`, `windows` and should also work on `macos` (not tested yet)

## Notes

- If you like it and want to use it often, **run it at startup**
- If you add one process to the pm2 list, **just restart the app**
