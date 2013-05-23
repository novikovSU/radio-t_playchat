#!/bin/sh

PATH="$PATH:/home/novikov/apps/s3cmd-1.5.0"

s3cmd put --acl-public index.html s3://playchat.radio-t.com/index.html
#s3cmd put --acl-public assets/images/favicon.png s3://playchat.radio-t.com/assets/images/favicon.png
#s3cmd put --acl-public --recursive assets s3://playchat.radio-t.com/
#s3cmd put --acl-public assets/candy/default.css s3://playchat.radio-t.com/assets/candy/
