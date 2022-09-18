
/**
 * 文件去重
 **/

//检测重复 - 文件去重
function managerCheckRepeat(){
    var platformId = $(#platform_repeat).value;

    if (platformId == "0"){
        alert(CONF.Lang.SelectPlatform);
        return;
    }

    var romjson = mainView.CheckRomRepeat(platformId);
    var romobj = JSON.parse(romjson);

    if(romobj.length == 0 ){
        $(#romlist_repeat tbody).html = "";
        alert(CONF.Lang.DuplicateFilesNotFound);
        return;
    }

    var max= 3000;
    if(romobj.length > max ){
        alert(CONF.Lang.TipDuplicateFilesFound);
    }

    var romData = "";
    var upSize = 0;
    var color1 = "grp1";
    var color2 = "grp2";
    var clr = color1;
    for(var i=0;i<=max;i++) {
        if(romobj[i] == undefined){
            break;
        }
        //颜色
        if(upSize == 0){upSize = romobj[i].size;}
        if (romobj[i].size != upSize){clr = (clr == color1) ? color2 : color1}
        upSize = romobj[i].size

        romData += "<tr opt=\"" + romobj[i].path + "\" rid='"+ romobj[i].id +"' class='"+clr+"'>";
        romData += "<td><input class='check_repeat' type='checkbox' value='"+ romobj[i].id +"' /></td>";
        romData += "<td>"+ romobj[i].name +"</td>";
        romData += "<td>"+ romobj[i].path +"</td>";
        romData += "<td>"+ romobj[i].size +"</td>";
        romData += "<td><button class='rungame'>"+CONF.Lang.Run+"</button></td>";
        romData += "</tr>"
    }
    $(#romlist_repeat tbody).html = romData;



}

//移动文件 - 文件去重
function managerFileRepeatMoveRom(){

    var checks = $$(#romlist_repeat .check_repeat:checked);
    if(checks.length == 0){
        alert(CONF.Lang.NotSelectGames);
        return;
    }

    var paths = [];
    for(var c in checks){
        paths.push(c.parent.parent.attributes["opt"]);
    }

    var url = view.selectFolder(CONF.Lang.SelectFolder);
    if(url){
        url = URL.toPath(url);
        url = url.split("\/").join(SEPARATOR);
        url = url.split(CONF.RootPath.toString()).join("");

        mainView.MoveRomByFileManager(JSON.stringify(paths),url);
    }

    alert(CONF.Lang.RomMoveSuccess);
}

//删除文件 - 文件去重
function managerFileRepeatDelete(){

    var checks = $$(#romlist_repeat .check_repeat:checked);
    if(checks.length == 0){
    alert(CONF.Lang.NotSelectGames);
        return;
    }

    var arr = [];
    for(var c in checks){
        arr.push(c.value);
    }


  //确认窗口
    var result = confirm(CONF.Lang.DeleteRomConfirm,CONF.Lang.DeleteRom);
    if (result != "yes"){
        return true;
    }


    for(var id in arr) {
        mainView.DeleteRom(id); //删除实体文件
        //删除页面上的dom
    }
    alert(CONF.Lang.DeleteSuccess);

}

//运行游戏 - 文件去重
function managerRunGame(obj){
    mainView.RunGame(obj.parent.parent.attributes["rid"],"");
}