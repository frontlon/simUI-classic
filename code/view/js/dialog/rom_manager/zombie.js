
/**
 * 无效资源清理
 **/

 //检测无效资源 - 无效资源清理
function managerCheckZombie(){

    var platformId = $(#platform_zombie).value;

    if (platformId == "0"){
        alert(CONF.Lang.SelectPlatform);
        return;
    }

    var romjson = mainView.CheckRomZombie(platformId);
    var romobj = JSON.parse(romjson);

    if(romobj.length == 0 ){
        $(#romlist_zombie tbody).html = "";
        alert(CONF.Lang.NoInvalidFilesNotFound);
        return;
    }
    var max = 3000;
    if(romobj.length > max ){
        alert(CONF.Lang.TipNoInvalidFilesFound);
    }

    var romData = "";
    var type = "";
    var desc = "";
    for(var i=0;i<=max;i++) {
        if(romobj[i] == undefined){
            break;
        }
        if(romobj[i].type == 1){
            type = CONF.Lang.NoInvalidFile;
            desc = CONF.Lang.TipNoInvalidFile;
        }else if(romobj[i].type == 2){
            type = CONF.Lang.DuplicateFiles;
            desc = CONF.Lang.TipDuplicateFiles;
        }else{
            type = CONF.Lang.SubGameFiles;
            desc = CONF.Lang.TipSubGameFiles;
        }

        romData += "<tr>";
        romData += "<td><input class='check_zombie' type='checkbox' value=\""+ romobj[i].path.toHtmlString() +"\" /></td>";
        romData += "<td title='"+ desc +"'>"+ type +"</td>";
        romData += "<td>"+ romobj[i].path.toHtmlString() +"</td>";
        romData += "</tr>"
    }
    $(#romlist_zombie tbody).html = romData;
}

//移动文件 - 无效资源清理
function managerFileZombieMoveRom(){

    var checks = $$(#romlist_zombie .check_zombie:checked);
    if(checks.length == 0){
        alert(CONF.Lang.NotSelectGames);
        return;
    }

    var paths = [];
    for(var c in checks){
        paths.push(c.value);
    }

    var url = view.selectFolder(CONF.Lang.SelectFolder);
    if(url){
        url = URL.toPath(url);
        url = url.split("\/").join(SEPARATOR);
        url = url.split(CONF.RootPath.toString()).join("");

        mainView.MoveZombieByFileManager(JSON.stringify(paths),url);

        alert(CONF.Lang.RomMoveSuccess);
    }

}

//删除文件 - 无效资源清理
function managerFileZombieDelete(){

    var checks = $$(#romlist_zombie .check_zombie:checked);
    if(checks.length == 0){
    alert(CONF.Lang.NotSelectGames);
        return;
    }

    var arr = [];
    for(var c in checks){
        arr.push(c.value);
    }


    //确认窗口
    var result = confirm(CONF.Lang.DeleteFile,CONF.Lang.DeleteFileConfirm);
    if (result != "yes"){
        return true;
    }

    for(var p in arr) {
        mainView.DeleteZombieByFileManager(p); //删除实体文件
    }

    alert(CONF.Lang.DeleteSuccess);

}