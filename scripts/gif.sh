#!/bin/bash

# mov to gif
ffmpeg -i input.mov -pix_fmt rgb8 -r 10 output.gif && gifsicle -O3 output.gif -o output.gif

# speed up gif
# 1x30 is ImageMagick’s notation for 1 times 1/30th of a second
convert -delay 1x100 output.gif flemq.gif

exit 0
