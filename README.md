## pion-radio-example

> Audio (G.711, PCMU) Streaming Example


### Usage
1. start the server
    ```
    $ make rs (run-server)
    ```

2. shoot audio data (sample.wav) to udp port (1234)
    ```
    $ cd audio
    $ make fs (ffmpeg-streaming)
    ```

3. open the web page
    ```
    open http://localhost:8080/static/rtp-radio.html
    press "start session"
    ```

### Reference
- [scottstensland/websockets-streaming-audio](https://github.com/scottstensland/websockets-streaming-audio)


