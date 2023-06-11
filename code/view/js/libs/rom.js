//创建游戏rom dom
function createRomList(cls = 0) {

    if (ROMJSON == undefined || ROMJSON == "") {
        return;
    }

    var romobj = JSON.parse(ROMJSON);
    var romlist = self.$(#romlist);
    var switch_romlist = $(#switch_romlist).attributes["value"];
    var data = "";

    //如果清理参数=1，则清理掉所有数据
    if (cls == 1) {

        //清空rom选定
        if ($(#romlist li:current) != undefined) {
            $(#romlist li:current).state.current = false;
        }

        //重置翻页数据
        resetPageScroll();
        //清理游戏列表
        romlist.clear();

        //显示出加载更多按钮
        //$(#load_more).style["display"] = "block";

    }

    //rom列表缩略图方向
    initRomlistThumbDirection();

    //更新romlist字体大小
    initRomlistFontsize();

    //更新romlist模块间距
    initRomlistMargin();

    //更新romlist 模块大小
    initRomlistSize();

    if (switch_romlist == 1) {
        //模块模式
        data = _createRomListBlock(romobj);
    } else {
        //列表模式
        data = _createRomListList(romobj, cls);

    }

    romlist.append(data);

    //动态载入rom
    romlist = $$(#romlist li);

    SCROLL_PAGE++;

    //验证是否显示【加载更多】
    //CONST_ROM_LIST_PAGE_SIZE

    try {
        //异步加载图片
        if (switch_romlist == 1) {
            dynamicLoadImages(romlist, romobj.length, 31);
        }
    } catch (e) {
        alert(e);
    }

}

var asyncMaxNum = maxNum;
var asyncRunning = 0;
var asyncQueue = [];

function dynamicLoadImages(romlist, romMaxLength, maxNum) {

    if (romlist.length == 0) {
        return;
    }

    var sum = romlist.length - romMaxLength < 0 ? 0 : romlist.length - romMaxLength;
    if (sum > 0) {
        sum--;
    }
    if (sum > 0) {
        romlist.splice(0, sum);
    }
    asyncMaxNum = maxNum;
    asyncQueue = romlist;
    asyncRunning = 0;
    dynamicLoadImagesNext();

}

function dynamicLoadImagesNext() {
    while (asyncRunning < asyncMaxNum && asyncQueue.length) {
        const task = asyncQueue.shift();
        var imageDom = task.select(".rom_thumb");
        if (imageDom != undefined && imageDom.attributes["data-src"] != "" && imageDom.attributes["data-src"] != undefined) {
            asyncRunning++;
            task.timer(1ms, async function () {
                await imageDom.attributes["src"] = imageDom.attributes["data-src"];
                imageDom.attributes["data-src"] = "";
                asyncRunning--;
                dynamicLoadImagesNext();
            });
        }
    }
}

//列表样式数据 - 模块
function _createRomListBlock(romobj) {
    var title = "";
    var thumb = "";
    var data = "";
    var name = "";

    for (var obj in romobj) {
        var romName = getRomName(obj.RomPath)
        thumb = getRomPicPath("", obj.Platform, romName);

        //没有找到图片，读取默认缩略图
        if (thumb == "" && CONF.Theme[CONF.Default.Theme].Params["default-thumb-image"] != undefined) {
            thumb = URL.fromPath(CONF.RootPath + "theme/" + CONF.Theme[CONF.Default.Theme].Params["default-thumb-image"]);
        }

        if (CONF.Default.RomlistNameType == 0) {
            name = obj.Name.toHtmlString(); //别名
        } else {
            name = romName.toHtmlString(); //文件名
        }

        data += "<li class='romitem' rid='" + obj.Id + "' title='" + title + "'>";
        data += "<div><img class='rom_thumb' src=\"\" data-src='" + thumb + "' /></div>";
        data += "<div.name>" + name + "</div>";
        if (obj.Star == 1) { data += "<span class='rom_star' style='display:block'></span>"; }
        else { data += "<span class='rom_star' style='display:none'></span>"; }
        data += "</li>";
    }
    return data;
}

//列表样式数据 - 列表
function _createRomListList(romobj, cls) {
    var title = "";
    var romPath = "";
    var data = "";

    //标题列
    var column = CONF.Default.RomlistColumn.split(",");
    if (cls == 1) {
        data += "<li class='romth'>";
        data += "<div>" + CONF.Lang["BaseName"] + "</div>";
        if (column[0] == 1) { data += "<div.menu>" + CONF.Lang["FilterMenu"] + "</div>"; }
        if (column[1] == 1) { data += "<div.type>" + CONF.Lang["BaseType"] + "</div>"; }
        if (column[2] == 1) { data += "<div.year>" + CONF.Lang["BaseYear"] + "</div>"; }
        if (column[7] == 1) { data += "<div.producer>" + CONF.Lang["BaseProducer"] + "</div>"; }
        if (column[3] == 1) { data += "<div.publisher>" + CONF.Lang["BasePublisher"] + "</div>"; }
        if (column[4] == 1) { data += "<div.country>" + CONF.Lang["BaseCountry"] + "</div>"; }
        if (column[5] == 1) { data += "<div.translate>" + CONF.Lang["BaseTranslate"] + "</div>"; }
        if (column[6] == 1) { data += "<div.version>" + CONF.Lang["BaseVersion"] + "</div>"; }
        if (column[9] == 1) { data += "<div.score>" + CONF.Lang["Score"] + "</div>"; }
        if (column[10] == 1) { data += "<div.complete>" + CONF.Lang["GameCompleteState"] + "</div>"; }
        if (column[8] == 1) { data += "<div.path>" + CONF.Lang["FilePath"] + "</div>"; }


        data += "</li>";
    }
    for (var obj in romobj) {

        if (obj.Name.toHtmlString().length > 20) {
            title = obj.Name.toHtmlString();
        } else {
            title = "";
        }

        var romName = getRomName(obj.RomPath);
        var name = "";
        if (CONF.Default.RomlistNameType == 0) {
            name = obj.Name.toHtmlString(); //别名
        } else {
            name = romName.toHtmlString(); //文件名
        }

        var fav = obj.Star == 1 ? ".fav" : "";

        var complete = "";
        if (obj.Complete == 0) complete = CONF.Lang.GamePlaying;
        else if (obj.Complete == 1) complete = CONF.Lang.GameComplete;
        else if (obj.Complete == 2) complete = CONF.Lang.GamePlatinumComplete;

        data += "<li class='romitem' rid='" + obj.Id + "' title='" + title + "'>";
        data += "<div.name" + fav + ">" + name + "</div>";
        if (column[0] == 1) {
            var menu = obj.Menu == "_7b9" ? "" : obj.Menu;
            data += "<div.menu>" + menu + "</div>";
        }
        if (column[1] == 1) { data += "<div.type>" + obj.BaseType + "</div>"; }
        if (column[2] == 1) { data += "<div.year>" + obj.BaseYear + "</div>"; }
        if (column[7] == 1) { data += "<div.producer>" + obj.BaseProducer + "</div>"; }
        if (column[3] == 1) { data += "<div.publisher>" + obj.BasePublisher + "</div>"; }
        if (column[4] == 1) { data += "<div.country>" + obj.BaseCountry + "</div>"; }
        if (column[5] == 1) { data += "<div.translate>" + obj.BaseTranslate + "</div>"; }
        if (column[6] == 1) { data += "<div.version>" + obj.BaseVersion + "</div>"; }
        if (column[9] == 1) { data += "<div.score>" + obj.Score + "</div>"; }
        if (column[10] == 1) { data += "<div.complete>" + complete + "</div>"; }
        if (column[8] == 1) { data += "<div.path>" + obj.RomPath + "</div>"; }
        data += "</li>"
    }
    return data;

}

//rom分页
function scrollLoadRom(evtPos) {

    var scrollPos = evtPos.toInteger() + $(#center_content).box(#height);
    var boxHeight = $(#romwrapper).box(#height);

    if (SCROLL_POS == 0) {
        SCROLL_POS = scrollPos;
    }

    if ((boxHeight - scrollPos <= 100) && (scrollPos > SCROLL_POS)) {
        //加载一页rom
        loadPageRom();
    }

}

//加载一页rom
function loadPageRom() {
    //如果加锁中，则不执行后续逻辑，防止重复触发
    if (SCROLL_LOCK == true) {
        return;
    }
    SCROLL_LOCK = true; //加锁

    var menu = $(#menulist).select("dd:current").attributes["opt"];
    var search = $(#search_input).value;

    var num = $(#num_search).select("li:current").html;

    if (num == "ALL") {
        num = "";
    }

    //生成游戏列表
    var filter_type = $(#filter_type).value === undefined ? "" : $(#filter_type).value;
    var filter_publisher = $(#filter_publisher).value === undefined ? "" : $(#filter_publisher).value;
    var filter_year = $(#filter_year).value === undefined ? "" : $(#filter_year).value;
    var filter_country = $(#filter_country).value === undefined ? "" : $(#filter_country).value;
    var filter_translate = $(#filter_translate).value === undefined ? "" : $(#filter_translate).value;
    var filter_version = $(#filter_version).value === undefined ? "" : $(#filter_version).value;
    var filter_producer = $(#filter_producer).value === undefined ? "" : $(#filter_producer).value;
    var filter_score = $(#filter_score).value === undefined ? "" : $(#filter_score).value;
    var filter_complete = $(#filter_complete).value === undefined ? "" : $(#filter_complete).value;


    var hide = ACTIVE_MENU == "hide" ? 1 : 0;

    var req = {
        "showHide": hide,
        "platform": ACTIVE_PLATFORM.toInteger(),
        "catname": menu.toString(),
        "keyword": search.toString(),
        "num": num,
        "page": SCROLL_PAGE,
        "baseType": filter_type.toString(),
        "basePublisher": filter_publisher.toString(),
        "baseYear": filter_year.toString(),
        "baseCountry": filter_country.toString(),
        "baseTranslate": filter_translate.toString(),
        "baseVersion": filter_version.toString(),
        "baseProducer": filter_producer.toString(),
        "score": filter_score.toString(),
        "complete": filter_complete.toString(),
    };

    var request = JSON.stringify(req);

    ROMJSON = view.GetGameList(request);

    if (ROMJSON != "[]") {
        createRomList(0);
        SCROLL_POS = scrollPos;
        SCROLL_LOCK = false;
    } else {
        $(#load_more).style["display"] = "none";
    }
}


//初始化滚动分页功能
function resetScroll() {
    SCROLL_PAGE = 0;
    SCROLL_LOCK = false;
    SCROLL_POS = 0
}

//字母搜索
function numSearch(evt) {

    //重复点击拦截
    var current = $(#num_search).select("li:current");
    if (current != undefined && current.html == evt.html) {
        return;
    }

    var menu = $(#menulist).select("dd:current").attributes["opt"];
    var num = evt.html;

    if (num == "ALL") {
        num = "";
    }

    //重置滚动翻页数据
    resetScroll();

    //重置筛选项
    resetFilterOptions();
    var hide = ACTIVE_MENU == "hide" ? 1 : 0;


    var req = {
        "showHide": hide,
        "platform": ACTIVE_PLATFORM.toInteger(),
        "catname": menu,
        "page": SCROLL_PAGE,
        "num": num,
    };
    var request = JSON.stringify(req);

    ROMJSON = view.GetGameList(request);
    createRomList(1);

    //生成rom数量
    var romCount = view.GetGameCount(request);
    $(#rom_count_num).html = romCount; //游戏数量

    //激活当前，改变样式
    evt.state.current = true;
}

//rom列表运行游戏
function romListRunGame(evt) {
    if (evt == undefined) {
        return;
    }
    //videoPause();
    var getjson = view.GetGameById(evt.attributes["rid"]);
    var info = JSON.parse(getjson);
    view.RunGame(evt.attributes["rid"], info.SimId);
}

//数字按键启动游戏
function numRunGame(num) {
    var simDom = $(#sim_select).select("li:nth-child(" + num + ")");
    if (simDom != undefined) {
        var romid = simDom.attributes["rom"];
        var simid = simDom.attributes["sim"];
        view.RunGame(romid, simid);
    }

}