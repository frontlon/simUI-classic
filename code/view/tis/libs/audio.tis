//播放单个音频
function playAudio(url){
    var data = [];
    if(url != ""){
        data[0] = url;
    }else{
        var audio = $$(#audio li);
        var i = 0;
        for(var a in audio) {
            data[i] = a.attributes["path"];
            i++;
        }
    }    
    view.PlayAudio(JSON.stringify(data));
}