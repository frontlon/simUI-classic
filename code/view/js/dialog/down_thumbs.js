var ID;
var SID;
var PLATFORM;
var SOURCE;
var PAGE=0;
var PAGENUM = 30;
var SCROLL_POS = 0;
var SCROLL_LOCK = false;
view.root.on("ready", function(){
    Init(); //初始化
});

//初始化
function Init(){
    view.windowCaption = CONF.Lang.ThumbsDown;

    var keyword = view.parameters.keyword;
    ID = view.parameters.id;
    SID = view.parameters.sid;
    PLATFORM = view.parameters.platform;
    SOURCE = $(#source).select("li:nth-child(1)").attributes['opt'];
    //创建语言
    createLang();

    $(#source).select("li:nth-child(1)").state.current = true;

    $(#search_input).value = keyword;
    searchPic(keyword);

    //分页
    self.onScroll = function(evt) {
        pages(evt);
    };
}

//分页
function pages(evt){

    var scrollPos = evt.scrollPos + self.box(#height);
    var boxHeight = $(#down_thumb_list).box(#height);

    if(SCROLL_POS == 0){
        SCROLL_POS = scrollPos;
    }

    if ((boxHeight - scrollPos <=50) && (scrollPos > SCROLL_POS)){

        //如果加锁中，则不执行后续逻辑，防止重复触发
        if (SCROLL_LOCK == true){
            return;
        }
        SCROLL_LOCK = true; //加锁

        var keyword = $(#search_input).value;
        PAGE++;
        var result = searchPic(keyword);
        if(result == true){
            SCROLL_POS = scrollPos;
        }
        SCROLL_LOCK = false;
    }
}


//点击图片
function DownThumb(evt){
    var ctype = view.parameters.type;

    var caption = "";
    var content = "";
    caption = CONF.Lang.Thumb;
    content = CONF.Lang.SetThumbConfirm;

    //下载图片
    var url = evt.select("img").attributes["src"];
    var ext = evt.select("img").attributes["ext"];

    //更改图片
    createFileDropZone(PLATFORM,ctype,SID,url,ext);
}

//搜索图片
function searchThumb(){
    var keyword = $(#search_input).value;

    if(keyword == ""){
        alert(CONF.Lang.InputKeyword);
        return true;
    }

    //清空数据
    $(#down_thumb_list).clear();

    PAGE = 0;

    //搜索内容
    searchPic(keyword);
}

//图片搜索
function searchPic(keyword){

    startLoading();

    //过滤特殊字符
    keyword = keyword.replace(/[\~\!\@\#\$\%\^\&\*\(\)\[\]\{\}\;\'\:\"\'\,\.\/\<\>\?\-\=\_\+]/g,"");

    var response = "{}";
    if(SOURCE == "baidu"){
        response = mainView.SearchThumbsForBaidu(keyword,PAGE);
    }else if (SOURCE == "hfsdb"){
        response = mainView.SearchThumbsForHfsDb(keyword,PAGE);
    }else{
        alert("source not exists.");
    }

    if(response == "{}"){
        alert(CONF.Lang.SearchEmpty);
    }else if(response != undefined){
        var data = JSON.parse(response);
        var content = "";
        for(var typ in data) {
            content +="<strong>"+ typ +"</strong><ul>";
            for(var item in data[typ]) {
                var desc = item.Width + "X" + item.Height + " ." + item.Ext;
                content +="<li><img ext='."+ item.Ext +"' src='"+item.ImgUrl+"'/><p>"+ desc +"</p></li>";
            }
            content += "</ul>"
        }
        if (content != ""){
            $(#down_thumb_list).append(content);
        }
    }
    endLoading();
    return false;
}

//切换来源
function changeSource(evt){
    evt.state.current = true;

    var source = evt.attributes["opt"];

    if (SOURCE == source){
        return;
    }
    SOURCE = source;
    PAGE = 0;
    $(#down_thumb_list).html = "";
    if(source == "baidu"){
        $(.description).html = CONF.Lang["BaiduDescription"];
    }else if (source == "hfsdb"){
        $(.description).html = CONF.Lang["HfsDBDescription"];
    }

    //搜索图片
    searchThumb();
}