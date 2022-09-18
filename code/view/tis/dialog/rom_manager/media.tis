

/**
 * 资源管理
 */
/*
//切换平台 - 资源管理
function managerChangeMediaPlatform(platformId){

    if (platformId == 0){
        return;
    }

    //创建菜单
    var dd = createMenuOption(platformId);
    var menulist = self.$(#menu_media);
    menulist.options.clear();
    if(dd != ""){
        menulist.options.html = dd; //生成dom
    }else{
        menulist.options.html = "<option value=''>"+CONF.Lang.AllGames+"</option>"; //生成dom
    }
    menulist.value = "";

    //生成默认游戏列表
    var req = {
        "platform" : platformId.toInteger(),
    };
    var request = JSON.stringify(req);
    var romjson = mainView.GetGameList(request);
    $(#romlist_media tbody).html = "";
    managerCreateMedia(romjson);

    //生成分页数据
    var romCount = mainView.GetGameCount(request); //必须在获取rom的方法下面
    var pages = managerCreatePages(romCount);
    $(#pages_media).html = pages;
    if (pages != ""){
        $(#pages_media).select("li:nth-child(1)").state.current = true;
    }
    if($$(#pages_media li).length == 1){
        $(#pages_media).style["display"] = "none";
    }
}

//创建资料rom列表 - 资源管理
function managerCreateMedia(romjson){
    var romobj = JSON.parse(romjson);
    var romData = "";
    for(var obj in romobj) {
        if(obj.Menu == "_7b9"){
            obj.Menu = "";
        }

        romData += "<tr rid='"+ obj.Id +"'>";
        var romName = obj.RomPath.replace("\\","/");
        var romPath = romName.split("/");
        var filename = romPath[romPath.length - 1].split(".");

        var thumb = getRomPicPath("thumb",obj.Platform,filename[0]);
        var snap = getRomPicPath("snap",obj.Platform,filename[0]);
        var title = getRomPicPath("title",obj.Platform,filename[0]);
        var poster = getRomPicPath("poster",obj.Platform,filename[0]);
        var packing = getRomPicPath("packing",obj.Platform,filename[0]);
        var cassette = getRomPicPath("cassette",obj.Platform,filename[0]);
        var icon = getRomPicPath("icon",obj.Platform,filename[0]);
        var gif = getRomPicPath("gif",obj.Platform,filename[0]);
        var background = getRomPicPath("background",obj.Platform,filename[0]);
        var wallpaper = getRomPicPath("wallpaper",obj.Platform,filename[0]);
        var video = getRomPicPath("video",obj.Platform,filename[0]);
       // var doc = mainView.GetGameDoc("doc",obj.Id);
        //var strategy = mainView.GetGameDoc("strategy",obj.Id);

        romData += "<td>"+ filename[0] +"</td>";
        romData += "<td>"+ obj.Name.toHtmlString() +"</td>";
        romData += "<td.openfile><input|text.thumb value='"+ thumb +"' /></td>";
        romData += "<td.openfile><input|text.snap value='"+ snap +"' /></td>";
        romData += "<td.openfile><input|text.title value='"+ title +"' /></td>";
        romData += "<td.openfile><input|text.poster value='"+ poster +"' /></td>";
        romData += "<td.openfile><input|text.packing value='"+ packing +"' /></td>";
        romData += "<td.openfile><input|text.cassette value='"+ cassette +"' /></td>";
        romData += "<td.openfile><input|text.icon value='"+ icon +"' /></td>";
        romData += "<td.openfile><input|text.gif value='"+ gif +"' /></td>";
        romData += "<td.openfile><input|text.background value='"+ background +"' /></td>";
        romData += "<td.openfile><input|text.wallpaper value='"+ wallpaper +"' /></td>";
        romData += "<td.openfile><input|text.video value='"+ video +"' /></td>";
        romData += "</tr>"
    }
    $(#romlist_media tbody).html = romData;
}

//切换目录 - 资源管理
function managerChangeMediaMenu(menuName){

    var platformId = $(#platform_media).value;

    //生成默认游戏列表
    var req = {
        "platform" : platformId.toInteger(),
        "catname" : menuName,
    };
    var request = JSON.stringify(req);
    var romjson = mainView.GetGameList(request);
    $(#romlist_media tbody).html = "";
    managerCreateMedia(romjson);

    var romCount = mainView.GetGameCount(request); //必须在获取rom的方法下面
    var pages = managerCreatePages(romCount);
    $(#pages_media).html = pages;
    if(pages != ""){
        $(#pages_media).select("li:nth-child(1)").state.current = true;
    }
    if($$(#pages_media li).length == 1){
        $(#pages_media).style["display"] = "none";
    }
}


//点击分页按钮 - 资源管理
function managerCreateMediaByPages(obj){
    if(obj.state.current == true){
        return;
    }
    var platformId = $(#platform_media).value;
    var id = obj.parent.id;
    var page = obj.html.toInteger() - 1;
    var req = {
        "platform" : platformId.toInteger(),
        "page" : page,
    };
    var request = JSON.stringify(req);
    var romjson = mainView.GetGameList(request);
    $(#romlist_media tbody).html = "";
    managerCreateMedia(romjson);
    obj.state.current = true;
    $(body).scrollTo(0,0, false);
}

//选择文件 - 资源管理
function managerSetMedia(obj) {
    var classes = obj.attributes["class"].split(" ");
    var id = obj.parent.parent.attributes["rid"];
    var uri = obj.value;
    var opt = classes[0];
    obj.attributes.removeClass("nosave");
    if(uri == ""){
        //mainView.DeleteThumbs(opt,id);
        alert(CONF.Lang.DeletePicSuccess);
    }else{
        mainView.EditRomThumbs(opt,id,uri);
    }
    
}
*/