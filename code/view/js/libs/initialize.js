//初始化
function Init() {

    //初始化全局配置
    initConfig("0");

    //改变默认窗口尺寸
    view.state = CONF.Default.WindowState.toInteger();

    //窗口面板显示状态
    if (CONF.Default.PanelPlatform.toInteger() == 0) {
        $(#left_platform).style["display"] = "none"
        $(#platform_splitter).style["display"] = "none";
    }
    if (CONF.Default.PanelMenu.toInteger() == 0) {
        $(#left_menu).style["display"] = "none";
        $(#menu_splitter).style["display"] = "none";
    }
    if (CONF.Default.PanelSidebar.toInteger() == 0) {
        $(#right).style["display"] = "none";
        $(#right_splitter).style["display"] = "none";
    }

    //窗口框架宽度
    if (CONF.Default.PanelPlatformWidth != "") {
        $(#left_platform).style["width"] = px(CONF.Default.PanelPlatformWidth.toInteger());
    }

    if (CONF.Default.PanelMenuWidth != "") {
        $(#left_menu).style["width"] = px(CONF.Default.PanelMenuWidth.toInteger());
    }

    if (CONF.Default.PanelSidebarWidth != "") {
        $(#right).style["width"] = CONF.Default.PanelSidebarWidth < 230 ? px(230) : px(CONF.Default.PanelSidebarWidth.toInteger());
        if (CONF.Default.PanelSidebarWidth >= 500) {
            $(#right).attributes.addClass("widthMode");
        }
    }

    //初始化rom
    var pf = CONF.Default.Platform;
    var menu = $(#menulist).select("dd:current").attributes["opt"];

    //设置全局变量
    ACTIVE_PLATFORM = CONF.Default.Platform;

    var req = {
        "platform": pf.toInteger(),
        "catname": menu,
    };
    var request = JSON.stringify(req);

    ROMJSON = view.GetGameList(request);

    if (pf == 0) {
        $(#add_rom).style["display"] = "none";
    } else {
        $(#add_rom).style["display"] = "block";
    }

    createRomList(1); //生成rom列表
    //生成rom数量
    var romCount = view.GetGameCount(request);
    $(#rom_count_num).html = romCount; //初始化游戏数量

    //更新侧边栏平台介绍
    sidebarPlatformDesc(pf);

    //初始化音量控制
    VIDEO_VALUME = CONF.Default.VideoVolume;

    $(body).state.focus = true;

    //如果没有平台，则自动打开平台设置窗口
    if (CONF.PlatformList.length == 0) {
        openPlatformConfig();
    }

}

//初始化配置
function initConfig(isfresh) {

    var confstr = view.InitData("config", isfresh);
    CONF = JSON.parse(confstr);

    //如果初始路径不同，则刷新数据
    if (CONF.RootPath != CONF.Default.RootPath) {
        view.UpdateConfig("root_path", CONF.RootPath);
    }

    //软件名称
    if (System.scanFiles(CONF.Default.SoftName)) { //图片logo
        $(#logo).html = "<img src='" + URL.fromPath(CONF.Default.SoftName) + "' />";
    } else {
        $(#logo).html = CONF.Default.SoftName; //文字logo
    }

    //创建语言
    if (isfresh == "0") {
        createLang();
    }

    initTheme(); //创建主题
    initInterface(); //创建界面设置（背景、透明度）
    initPlatform(); //生成平台列表
    initShortcutList(); //初始化快捷工具列表
    initRomListStyle(); //初始化游戏列表样式
    initBaseFilter(CONF.Default.Platform); //初始化基本信息过滤器
    //生成菜单
    var menujson = view.GetMenuList(CONF.Default.Platform, 0);
    createMenuList(menujson);

    //激活第一个字母索引
    if($(#num_search).select("li:current") == undefined){
        $(#num_search li).state.current = true;
    }

}

//创建界面设置
function initInterface() {
    //初始化选项
    $(#config_font_background).select("li[opt=" + CONF.Default.RomlistFontBackground + "]").attributes.addClass("active");
    $(#config_orders).select("li[opt=" + CONF.Default.RomlistOrders + "]").attributes.addClass("active");
    $(#config_show_name_type).select("li[opt=" + CONF.Default.RomlistNameType + "]").attributes.addClass("active");


    //各平台缩略图类型
    createConfigThumb();

    //各平台缩略图方向
    createConfigThumbDirection();

    //初始化展示图字体大小
    createRomlistFontsize();

    //初始化模块间距
    createRomlistMargin();

    //初始化模块大小
    createRomlistSize();

    //更换鼠标指针
    $(body).style["cursor"] = [url: URL.fromPath(CONF.Default.Cursor)];

    //更新背景图
    if (CONF.Default.BackgroundImage != "") {
        showListBackground("");
    }

    //更新背景循环方式
    var repeat = CONF.Default.BackgroundRepeat;
    if (repeat == "cover") {
        $(#center).style["background-size"] = repeat;
        $(#center).style["background-repeat"] = "no-repeat";
    } else if (repeat == "repeat") {
        $(#center).style["background-repeat"] = repeat;
        $(#center).style["background-size"] = "auto auto";
    } else {
        $(#center).style["background-repeat"] = repeat;
        $(#center).style["background-size"] = "contain";
    }

}

//初始化游戏列表样式
function initRomListStyle() {

    var romlist = $(#romlist);

    //初始化游戏列表样式
    if (CONF.Default.RomlistStyle == "1") {
        $(#switch_romlist).html = "<i.block.minsize></i>";
        romlist.attributes.addClass("romblock");
    } else {
        $(#switch_romlist).html = "<i.list.minsize></i>";
        romlist.attributes.addClass("romlist");
    }

    //初始化rom列表模块大小
    romlist.attributes.addClass("zoom" + CONF.Default.RomlistSize);
    $(#switch_romlist).attributes["value"] = CONF.Default.RomlistStyle;

    //初始化rom列表字体大小
    romlist.attributes.addClass("fontsize" + CONF.Default.RomlistFontSize);

    //初始化列表模块间距
    romlist.attributes.addClass("margin" + CONF.Default.RomlistMargin);

    //是否显示rom列表的标题背景颜色
    romlist.attributes.addClass("bg" + CONF.Default.RomlistFontBackground);

    //字体阴影
    if (CONF.Default.RomlistFontBackground == 0) {
        romlist.attributes.addClass("textShadow");
    }

    //rom列表缩略图方向
    initRomlistThumbDirection();

    //初始化列表列显示项
    var column = CONF.Default.RomlistColumn.split(",");
    var dom = $$(.romlist_column);
    var i = 0;
    for (var obj in dom) {
        obj.checked = column[i] == 1 ? true : false;
        i++;
    }

}

//生成快捷工具列表
function initShortcutList() {
    var li = "";
    for (var obj in CONF.Shortcut) {
        li += "<li path='" + obj.Path + "'>" + obj.Name + "</li>";
    }
    $(#shortcut menu).html = li; //生成dom
}

//创建主题
function initTheme() {

    //默认主题
    self.attributes["theme"] = CONF.Default.Theme;


    //渲染页面主题
    initUiTheme();

    //生成主题列表
    var menu = "";
    for (var themeId in CONF.Theme) {
        menu += "<li id=" + themeId + ">" + CONF.Theme[themeId].Name + "</li>"; //填充菜单
    }
    $(#theme menu).html = menu;

    //设置首页主题样式
    setTheme(CONF.Default.Theme);
}

//初始化过滤器
function initBaseFilter(platform) {

    $(#search_input).value = "";

    //填充数据
    var getjson = view.GetFilter(platform);
    var jsonObj = JSON.parse(getjson);
    for (var type in jsonObj) {
        createBaseFilterOptions(type, jsonObj[type]);
    }


    //是否显示和隐藏
    var column = CONF.Default.RomlistColumn.split(",");
    $(#filter_type).style["display"] = column[1] == 0 ? "none" : "inline-block";
    $(#filter_year).style["display"] = column[2] == 0 ? "none" : "inline-block";
    $(#filter_publisher).style["display"] = column[3] == 0 ? "none" : "inline-block";
    $(#filter_country).style["display"] = column[4] == 0 ? "none" : "inline-block";
    $(#filter_translate).style["display"] = column[5] == 0 ? "none" : "inline-block";
    $(#filter_version).style["display"] = column[6] == 0 ? "none" : "inline-block";
    $(#filter_producer).style["display"] = column[7] == 0 ? "none" : "inline-block";
    $(#filter_score).style["display"] = column[9] == 0 ? "none" : "inline-block";
    $(#filter_complete).style["display"] = column[10] == 0 ? "none" : "inline-block";
}

//初始化过滤器选项
function createBaseFilterOptions(type, data) {

    var dom;
    var title = "";

    switch (type) {
        case "base_type":
            dom = $(#filter_type);
            title = CONF.Lang.BaseType;
            break;
        case "base_year":
            dom = $(#filter_year);
            title = CONF.Lang.BaseYear;
            break;
        case "base_producer":
            dom = $(#filter_producer);
            title = CONF.Lang.BaseProducer;
            break;
        case "base_publisher":
            dom = $(#filter_publisher);
            title = CONF.Lang.BasePublisher;
            break;
        case "base_country":
            dom = $(#filter_country);
            title = CONF.Lang.BaseCountry;
            break;
        case "base_translate":
            dom = $(#filter_translate);
            title = CONF.Lang.BaseTranslate;
            break;
        case "base_version":
            dom = $(#filter_version);
            title = CONF.Lang.BaseVersion;
            break;
        case "score":
            dom = $(#filter_score);
            title = CONF.Lang.Score;
            break;
        case "complete":
            dom = $(#filter_complete);
            title = CONF.Lang.GameCompleteState;
            break;
    }

    dom.value = "";
    dom.options.clear();
    var domData = "<option value=''>" + title + "</option>";

    if (data != undefined && data.length > 0) {
        for (var val in data) {
            var name = "";
            if (type == "complete") {
                //通关状态值是数字
                if (val == 0) name = CONF.Lang.GamePlaying;
                else if (val == 1) name = CONF.Lang.GameComplete;
                else if (val == 2) name = CONF.Lang.GamePlatinumComplete;
            } else {
                name = val;
            }

            domData += "<option value=\"" + val + "\">" + name + "</option>";
        }
    }
    dom.options.html = domData;
    dom.value = "";
}

//初始化展示图类型
function createConfigThumb() {
    var list = "";
    list += "<li opt='optimized'>" + CONF.Lang.OptimizedCache + "</li>";
    list += "<li opt='thumb'>" + CONF.Lang.Thumb + "</li>";
    list += "<li opt='snap'>" + CONF.Lang.Snap + "</li>";
    list += "<li opt='poster'>" + CONF.Lang.Poster + "</li>";
    list += "<li opt='packing'>" + CONF.Lang.Packing + "</li>";
    list += "<li opt='title'>" + CONF.Lang.TitlePic + "</li>";
    list += "<li opt='cassette'>" + CONF.Lang.CassettePic + "</li>";
    list += "<li opt='icon'>" + CONF.Lang.IconPic + "</li>";
    list += "<li opt='gif'>" + CONF.Lang.GifPic + "</li>";
    list += "<li opt='wallpaper'>" + CONF.Lang.WallpaperPic + "</li>";
    list += "<li opt='background'>" + CONF.Lang.BackgroundPic + "</li>";

    var data = "";
    data += "<li opt='all'>" + CONF.Lang.AllPlatform + "<menu>" + list + "</menu></li>";
    for (var pf in CONF.PlatformList) {
        data += "<li opt='" + pf.Id + "'>" + pf.Name + "<menu>" + list + "</menu></li>";
    }

    var menucontext = $(#config_thumb menu);
    menucontext.clear();
    menucontext.append(data);

    //激活
    var configThumb = $(#config_thumb);
    var allDom = configThumb.select("li[opt=all] li[opt=" + CONF.Default.Thumb + "]")
    if (allDom != undefined) {
        allDom.attributes.addClass("active");
    }
    for (var pf in CONF.PlatformList) {
        if (pf.Thumb != "") {
            var dom = configThumb.select("li[opt=" + pf.Id + "] li[opt=" + pf.Thumb + "]");
            if (dom != undefined) {
                dom.attributes.addClass("active");
            }

        }
    }
}

//初始化展示图方向选项列表
function createConfigThumbDirection() {
    var list = "";
    list += "<li opt='0'>" + CONF.Lang.ThumbAuto + "</li>";
    list += "<li opt='1'>" + CONF.Lang.ThumbHorizontal + "</li>";
    list += "<li opt='2'>" + CONF.Lang.ThumbVertical + "</li>";

    var data = "";
    data += "<li opt='all'>" + CONF.Lang.AllPlatform + "<menu>" + list + "</menu></li>";
    for (var pf in CONF.PlatformList) {
        data += "<li opt='" + pf.Id + "'>" + pf.Name + "<menu>" + list + "</menu></li>";
    }

    var menucontext = $(#config_romlist_direction menu);
    menucontext.clear();
    menucontext.append(data);

    //激活
    var configDirection = $(#config_romlist_direction);
    var dom = configDirection.select("li[opt=all] li[opt=" + CONF.Default.RomlistDirection + "]");
    if (dom != undefined) {
        dom.attributes.addClass("active");
    }

    for (var pf in CONF.PlatformList) {
        if (pf.ThumbDirection != "") {
            var dirDom = configDirection.select("li[opt=" + pf.Id + "] li[opt=" + pf.ThumbDirection + "]");
            if (dirDom != undefined) {
                dirDom.attributes.addClass("active");
            }
        }
    }
}

//更新romlist展示图方向
function initRomlistThumbDirection() {

    var romlist = $(#romlist);
    if (romlist.attributes.hasClass("horizontal")) {
        romlist.attributes.removeClass("horizontal");
    }
    if (romlist.attributes.hasClass("vertical")) {
        romlist.attributes.removeClass("vertical");
    }

    var dir = CONF.Default.RomlistDirection;
    if (CONF.Platform[ACTIVE_PLATFORM.toString()] != undefined && CONF.Platform[ACTIVE_PLATFORM.toString()].ThumbDirection != "") {
        dir = CONF.Platform[ACTIVE_PLATFORM.toString()].ThumbDirection;
    }

    if (dir == 1) {
        romlist.attributes.addClass("horizontal");
    } else if (dir == 2) {
        romlist.attributes.addClass("vertical");
    }

}

//初始化展示图字体大小
function createRomlistFontsize() {
    var list = "";
    for (var i = 1; i <= ROOMLIST_SIZE_NUM; i++) {
        list += "<li opt='" + i + "'>" + i + "</li>";
    }

    var data = "";
    data += "<li opt='all'>" + CONF.Lang.AllPlatform + "<menu>" + list + "</menu></li>";
    for (var pf in CONF.PlatformList) {
        data += "<li opt='" + pf.Id + "'>" + pf.Name + "<menu>" + list + "</menu></li>";
    }

    var menucontext = $(#config_title_fontsize menu);
    menucontext.clear();
    menucontext.append(data);

    //激活
    var configFontSize = $(#config_title_fontsize);
    var configRomlistFontSize = configFontSize.select("li[opt=all] li[opt=" + CONF.Default.RomlistFontSize + "]");
    if (configRomlistFontSize != undefined) {
        configRomlistFontSize.attributes.addClass("active");
    }
    for (var pf in CONF.PlatformList) {
        if (pf.ThumbFontSize != "") {
            var fontsizeDom = configFontSize.select("li[opt=" + pf.Id + "] li[opt=" + pf.ThumbFontSize + "]");
            if (fontsizeDom != undefined) {
                fontsizeDom.attributes.addClass("active");
            }
        }
    }
}

//更新romlist展示图字体大小
function initRomlistFontsize() {
    var dom = $(#romlist);
    for (var i = 1; i <= ROOMLIST_SIZE_NUM; i++) {
        var s = 'fontsize' + i.toString();
        if (dom.attributes.hasClass(s)) {
            dom.attributes.removeClass(s);
        }
    }

    var size = CONF.Default.RomlistFontSize;
    if (CONF.Platform[ACTIVE_PLATFORM.toString()] != undefined && CONF.Platform[ACTIVE_PLATFORM.toString()].ThumbFontSize != "") {
        size = CONF.Platform[ACTIVE_PLATFORM.toString()].ThumbFontSize;
    }
    dom.attributes.addClass("fontsize" + size);
}

//初始化展示图模块间距
function createRomlistMargin() {
    var list = "";
    for (var i = 1; i <= ROOMLIST_MARGIN_NUM; i++) {
        list += "<li opt='" + i + "'>" + i + "</li>";
    }

    var data = "";
    data += "<li opt='all'>" + CONF.Lang.AllPlatform + "<menu>" + list + "</menu></li>";
    for (var pf in CONF.PlatformList) {
        data += "<li opt='" + pf.Id + "'>" + pf.Name + "<menu>" + list + "</menu></li>";
    }

    var menucontext = $(#config_romlist_margin menu);
    menucontext.clear();
    menucontext.append(data);

    //激活
    var dom = $(#config_romlist_margin);
    var configRomlistMargin = dom.select("li[opt=" + CONF.Default.RomlistMargin + "]");
    if (configRomlistMargin != undefined) {
        configRomlistMargin.attributes.addClass("active");
    }
    for (var pf in CONF.PlatformList) {
        if (pf.ThumbMargin != "") {
            var marginDom = dom.select("li[opt=" + pf.Id + "] li[opt=" + pf.ThumbMargin + "]");
            if (marginDom != undefined) {
                marginDom.attributes.addClass("active");
            }
        }
    }
}

//更新romlist模块间距
function initRomlistMargin() {
    var dom = $(#romlist);
    for (var i = 1; i <= ROOMLIST_MARGIN_NUM; i++) {
        var s = 'margin' + i.toString();
        if (dom.attributes.hasClass(s)) {
            dom.attributes.removeClass(s);
        }
    }

    var size = CONF.Default.RomlistMargin;
    if (CONF.Platform[ACTIVE_PLATFORM.toString()] != undefined && CONF.Platform[ACTIVE_PLATFORM.toString()].ThumbMargin != "") {
        size = CONF.Platform[ACTIVE_PLATFORM.toString()].ThumbMargin;
    }
    dom.attributes.addClass("margin" + size);
}

//初始化展示图模块大小
function createRomlistSize() {
    var sizeNum = 15;
    var list = '';
    list += '<li opt="1">32px</li>';
    list += '<li opt="2">80px</li>';
    list += '<li opt="3">100px</li>';
    list += '<li opt="4">120px</li>';
    list += '<li opt="5">140px</li>';
    list += '<li opt="6">170px</li>';
    list += '<li opt="7">200px</li>';
    list += '<li opt="8">240px</li>';
    list += '<li opt="9">280px</li>';
    list += '<li opt="10">320px</li>';
    list += '<li opt="11">480px</li>';
    list += '<li opt="12">640px</li>';
    list += '<li opt="13">960px</li>';
    list += '<li opt="14">1280px</li>';
    list += '<li opt="15">1680px</li>';
    list += '<li opt="16">100%</li>';

    var data = "<li opt='all'>" + CONF.Lang.AllPlatform + "<menu>" + list + "</menu></li>";
    for (var pf in CONF.PlatformList) {
        var active = pf.ThumbSize == "" ? "" : pf.ThumbSize;
        data += "<li opt='" + pf.Id + "' class='" + active + "'>" + pf.Name + "<menu>" + list + "</menu></li>";
    }

    var menucontext = $(#config_romlist_size menu);
    menucontext.clear();
    menucontext.append(data);

    //激活
    var dom = $(#config_romlist_size);
    var configRomlistSize = dom.select("li[opt=" + CONF.Default.RomlistSize + "]");
    if (configRomlistSize != undefined) {
        configRomlistSize.attributes.addClass("active");
    }
    for (var pf in CONF.PlatformList) {
        if (pf.ThumbSize != "") {
            var sizeDom = dom.select("li[opt=" + pf.Id + "] li[opt=" + pf.ThumbSize + "]");
            if (sizeDom != undefined) {
                sizeDom.attributes.addClass("active");
            }
        }
    }
}

//更新romlist 模块大小
function initRomlistSize() {
    var dom = $(#romlist);
    for (var i = 1; i <= ROOMLIST_ROOM_SIZE_NUM; i++) {
        var s = 'zoom' + i.toString();
        if (dom.attributes.hasClass(s)) {
            dom.attributes.removeClass(s);
        }
    }

    var size = CONF.Default.RomlistSize;
    if (CONF.Platform[ACTIVE_PLATFORM.toString()] != undefined && CONF.Platform[ACTIVE_PLATFORM.toString()].ThumbSize != "") {
        size = CONF.Platform[ACTIVE_PLATFORM.toString()].ThumbSize;
    }
    dom.attributes.addClass("zoom" + size);
}