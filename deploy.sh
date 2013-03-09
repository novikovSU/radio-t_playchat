#!/bin/sh

PATH="$PATH:/home/novikov/apps/s3cmd-1.1.0-beta3"

s3cmd put --acl-public index.html s3://dev.playchat.radio-t.com/index.html
#s3cmd put --acl-public assets/images/favicon.png s3://dev.playchat.radio-t.com/assets/images/favicon.png
#s3cmd put --acl-public assets/jplayer/jquery.jplayer.min.js s3://dev.playchat.radio-t.com/assets/jplayer/jquery.jplayer.min.js
#s3cmd put --acl-public assets/jplayer/Jplayer.swf s3://dev.playchat.radio-t.com/assets/jplayer/Jplayer.swf
#s3cmd put --acl-public --recursive assets/jplayer/skins s3://dev.playchat.radio-t.com/assets/jplayer/
