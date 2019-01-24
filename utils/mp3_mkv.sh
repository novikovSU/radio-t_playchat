#!/bin/bash

ffmpeg -loop 1 -framerate 1 -i rt_podcast631_cover.jpg -i rt_podcast631.mp3 -c:v libx264 -crf 0 -preset veryfast -tune stillimage -c:a copy -shortest rt_podcast631.mkv