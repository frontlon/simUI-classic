view.root.on("ready", function(){

    //模拟器拖拽排序
    DragDrop{
        what      : "#sim_select>li:not(.fixed)",
        where     : "#sim_select",
        container : "#sim_select",
        notBefore : "#sim_select>li.fixed",
        dropped : function(){
            var lis = $$(#sim_select > li);
            var id = "";
            var newObj = {};
            for(var li in lis) {
                id = li.attributes["sim"];
                newObj[id] = li.index;
            }
            //更新到数据库
            var datastr = JSON.stringify(newObj);
            var result = view.UpdateSimSort(datastr);
        }
    }

});

//单击游戏侧边栏详情
function openSidebar(obj){

    if(obj.state.current == true){
        return;
    }

    obj.state.current = true;
    
    var getjson = view.GetGameDetail(obj.attributes["rid"]);
    var detailObj = JSON.parse(getjson);
    var info = detailObj.Info;
    var romPath = URL.parse(info.RomPath);

    //列表背景图片
    var ext = "";
    var background = getRomPicPath("background",info.Platform,romPath.name); //背景图
    showListBackground(background);

    //侧边栏隐藏的情况下，不执行后续操作
    if ($(#right).style["display"] == "none"){
        return;
    }

    //隐藏侧边栏遮罩
    $(#mask).style["display"] = "none";
    $(#platform_desc).style["display"] = "none";
    

   //侧边栏背景图片
    var wallpaper = getRomPicPath("wallpaper",info.Platform,romPath.name); //背景图
    if(wallpaper != ""){
        $(#right_bg).style["background-image"] = [url:URL.fromPath(wallpaper)];
    }else{
        if(CONF.Default.WallpaperImage != ""){
            $(#right_bg).style["background-image"] = [url:URL.fromPath(CONF.Default.WallpaperImage)];
        }else{
            $(#right_bg).style["background-image"] = "none";
        }

    }

    //侧边栏指定rom id
    ACTIVE_ROM_ID = info.Id;
    
    //创建滑动图
    _createRotate(info);
    
    //生成模拟器列表
    $(#sim_select).html = ""; //先清空所有模拟器
    var pfObj = CONF.Platform[info.Platform.toString()];

    var simlist = detailObj.Simlist;
    
    if(simlist.length > 0 ){
        var simliststr = "";
        for(var obj in simlist) {
            simliststr += "<li filename=\""+ URL.toPath(obj.Path) +"\" rom='"+ info.Id +"' sim='"+ obj.Id +"'><p>"+ obj.Name +"</p><button>"+ CONF.Lang.Edit +"</button></li>";
        }
        $(#sim_select).html = simliststr;

    }

    //如果该平台有模拟器
    if($(#sim_select li) != undefined){
        //先找rom对应的模拟器
        var SimDom = $(#sim_select).select("li[sim="+info.SimId+"]");

        //找不到rom对应的模拟器，找平台默认模拟器
        if(info.SimId == 0 || SimDom == undefined){
            SimDom = $(#sim_select).select("li[sim="+ pfObj.UseSim["Id"] +"]");
        }

        //如果平台模拟器也不存在，则读取第一个模拟器
        if(SimDom == undefined){
            SimDom = $(#sim_select li)
        }

        //选定模拟器
        SimDom.state.current = true;
    }

    //生成运行游戏按钮
    var btndata = "<li class='run' rid='"+ info.Id +"'>"+info.Name.toHtmlString()+"</li>";
    for(var sub in detailObj.Sublist) {
        btndata += "<li class='run' rid='"+ sub.Id +"'>"+sub.Name.toHtmlString() + "</li>";
    }
    self.$(#buttons).html = btndata;

    //显示游戏基础信息
    var rombase = "";
    var exts = ["bat","cmd","lnk","path","slnk"];
    var fileExt = getFileExt(info.RomPath);
    if(info.BaseType != "")      {rombase += "<li readonly title=\""+ CONF.Lang.BaseType +"\">" + info.BaseType.toHtmlString() + "</li>";}
    if(info.BaseYear != "")      {rombase += "<li readonly title=\""+ CONF.Lang.BaseYear +"\">" + info.BaseYear.toHtmlString() + "</li>";}
    if(info.BaseProducer != "")  {rombase += "<li readonly title=\""+ CONF.Lang.BaseProducer +"\">" + info.BaseProducer.toHtmlString() + "</li>";}
    if(info.BasePublisher != "") {rombase += "<li readonly title=\""+ CONF.Lang.BasePublisher +"\">" + info.BasePublisher.toHtmlString() + "</li>";}
    if(info.BaseCountry != "")   {rombase += "<li readonly title=\""+ CONF.Lang.BaseCountry +"\">" + info.BaseCountry.toHtmlString() + "</li>";}
    if(info.BaseTranslate != "") {rombase += "<li readonly title=\""+ CONF.Lang.BaseTranslate +"\">" + info.BaseTranslate.toHtmlString() + "</li>";}
    if(info.BaseVersion != "")   {rombase += "<li readonly title=\""+ CONF.Lang.BaseVersion +"\">" + info.BaseVersion.toHtmlString() + "</li>";}
    if(info.BaseNameEn != "")    {rombase += "<li readonly title=\""+ CONF.Lang.BaseNameEn +"\">" + info.BaseNameEn.toHtmlString()      + "</li>";}
    if(info.BaseNameJp != "")    {rombase += "<li readonly title=\""+ CONF.Lang.BaseNameJp +"\">" + info.BaseNameJp.toHtmlString()      + "</li>";}
    if(info.BaseOtherA != "")    {rombase += "<li readonly >" + info.BaseOtherA.toHtmlString() + "</li>";}
    if(info.BaseOtherB != "")    {rombase += "<li readonly >" + info.BaseOtherB.toHtmlString() + "</li>";}
    if(info.BaseOtherC != "")    {rombase += "<li readonly >" + info.BaseOtherC.toHtmlString() + "</li>";}
    if(info.BaseOtherD != "")    {rombase += "<li readonly >" + info.BaseOtherD.toHtmlString() + "</li>";}
    if(info.Size != "" && exts.indexOf(fileExt) == -1) {rombase += "<li readonly title=\""+ CONF.Lang.FileSize +"\">" + info.Size.toHtmlString() + "</li>";}
    
    rombase += "<li readonly title=\""+ CONF.Lang.Platform +"\">" + CONF.Platform[info.Platform.toString()].Name + "</li>";
    rombase += "<li readonly title=\""+ CONF.Lang.FileName +"\">" + info.RomPath.toHtmlString() + "</li>";

    $(#right_rombase).html = rombase;


   //星级
    createScore(info.Score);

    //游戏运行信息
    $(#run_num).html = info.RunNum.toString();
    let lasttime = timestampToYmd(info.RunLasttime);
    $(#run_lasttime).html = lasttime == "" ? "-" : lasttime;

    //显示简介
    if(detailObj.DocContent != ""){
        $(#doc).attributes.removeClass("nodoc");
        $(#doc).style["display"] = "block";
        //$(#doc).html = markdown.parse(detailObj.DocContent);
        $(#doc).html = detailObj.DocContent;
        $(#doc_title).style["display"] = "block";
    }else{
        if($(#right).attributes.hasClass("widthMode")){
            $(#doc_title).style["display"] = "block";
            $(#doc).html = CONF.Lang.DocEmpty;
            $(#doc).attributes.addClass("nodoc");
        }else{
            $(#doc_title).style["display"] = "none";
            $(#doc).style["display"] = "none";
        }
    }

    //通关状态
    createGameComplete(info.Complete);

    //读取相关游戏
    var related = view.GetRelatedGames(obj.attributes["rid"]);
    var relatedObj = JSON.parse(related);
    var relatedDom = "";
    for(var obj in relatedObj) {

        var romName = getRomName(obj.RomPath);
        var thumb = getRomPicPath("",obj.Platform,romName);

        //没有找到图片，读取默认缩略图
        if(thumb == "" && CONF.Theme[CONF.Default.Theme].Params["default-thumb-image"] != undefined){
            thumb = URL.fromPath(CONF.RootPath + "theme/" + CONF.Theme[CONF.Default.Theme].Params["default-thumb-image"]);
        }

        relatedDom += "<li rid='"+obj.Id+"'>";
        relatedDom += "<img src='"+ thumb +"'/>";
        relatedDom += "<p>"+ obj.Name.toHtmlString() +"</p></li>";
    }
    $(#related).html = relatedDom;

    //选项卡2图集
    _createSecondThumbs(info);

    //攻略文件
    var strategyFiles = detailObj.StrategyFiles;
    var files = "";
        for(var f in strategyFiles) {
            var fileExt = getFileExt(f.path);
            var ico = "this://app/images/filetype/"+ fileExt +".png";
            files += "<li path=\""+ f.path +"\" style='background-image:url("+ ico +")'>"+ f.name +"</li>";
        }
    $(#third_files).html = files;

    //攻略
    var strategy = view.GetGameDoc("strategy",obj.attributes["rid"],1);
    if(strategy == "" && files == ""){
        $(#idx_label_third).attributes.addClass("disable");
    }else{
        $(#idx_label_third).attributes.removeClass("disable");
    }

    if(strategy != ""){
        strategy = strategy.replace("_+","<img class='ctl' src='this://app/images/ctl_+.png' />");
        strategy = strategy.replace("_1","<img class='ctl' src='this://app/images/ctl_1.png' />");
        strategy = strategy.replace("_2","<img class='ctl' src='this://app/images/ctl_2.png' />");
        strategy = strategy.replace("_3","<img class='ctl' src='this://app/images/ctl_3.png' />");
        strategy = strategy.replace("_4","<img class='ctl' src='this://app/images/ctl_4.png' />");
        strategy = strategy.replace("_5","<img class='ctl' src='this://app/images/ctl_5.png' />");
        strategy = strategy.replace("_6","<img class='ctl' src='this://app/images/ctl_6.png' />");
        strategy = strategy.replace("_7","<img class='ctl' src='this://app/images/ctl_7.png' />");
        strategy = strategy.replace("_8","<img class='ctl' src='this://app/images/ctl_8.png' />");
        strategy = strategy.replace("_9","<img class='ctl' src='this://app/images/ctl_9.png' />");
        strategy = strategy.replace("_A","<img class='ctl' src='this://app/images/ctl_A.png' />");
        strategy = strategy.replace("_B","<img class='ctl' src='this://app/images/ctl_B.png' />");
        strategy = strategy.replace("_C","<img class='ctl' src='this://app/images/ctl_C.png' />");
        strategy = strategy.replace("_D","<img class='ctl' src='this://app/images/ctl_D.png' />");
        strategy = strategy.replace("_E","<img class='ctl' src='this://app/images/ctl_E.png' />");
        strategy = strategy.replace("_F","<img class='ctl' src='this://app/images/ctl_F.png' />");
        strategy = strategy.replace("_N","<img class='ctl' src='this://app/images/ctl_N.png' />");
        strategy = strategy.replace("_S","<img class='ctl' src='this://app/images/ctl_S.png' />");
        strategy = strategy.replace("_P","<img class='ctl' src='this://app/images/ctl_P.png' />");
        strategy = strategy.replace("_K","<img class='ctl' src='this://app/images/ctl_K.png' />");
        /*
        strategy = strategy.replace("_a","<img class='ctl' src='this://app/images/ctl_A.png' />");
        strategy = strategy.replace("_b","<img class='ctl' src='this://app/images/ctl_B.png' />");
        strategy = strategy.replace("_c","<img class='ctl' src='this://app/images/ctl_C.png' />");
        strategy = strategy.replace("_d","<img class='ctl' src='this://app/images/ctl_D.png' />");
        strategy = strategy.replace("_e","<img class='ctl' src='this://app/images/ctl_E.png' />");
        strategy = strategy.replace("_f","<img class='ctl' src='this://app/images/ctl_F.png' />");
        strategy = strategy.replace("_n","<img class='ctl' src='this://app/images/ctl_N.png' />");
        strategy = strategy.replace("_s","<img class='ctl' src='this://app/images/ctl_S.png' />");
        strategy = strategy.replace("_p","<img class='ctl' src='this://app/images/ctl_P.png' />");
        strategy = strategy.replace("_k","<img class='ctl' src='this://app/images/ctl_K.png' />");
        */
        var background = getRomPicPath("background",info.Platform,romPath.name); //背景图
        $(#third_strategy).html = strategy;
    }else{
        $(#third_strategy).html = "<div class='sidebar_empty'>"+ CONF.Lang.StrategyEmpty +"</div>";
    }

    //音频文件
    var audioList = detailObj.AudioList;
    var audio = "";
    for(var f in audioList) {
        audio += "<li path=\""+ f.path +"\">"+ f.name +"</li>";
    }
    $(#audio).html = audio;
    
    //是否显示音频标题
    if(audio == ""){
        $(#audio_title_wrapper).style["display"] = "none";
    }else{
        $(#audio_title_wrapper).style["display"] = "block";
    }
}

//视频
function _setVideoControl(movies){
    var movie = "";
    if(movies.length > 0){
        movie = movies[0];
    }

    if(movie == ""){
        return;
    }

    if(VIDEO_ALLOW_PLAY == 0){
        return;
    }

    var volume = VIDEO_VALUME == 0 ? "<i.sound_on></i>" :"<i.sound_off></i>";
    var play = "<i.video_pause></i>";
    var pics = "<div#right_video_wrapper>";
    pics += "<video#right_video src='"+movie+"'/>";
    pics += "<div.bar><button#video_play>"+ play +"</button><button#right_volume>"+ volume +"</button></div>";
    pics += "</div>";
    $(#stack).prepend(pics);
    //滚动图
    $(#stack).refresh();
    $(#stack).update();
    rotateAttached($(#stack));

    var vdom = $(#right_video);
    try{
        vdom.videoStop();
        vdom.videoUnload();
        vdom.videoLoad(movie)
        vdom.videoPlay();
        vdom.audioVolume(VIDEO_VALUME);
        VIDEO_PLAY_STATE = 1;
        var widthMode = $(#right).attributes.hasClass("widthMode");
        var videoWidth = vdom.videoWidth();
        var videoHeight = vdom.videoHeight();
        var wrapperWidth = $(#right_video_wrapper).box(#width);
        if(widthMode){
            //宽模式
            var maxHeight = 200;
            var prop = videoWidth.toFloat() / videoHeight.toFloat()
            prop = Math.floor(prop * 100) / 100;
            var realWidth = maxHeight * prop;
            $(#right_video_wrapper).style["width"] = realWidth;
            $(#right_video_wrapper).style["height"] = maxHeight;
        }else{
            //窄模式
            var prop = videoHeight.toFloat() / videoWidth.toFloat()
            prop = Math.floor(prop * 100) / 100;
            //var realHeight = (wrapperWidth.toFloat() * prop) * 0.8;
            var realHeight = (wrapperWidth.toFloat() * prop);
            realHeight = realHeight.toInteger();
            $(#right_video_wrapper).style["height"] = realHeight + "dip";
        }
    }catch(e){}
   
    $(#right_video).onControlEvent = function(evt) {
        if(evt.type == Event.VIDEO_STOPPED){
            if ( this.videoDuration() - this.videoPosition()  < 1){
                this.videoPlay(0.0);
            }
        }
    }
}

//滑动图
function _createRotate(info){
    var pic = {}; //图片列表
    var sort = JSON.parse(CONF.Default.ThumbOrders);
    var defSort = ["title","thumb","snap","poster","packing","cassette","icon","gif","background","wallpaper","video"]; //默认排序方式
    if(info != false){
        var name = getRomName(info.RomPath);

        pic["title"] = getRomPicPath("title",info.Platform,name,true); //标题图
        pic["thumb"] = getRomPicPath("thumb",info.Platform,name,true); //展示图
        pic["snap"] = getRomPicPath("snap",info.Platform,name,true); //插画
        pic["poster"] = getRomPicPath("poster",info.Platform,name,true); //海报
        pic["packing"] = getRomPicPath("packing",info.Platform,name,true); //包装盒图
        pic["cassette"] = getRomPicPath("cassette",info.Platform,name,true); //卡带图
        pic["icon"] = getRomPicPath("icon",info.Platform,name,true); //图标图
        pic["gif"] = getRomPicPath("gif",info.Platform,name,true); //gif图
        pic["background"] = getRomPicPath("background",info.Platform,name),true; //背景图
        pic["wallpaper"] = getRomPicPath("wallpaper",info.Platform,name,true); //壁纸
        pic["video"] = getRomPicPath("video",info.Platform,name,true); //视频
    }else{
        $(#rotate).style["display"] = "none";
        $(#stack).html = "";
        return;
    }

    var stackHtml = "";    
    for (var s in sort){
        if(pic[s] != undefined && pic[s].length > 0){
            for (var ps in pic[s]){
                if (s == "video"){
                    continue;
                }
                stackHtml += "<img src=\""+ ps +"\" />";
            }
        }
    }

    //防止排序字段不足导致的图片无法显示
    if(pic.length > sort.length){
        for (var p in defSort){
            if(sort.indexOf(p) == -1 && pic[p].length > 0){
                for (var ps in pic[p]){
                    if (p == "video"){
                        continue;
                    }
                    stackHtml += "<img src=\""+ ps +"\" />";
                }
            }
        }
    }

    //如果没有图片，则显示默认图片
     if(stackHtml == "" && CONF.Theme[CONF.Default.Theme].Params["default-thumb-image"] != undefined){
        stackHtml = "<img src=\""+ URL.fromPath(CONF.RootPath + "theme/" + CONF.Theme[CONF.Default.Theme].Params["default-thumb-image"]) +"\" />";
     }
    //滚动图
    $(#stack).html = stackHtml;
    $(#stack).refresh();
    $(#stack).update();
    rotateAttached($(#stack));

    //重新计算下尺寸，防止标题图过宽
    var stackList = $$(#stack img);
    for (var s in stackList){
        var width = s.box(#width);
        var height = s.box(#height);
        if(width == 0 || height == 0){
            continue;
        }
        if(width / height > 1.5){
            s.style["max-width"] = "300dip";
            s.style["height"] = "auto";
        }
    }

    //加载视频和视频状态控制
    VIDEO_ALLOW_PLAY = 1;
    self.timer(1s, function(){
        _setVideoControl(pic["video"]);
    });

    if($$(#stack > *).length <= 1){
        $(#rotate).style["display"] = "none";
    }else{
        $(#rotate).style["display"] = undefined;
    }

    //预先加载第一张图
    var firstDom = $(#stack img);
    if(firstDom != undefined && firstDom.attributes["data-src"] != "" && firstDom.attributes["data-src"] != undefined){
        firstDom.attributes["src"] = firstDom.attributes["data-src"];
        firstDom.attributes["data-src"] = "";
    }

    return pics;
}

//图集选项卡
function _createSecondThumbs(info){

    var getjson = mainView.GetGameById(ACTIVE_ROM_ID);
    var info = JSON.parse(getjson);
    var pfid = info.Platform;
    var name = getRomName(info.RomPath)
    var platform = CONF.Platform[pfid.toString()];

    var pic = {};

    pic["title"] = getRomPicPath("title",pfid,name,true);
    pic["thumb"] = getRomPicPath("thumb",pfid,name,true);
    pic["snap"] = getRomPicPath("snap",pfid,name,true);
    pic["poster"] = getRomPicPath("poster",pfid,name,true);
    pic["packing"] = getRomPicPath("packing",pfid,name,true);
    pic["cassette"] = getRomPicPath("cassette",pfid,name,true);
    pic["icon"] = getRomPicPath("icon",pfid,name,true);
    pic["gif"] = getRomPicPath("gif",pfid,name,true);
    pic["background"] = getRomPicPath("background",pfid,name,true);
    pic["wallpaper"] = getRomPicPath("wallpaper",pfid,name,true);
    pic["video"] = getRomPicPath("video",pfid,name,true);

    var sort = JSON.parse(CONF.Default.ThumbOrders);
    var defSort = ["title","thumb","snap","poster","packing","cassette","icon","gif","background","wallpaper","video"]; //默认排序方式

    var html = "";    
    for (var s in sort){
        if(pic[s] != undefined){
            html += _getSideSecondDom(s,pic[s]);
        }
    }

    //防止排序字段不足导致的图片无法显示
    if(pic.length > sort.length){
        for (var p in defSort){     
            if(sort.indexOf(p) == -1){
                html += _getSideSecondDom(p,pic[p]);
            }
        }
    }

    $(#second_thumbs).html = html;
}


//获取侧边栏图集
function _getSideSecondDom(type,urls){
    var html = "";
    var title = ""; //标题
    var imgContent = "";
    switch(type){
        case "thumb":
            title = CONF.Lang.Thumb;
            break;
        case "snap":
            title = CONF.Lang.Snap;
            break;
        case "title":
            title = CONF.Lang.TitlePic;
            break;
        case "poster":
            title = CONF.Lang.Poster;
            break;
        case "packing":
            title = CONF.Lang.Packing;
            break;
        case "cassette":
            title = CONF.Lang.CassettePic;
            break;
        case "icon":
            title = CONF.Lang.Icon;
            break;
        case "gif":
            title = CONF.Lang.GifPic;
            break;
        case "background":
            title = CONF.Lang.BackgroundPic;
            break;
        case "wallpaper":
            title = CONF.Lang.WallpaperPic;
            break;
        case "video":
            title = CONF.Lang.Video;
            break;
    }

    html += "<h2>"+title+"</h2>";
    html += "<div class='filedropzone_widget' id='filedropzone_"+ type +"'>";
    var i = 0;
    if(type  == "video"){
        for (var url in urls){
            var filenameArr = getFileName(url).split("__");
            var sid = filenameArr[1] == undefined ? "" : filenameArr[1];
            html += "<div.file-drop-zone opt='video' sid='"+ sid +"' accept-drop='"+ VIDEO_FILTER +"'>";
            html += "<div class='file-drop-zone-empty file-drop-zone-isset'>("+ CONF.Lang.VideoCanNotBePreviewed + ")<br>"+CONF.Lang.VideoUrl + "<br>" + URL.toPath(url) + "</div>";
            html += "</div>";
        }
    }else{
        var visible = $(#right_second).style["visibility"];
        for (var url in urls){
            var filenameArr = getFileName(url).split("__");
            var sid = filenameArr[1] == undefined ? "" : filenameArr[1];
            html += "<div.file-drop-zone opt='"+type+"' sid='"+ sid +"' accept-drop='"+ PIC_FILTER +"'>";
            if(visible == "visible"){
                html += "<img data-src=\"\" src='"+url+"'>";
            }else{
                html += "<img src=\"\" data-src='"+url+"'>";
            }
            html += "</div>";
        }
    }

    //新增空模块
    if(type == "video"){
        if(urls.length == 0) html += createEmptyFileDropZone(type);
    }else{
        html += createEmptyFileDropZone(type);
    }

    html += "</div>";
    return html;
}

//游戏启动（侧边栏）
function sidebarRunGame(evt){
    //videoPause();
    var simdom = $(#sim_select).select("li:current");
    var sim = "";
    if(simdom != undefined){
        sim = simdom.attributes["sim"];
    }
    view.RunGame(evt.attributes["rid"],sim);
}

//侧边栏切换模拟器
function switchRomSim(evt){
     evt.state.current = true;
    var romid = evt.attributes["rom"];
    var simid = evt.attributes["sim"];
    view.SetRomSimulator(romid,simid);
}

//设置rom的cmd
function SetRomCmd(evt){
   var obj = evt.parent;
        var romId = obj.attributes["rom"];
        var simId = obj.attributes["sim"];
        view.dialog({
            url:self.url(ROOTPATH + "edit_rom_cmd.html"),
            width:self.toPixels(500dip),
            height:self.toPixels(310dip),
            parameters: {
                romId:romId,
                simId:simId,
            }
        });
}

//侧边栏缩略图滑动特效
function thumbSlider(evt){
    var container = $(#stack);
    var next = container.shown.next || container.first;
    rotateTo($(#stack),next,false);
}


//侧边栏设置我的喜爱
function setFavorite(evt){
    var id = evt.attributes["rid"];
    var star = evt.attributes["value"];
    var rdom = $(#romlist).select("li[rid="+id+"]");

    if(star == "1"){
        star = "0";
        evt.html = "<img src='this://app/images/fileico_fav_0.png' />";
        if(rdom != undefined){
            if($(#switch_romlist).attributes["value"] == 1){
                rdom.select(".rom_star").style["display"] = "none";
            }else{
                rdom.select("div").attributes.removeClass("fav");
            }
        }
    }else{
        star = "1";
        evt.html = "<img src='this://app/images/fileico_fav_1.png' />";
        if(rdom != undefined){
            if($(#switch_romlist).attributes["value"] == 1){
                rdom.select(".rom_star").style["display"] = "block";
            }else{
                rdom.select("div").attributes.addClass("fav");
            }
        }
    }

    evt.attributes["value"] = star;

    var result = view.SetFavorite(id,star);
   
    if(result != "1"){
        alert(result);
    }

    //如果当前是喜好目录，则从rom列表中删除
    var menu = $(#menulist).select("dd:current").attributes["opt"];
    if(menu == "favorite" && rdom != undefined){
        rdom.remove();
    }
}

//侧边栏设置隐藏
function setHide(evt){
    var id = evt.attributes["rid"];
    var hide = evt.attributes["value"];

    if(hide == "1"){
        hide = "0";
        evt.html = "<img src='this://app/images/fileico_hide_0.png' />";

    }else{
        hide = "1";
        evt.html = "<img src='this://app/images/fileico_hide_1.png' />";     
    }

    evt.attributes["value"] = hide;

    var result = view.SetHide(id,hide);
   

    if(result != "1"){
        alert(result);
    }

    //从rom列表中删除dom
    var rdom = $(#romlist).select("li[rid="+id+"]");
    if(rdom != undefined){
        rdom.remove();
    }
}

//控制视频播放
function videoPlay(){
    var video = $(#right_video);
    var position = video.videoPosition();
    var duration = video.videoDuration();
    if (video.videoIsPlaying()){
        VIDEO_PLAY_STATE = 0;
        $(#video_play).html = "<i.video_play></i>";
        video.videoStop();
    }else{
        if (position == duration){
            video.videoPlay(0.0);
        }else{
            video.videoPlay(position);
        }
        VIDEO_PLAY_STATE = 1;
        $(#video_play).html = "<i.video_pause></i>";

    }
}

//控制视频播放 - 暂停
function videoPause(){
    try{
        var video = $(#right_video);
        if(video != undefined){
            var position = video.videoPosition();
            var duration = video.videoDuration();
            if (video.videoIsPlaying()){
                VIDEO_PLAY_STATE = 0;
                video.videoStop();
                if($(#video_play) != undefined){
                    $(#video_play).html = "<i.video_play></i>";
                }
            }
        }
    }catch(e){}
}

//控制视频音量
var VIDEO_VALUME;
function videoVolume(evt){
    var video = $(#right_video);
    if(VIDEO_VALUME == 0 ){
        VIDEO_VALUME = 1;
        evt.html = "<i.sound_off></i>";
    }else{
        VIDEO_VALUME = 0;    
        evt.html = "<i.sound_on></i>";
    }
    video.audioVolume(VIDEO_VALUME);
    view.UpdateConfig("video_volume",VIDEO_VALUME);
}

//打开文件夹菜单
function createOpenFolderMenu(info){

    var pfId = info.Platform.toString();
    var simId = 0;
    var simIco = "";
    if (info.SimId != 0){
        simId = info.SimId;
    }else if (CONF.Platform[pfId].UseSim != undefined){
        simId = CONF.Platform[pfId].UseSim["Id"];
    }
    var btndata = "<li class='folder' rid='"+ info.Id +"' opt='rom'>"+ CONF.Lang.LocationRomFile +"</li>";
    
    btndata += "<li class='folder'>"+ CONF.Lang.LocationResDir +"<menu>";

    if(CONF.Platform[pfId].ThumbPath != ""){
        btndata += "<li class='folder' rid='"+ info.Id +"' opt='thumb'>"+ CONF.Lang.LocationThumbFile +"</li>";
    }
    
    if(CONF.Platform[pfId].SnapPath != ""){
        btndata += "<li class='folder' rid='"+ info.Id +"' opt='snap'>"+ CONF.Lang.LocationSnapFile +"</li>";
    }

    if(CONF.Platform[pfId].PosterPath != ""){
        btndata += "<li class='folder' rid='"+ info.Id +"' opt='poster'>"+ CONF.Lang.OpenPosterFolder +"</li>";
    }

    if(CONF.Platform[pfId].PackingPath != ""){
        btndata += "<li class='folder' rid='"+ info.Id +"' opt='packing'>"+ CONF.Lang.OpenPackingFolder +"</li>";
    }
   if(CONF.Platform[pfId].TitlePath != ""){
        btndata += "<li class='folder' rid='"+ info.Id +"' opt='title'>"+ CONF.Lang.LocationTitleFile +"</li>";
    }

    if(CONF.Platform[pfId].CassettePath != ""){
            btndata += "<li class='folder' rid='"+ info.Id +"' opt='cassette'>"+ CONF.Lang.LocationCassetteFile +"</li>";
    }

    if(CONF.Platform[pfId].IconPath != ""){
            btndata += "<li class='folder' rid='"+ info.Id +"' opt='icon'>"+ CONF.Lang.LocationIconFile +"</li>";
    }

    if(CONF.Platform[pfId].GifPath != ""){
            btndata += "<li class='folder' rid='"+ info.Id +"' opt='gif'>"+ CONF.Lang.LocationGifFile +"</li>";
    }
    if(CONF.Platform[pfId].BackgroundPath != ""){
        btndata += "<li class='folder' rid='"+ info.Id +"' opt='background'>"+ CONF.Lang.LocationBackgroundFile +"</li>";
    }

    if(CONF.Platform[pfId].WallpaperPath != ""){
        btndata += "<li class='folder' rid='"+ info.Id +"' opt='wallpaper'>"+ CONF.Lang.LocationWallpaperFile +"</li>";
    }

    if(CONF.Platform[pfId].VideoPath != ""){
        btndata += "<li class='folder' rid='"+ info.Id +"' opt='video'>"+ CONF.Lang.LocationVideoFile +"</li>";
    }

    if(CONF.Platform[pfId].AudioPath != ""){
        btndata += "<li class='folder' rid='"+ info.Id +"' opt='audio'>"+ CONF.Lang.LocationAudioFile +"</li>";
    }

    btndata +="<hr>";

    if(CONF.Platform[pfId].DocPath != ""){
        btndata += "<li class='folder' rid='"+ info.Id +"' opt='doc'>"+ CONF.Lang.LocationDocFile +"</li>";
    }

    if(CONF.Platform[pfId].StrategyPath != ""){
        btndata += "<li class='folder' rid='"+ info.Id +"' opt='strategy'>"+ CONF.Lang.LocationStrategyFile +"</li>";
    }

    if(CONF.Platform[pfId].FilesPath != ""){
        btndata += "<li class='folder' rid='"+ info.Id +"' opt='files'>"+ CONF.Lang.LocationStrategyFiles +"</li>";
    }

    btndata += "</menu></li>"
    btndata += "<li class='folder'>"+ CONF.Lang.LocationSimDir +"<menu>"

    //定位模拟器目录
   if(CONF.Platform[pfId].SimList.length > 0){
        for(var simId in CONF.Platform[pfId].SimList) {
            btndata += "<li class='file folder' filename='"+ CONF.Platform[pfId].SimList[simId].Path +"' rid='"+ info.Id +"' sid='"+ CONF.Platform[pfId].SimList[simId].Id +"' opt='sim'>"+ CONF.Lang.LocationSim + CONF.Platform[pfId].SimList[simId].Name.toHtmlString() + "</li>";;
        }

    }

    btndata += "</menu></li>"
    return btndata;
}

//调整侧边栏宽度
function sidebarSize(){
    var width = $(#right).box(#width);
    width = Math.round(width * 0.8).toInteger();
    view.UpdateConfig("panel_sidebar_width",width.toInteger());

    //界面宽窄模式
    var right = $(#right);
    var has = right.attributes.hasClass("widthMode")
    if(width >= 500){
        if(has == false){
            $(#right).attributes.addClass("widthMode");
        }
    }else{
        if(has == true){
            $(#right).attributes.removeClass("widthMode");
        }
    }
}

//打开文件夹 - 图集按钮
//右键菜单定位目录
function openFolderBySideSecond(evt){
    if(evt.attributes["opt"] == "" || evt.attributes["opt"] == undefined){
        return;
    }

    var opt = evt.attributes["opt"];
    var sid = evt.attributes["sid"];

    if(sid == THUMB_EMPTY_SID){
        sid = "";
    }
    view.OpenFolder(ACTIVE_ROM_ID,opt,sid);
}


//更新侧边栏平台介绍
function sidebarPlatformDesc(platform){
    var desc = CONF.Platform[platform.toString()] == undefined ? "" : CONF.Platform[platform.toString()].Desc;

    $(#right).scrollTo(0,0, false);

    if(desc == ""){
        //没有简介，显示遮罩
        $(#mask).style["display"] = "block";
        $(#platform_desc).style["display"] = "none";
        $(#platform_desc).html = "";
    }else{

        //正则替换中文
        desc = desc.replace(/<img [^>]*src=['"]([^'"]+)[^>]*>/gi, function (match, capture) {
            return "<img src=\"" + URL.fromPath(capture) + "\" />";
        });

        //有简介，显示内容
        $(#mask).style["display"] = "none";
        $(#platform_desc).style["display"] = "block";
        $(#platform_desc).html = desc;
    }
}

//启动游戏后回调
function CB_runGame(){
    //禁止播放视频
    VIDEO_ALLOW_PLAY = 0;
}

//模拟器向下切换
function nextSim(){
    //如果界面被隐藏，则跳过此功能
    if ($(#right).style["display"] == "none"){
        return;
    }

    //如果没有或只有一个模拟器，则跳过此功能=
    if ($$(#sim_select li) == undefined){
        return;
    }

    if ($$(#sim_select li).length < 2){
        return;
    }

    var nextNum = $(#sim_select).select("li:current").index + 1;
    var count = $$(#sim_select li).length;
    var nextDom;
    if (nextNum +1 > count){
        nextDom = $(#sim_select).select("li:nth-child(1)");
    }else{
        nextNum ++ ;
        nextDom = $(#sim_select).select("li:nth-child("+ nextNum +")");
    }
    if(nextDom != undefined){
        switchRomSim(nextDom);
    }
}

//子游戏向下切换
function nextSubGame(){
    //如果界面被隐藏，则跳过此功能
    if ($(#right).style["display"] == "none"){
        return;
    }

    //如果没有或只有一个模拟器，则跳过此功能=
    if ($$(#buttons li) == undefined){
        return;
    }

    if ($$(#buttons li).length < 2){
        return;
    }

    if ($(#buttons li:current) == undefined){
        $(#buttons).select("li:nth-child(1)").state.current = true;
    }

    var nextNum = $(#buttons).select("li:current").index + 1;
        
    var count = $$(#buttons li).length;

    var nextDom;
    if (nextNum +1 > count){
        nextDom = $(#buttons).select("li:nth-child(0)");
    }else{
        nextNum ++ ;
        nextDom = $(#buttons).select("li:nth-child("+ nextNum +")");
    }
    if(nextDom != undefined){
        nextDom.state.current = true;
    }
}

//生成评分星级
function createScore(scoreStr){

    if(scoreStr == ""){
        scoreStr = "0.0";
    }
    var score = scoreStr.toFloat();

    if(score == "" || score == undefined){
        return;
    }

    var scoreHtml = score.toFixed(1);
    score = (Math.round(scoreHtml.toFloat() * 2)) / 2;

    var spans = $$(#score span);
    $(#score_num).html = scoreHtml;
    for (var span in spans){
        if(score >= 1){
            span.attributes["stat"] = "full";
            span.attributes["class"] = "full";
        }else if(score <= 0){
            span.attributes["stat"] = "empty";
            span.attributes["class"] = "empty";
        }else{
            span.attributes["stat"] = "half";
            span.attributes["class"] = "half";
        }
        score = score - 1 ;
    }

}

//设置星级
function setScore(){
    var spans = $$(#score span);
    var count = 0.0;
    for (var span in spans){
        var cls = span.attributes["class"];
        var setCls = "empty";
        if(cls == "full" || cls == "h-full"){
            setCls = "full";
            count += (1.0).toFloat();
        }else if(cls == "half"|| cls == "h-half") {
            setCls = "half";
            count += (0.5).toFloat();
        }
        span.attributes["stat"] = setCls;
        span.attributes["class"] = setCls;
    }

    var score = count.toFixed(1).toString();

    view.SetScore(ACTIVE_ROM_ID,score);
    $(#score_num).html = score;
    $(#score_num).attributes.addClass("actShake");
    $(#score_num).attributes.removeClass("actShake");

}

//设置星级 - 鼠标移入
function mousemoveScore(obj,e){
    var score = $(#score);
    var max = 5;
    var (scoreX,scoreY) = obj.box(#position, #inner, #view);
    var currentX = e.x < scoreX ? scoreX : e.x;
    var starWidth = obj.box(#width); //星星宽度
    var index = obj.index + 1;
    if(currentX - scoreX >= starWidth/2){
        obj.attributes["class"] = 'h-full';
    }else{
        obj.attributes["class"] = 'h-half';
    }
    if(index > 1){
        for (let k = 1; k < index; k++) {
            score.select("span:nth-child("+ k +")").attributes["class"] = 'h-full';
        }
    }
    if(index < max){
        for (let j = index+1; j <= 5; j++) {
            score.select("span:nth-child("+ j +")").attributes["class"] = 'empty';
        }
    }
};

//设置星级 - 鼠标移出
function mouseleaveScore(){
    var spans = $$(#score span);
    for (var span in spans){
        var stat = span.attributes["stat"] == "" ? "empty" : span.attributes["stat"];
        span.attributes["class"] = stat;
    }
};

//读取通关状态
function createGameComplete(stat){
    var complete = CONF.Lang.GamePlaying;
    if(stat == 1){
        complete = CONF.Lang.GameComplete;
    }else if(stat == 2){
        complete = CONF.Lang.GamePlatinumComplete;
    }
    $(#game_complete .info).html = complete;
    $(#game_complete).attributes["class"] = "complete" + stat;
    $(#game_complete).attributes["opt"] = stat;
}

//设置通关状态
function setGameComplete(obj){

    var stat = $(#game_complete).attributes["opt"];
    var newStat = 0;
    var complete = "";

    switch(stat){
        case "0":
            //未通关，改为已通关
            newStat = 1;
            complete = CONF.Lang.GameComplete;
            
            break;
        case "1":
            //已通关，改为完美通关
            newStat = 2;
            complete = CONF.Lang.GamePlatinumComplete;
            break;
        case "2":
            //完美通关，改为未通关
            newStat = 0;
            complete = CONF.Lang.GamePlaying;
            break;
    }

    view.SetComplete(ACTIVE_ROM_ID,newStat);

    $(#game_complete .info).html = complete;
    $(#game_complete).attributes["class"] = "complete" + newStat;
    $(#game_complete).attributes["opt"] = newStat;

    obj.attributes.addClass("actShake");
    obj.attributes.removeClass("actShake");

}
