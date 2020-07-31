# pion-radio-example

> Audio (G.711, PCMU) Streaming Example


### Usage
1. start the server
    ```
    $ make rs`
    ```

2. shoot audio data to udp port
    ```
    $ cd audio
    $ make fs
    ```

3. open the web page
    ```
    open http://localhost:8080/static/rtp-radio.html
    ```


