#
# Makefile
#
HOST=localhost
PORT=1234
all:
	@echo "usage: make [vlc|ffmpeg]"

## vlc
vlc:
	@echo "usage: make (vlc) [streaming|play]"

vlc-streaming vs:
	#vlc sample.wav --sout "#transcode{acodec=ulaw,channels=1,ar=8000}:rtp{dst=$(HOST),port-audio=$(PORT)}"
	vlc ThePolice.wav --sout "#transcode{acodec=ulaw,channels=1,ar=8000}:rtp{dst=$(HOST),port-audio=$(PORT)}"

vlc-play vp:
	vlc rtp://@:$(PORT)

## ffmpeg
ffmpeg:
	@echo "usage: make (ffmpeg) [streaming|play]"

ffmpeg-streaming fs:
	ffmpeg -re -stream_loop -1 -i sample.wav -ar 8k -ac 1 -acodec pcm_mulaw -f rtp rtp://$(HOST):$(PORT)

ffmpeg-play fp:
	ffplay rtp://$(HOST):$(PORT)

