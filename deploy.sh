#!/bin/sh

PATH="$PATH:/home/novikov/apps/s3cmd-1.1.0-beta3"

s3cmd put --acl-public index.html s3://dev.playchat.radio-t.com/index.html
s3cmd put --acl-public assets/audiojs/audio.js s3://dev.playchat.radio-t.com/assets/audiojs/audio.js
s3cmd put --acl-public assets/audiojs/audiojs.swf s3://dev.playchat.radio-t.com/assets/audiojs/audiojs.swf
s3cmd put --acl-public assets/audiojs/player-graphics.gif s3://dev.playchat.radio-t.com/assets/audiojs/player-graphics.gif
s3cmd put --acl-public assets/images/favicon.png s3://dev.playchat.radio-t.com/assets/images/favicon.png
