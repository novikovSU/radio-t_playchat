#!/bin/sh

PATH="$PATH:/home/novikov/apps/s3cmd-1.1.0-beta3"

s3cmd put --acl-public index.html s3://playchat.radio-t.com/index.html
s3cmd put --acl-public assets/images/favicon.png s3://playchat.radio-t.com/assets/images/favicon.png
s3cmd put --acl-public --recursive assets/jplayer s3://playchat.radio-t.com/assets/
