//设置展示图显示方向
function setThumbDirection(obj) {

    var pfId = obj.parent.parent.attributes["opt"];
    var dir = obj.attributes["opt"];
    var doms;
    //全部平台
    if (pfId == "all") {

        //确认窗口
        var result = confirm(CONF.Lang.WarningSetAllPlatformThumb, CONF.Lang.Warning);
        if (result != "yes") {
            return true;
        }
        view.UpdateConfig("romlist_direction", dir);
        view.ClearAllPlatformConfig("thumb_direction");
        CONF.Default.RomlistDirection = dir;
    } else {
        view.UpdatePlatformFieldById(pfId.toInteger(), "thumb_direction", dir);
        CONF.Platform[pfId.toString()].ThumbDirection = dir
    }

    initConfig("1");
    createRomList(1);

}

//更改rom列表展示图类型
function setThumbType(obj) {

    var pfId = obj.parent.parent.attributes["opt"];
    var thumb = obj.attributes["opt"];
    //全部平台
    if (pfId == "all") {

        //确认窗口
        var result = confirm(CONF.Lang.WarningSetAllPlatformThumb, CONF.Lang.Warning);
        if (result != "yes") {
            return true;
        }
        view.UpdateConfig("thumb", obj.attributes["opt"]);
        view.ClearAllPlatformConfig("thumb");
    } else {
        view.UpdatePlatformFieldById(pfId.toInteger(), "thumb", thumb);
    }

   //初始化rom\
   var req = {
       "platform": ACTIVE_PLATFORM,
       "catname": $(#menulist).select("dd:current").attributes["opt"],
   };
   var request = JSON.stringify(req);

   ROMJSON = view.GetGameList(request);

    initConfig("1");
    createRomList(1);
}

//是否显示rom的标题背景
function setFontBackgrond(obj) {
    var romlist = $(#romlist);
    romlist.attributes.removeClass("bgshow");
    romlist.attributes.removeClass("bghide");
    romlist.attributes.removeClass("textShadow");

    //先清除所有缩放样式
    for (var i = 0; i <= 2; i++) {
        $(#config_font_background).select("li[opt=" + i + "]").attributes.removeClass("active");
        if (romlist.attributes.hasClass("bg" + i)) {
            romlist.attributes.removeClass("bg" + i);
        }
    }

    if (obj.attributes["opt"] == 0) {
        romlist.attributes.addClass("textShadow");
    }

    $(#config_font_background).select("li[opt=" + obj.attributes["opt"] + "]").attributes.addClass("active");

    romlist.attributes.addClass("bg" + obj.attributes["opt"]);
    view.UpdateConfig("romlist_font_background", obj.attributes["opt"]);
}

//更新rom列表字体
function setFontsize(obj) {
    var pfId = obj.parent.parent.attributes["opt"];
    var size = obj.attributes["opt"];

    //全部平台
    if (pfId == "all") {

        //确认窗口
        var result = confirm(CONF.Lang.WarningSetAllPlatformFontSize, CONF.Lang.Warning);
        if (result != "yes") {
            return true;
        }
        view.UpdateConfig("romlist_font_size", obj.attributes["opt"]);
        view.ClearAllPlatformConfig("thumb_font_size");
    } else {
        view.UpdatePlatformFieldById(pfId.toInteger(), "thumb_font_size", size);
    }

    initConfig("1");
    createRomList(1);
}

//列表图标缩放
function setRomlistSize(obj) {
    var pfId = obj.parent.parent.attributes["opt"];
    var size = obj.attributes["opt"];

    //全部平台
    if (pfId == "all") {

        //确认窗口
        var result = confirm(CONF.Lang.WarningSetAllPlatformThumbSize, CONF.Lang.Warning);
        if (result != "yes") {
            return true;
        }
        view.UpdateConfig("romlist_size", obj.attributes["opt"]);
        view.ClearAllPlatformConfig("thumb_size");
    } else {
        view.UpdatePlatformFieldById(pfId.toInteger(), "thumb_size", size);
    }

    initConfig("1");
    createRomList(1);
}

//列表模块间距
function setRomMargin(obj) {

    var pfId = obj.parent.parent.attributes["opt"];
    var size = obj.attributes["opt"];

    //全部平台
    if (pfId == "all") {

        //确认窗口
        var result = confirm(CONF.Lang.WarningSetAllPlatformThumbMargin, CONF.Lang.Warning);
        if (result != "yes") {
            return true;
        }
        view.UpdateConfig("romlist_margin", obj.attributes["opt"]);
        view.ClearAllPlatformConfig("thumb_margin");
    } else {
        view.UpdatePlatformFieldById(pfId.toInteger(), "thumb_margin", size);
    }

    initConfig("1");
    createRomList(1);
}

//列表列设置
function setRomlistColumn(obj) {
    var column = CONF.Default.RomlistColumn.split(",");
    column[obj.attributes["opt"].toInteger()] = obj.checked == true ? 1 : 0;
    var str = column.join(",");
    view.UpdateConfig("romlist_column", str);
    //更新配置
    CB_createCache();
}

//列表名称显示类型
//0别名；1文件名
function setShowNameType(obj) {
    $(#config_show_name_type).select("li[opt=0]").attributes.removeClass("active");
    $(#config_show_name_type).select("li[opt=1]").attributes.removeClass("active");
    obj.attributes.addClass("active");
    view.UpdateConfig("romlist_name_type", obj.attributes["opt"]);
    //更新配置
    CB_createCache();
}

//列表排序方式
//0别名；1文件名
function setListSort(obj) {
    //高亮当前选项
    var lis = $$(#config_orders menu li);
    for (var s in lis) {
        s.attributes.removeClass("active");
    }
    obj.attributes.addClass("active");

    view.UpdateConfig("romlist_orders", obj.attributes["opt"]);
    //更新配置
    CB_createCache();
}