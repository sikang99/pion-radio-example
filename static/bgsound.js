function playWavFileInIE(src){
    if(/msie/i.test(navigator.userAgent) || /trident/i.test(navigator.userAgent)){
        var origPlayer = document.getElementById('currentWavPlayer');
        if(origPlayer){
            origPlayer.src = '';
            origPlayer.outerHtml = '';
            document.body.removeChild(origPlayer);
            delete origPlayer;
        }
        var newPlayer = document.createElement('bgsound');
        newPlayer.setAttribute('id', 'currentWavPlayer');
        newPlayer.setAttribute('src', src);
        document.body.appendChild(newPlayer);
        return false;
    }
}
