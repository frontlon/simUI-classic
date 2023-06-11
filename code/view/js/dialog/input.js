view.root.on("ready", function () {
    try {
        Init();
    } catch (e) {
        alert(e);
    }
});

//初始化方法
function Init() {

    //初始化主题
    initUiTheme();

    //渲染语言
    createLang();

    //定义标题
    view.windowCaption = CONF.Lang.RomOutput;

    //生成平台列表
    var platforms = "<option value='0'>" + CONF.Lang.SelectPlatform + "</option>";
    for (var pf in CONF.PlatformList) {
        platforms += "<option value='" + pf.Id + "'>" + pf.Name + "</option>";
    }
    $(#platform).options.html = platforms; //资料管理

}

function changeType(type) {

    if (type == "") {
        return;
    }

    //先清空选项
    clearOptions();

    $(#platform).style["display"] = "inline-block";
    $(#platform).value = "0";
    $(#menu_h2).style["display"] = "none";
    $(#menu_ul).style["display"] = "none";
    $(#game_h2).style["display"] = "none";
    $(#game_ul).style["display"] = "none";
    $(#search_h2).style["display"] = "none";
    $(#romname_wrapper).style["display"] = "none";
    $(#options_h2).style["display"] = "none";
    $(#options_ul).style["display"] = "none";
}

//切换平台
function changePlatform(platformId) {
    if (platformId == "0") {
        return;
    }

    //先清空选项
    clearOptions();

    $(#options_h2).style["display"] = "block";
    $(#options_ul).style["display"] = "block";

    var type = $(#type).value;
    if (type == "menu") {
        $(#menu_ul).html = createMenuOption(platformId);
        $(#menu_h2).style["display"] = "block";
        $(#menu_ul).style["display"] = "block";
        $(#game_h2).style["display"] = "none";
        $(#game_ul).style["display"] = "none";
        $(#search_h2).style["display"] = "none";
        $(#romname_wrapper).style["display"] = "none";
    } else if (type == "rom") {
        $(#search_h2).style["display"] = "block";
        $(#romname_wrapper).style["display"] = "block";
        $(#game_h2).style["display"] = "block";
        $(#game_ul).style["display"] = "block";
        $(#menu_h2).style["display"] = "none";
        $(#menu_ul).style["display"] = "none";
    }

}

//生成menu菜单
function createMenuOption(platformId) {
    var menujson = mainView.GetMenuList(platformId, 0);
    var menuobj = JSON.parse(menujson);
    var name = "";
    var html = "";

    //遍历数据，生成dom
    for (var obj in menuobj) {
        if (obj.Name == "_7b9") {
            name = CONF.Lang.Uncate;
        } else {
            name = obj.Name;
        }

        html += "<li><label><input type='checkbox' value=" + obj.Name + " />" + name + "</label></li>";
    }
    return html;
}

//模糊搜索rom
function searchRom(keyword) {
    if (keyword == "") {
        $(#name_ul).style["display"] = "none";
        return;
    }

    //生成游戏列表
    var platform = $(#platform).value;
    var request = {
        "platform": platform.toInteger(),
        "keyword": keyword.toString(),
    };
    var romJson = mainView.GetGameList(JSON.stringify(request));
    var romObj = JSON.parse(romJson);
    var html = "";

    for (var obj in romObj) {
        html += "<li opt='" + obj.Id + "'><div>" + obj.Name + "<span.small_info>" + obj.RomPath + "</span></div></li>";
    }

    if (html == "") {
        html += "<li opt='0'><div.empty>没有找到搜索结果</div></li>";
    }

    $(#name_ul).html = html;
    $(#name_ul).style["display"] = "block";
}

function selectGame(id, content) {
    if ($(#game_ul).select("li[opt=" + id + "]") == undefined) {
        var html = "<li opt='" + id + "'>" + content + "<button>删除</button></li>";
        $(#game_ul).append(html);
    }
    $(#name_ul).select("li[opt=" + id + "]").remove();
    $(#name_ul).style["display"] = "block";
}

function clearOptions() {
    $(#game_ul).html = "";
    $(#menu_ul).html = "";
    $(#name_ul).html = "";
    $(#romname).value = "";
}

function outputData() {
    var type = $(#type).value;
    var checks = $$(#options_ul input[type = 'checkbox']);
    var options = [];
    var request = {};
    var menus = [];
    var roms = [];
    for (var opt in checks) {
        if (opt.checked) {
            options.push(opt.value)
        }
    }

    if (type == "") {
        alert("没有选择导出类型");
        return;
    }

    if (options.length == 0) {
        alert("没有选择任何导出选项");
        return;
    }

    var platform = $(#platform).value;
    if (platform == 0) {
        alert("没有选择平台");
        return;
    }

    if (type == "menu") {
        var checks = $$(#menu_ul input[type = 'checkbox']);
        for (var opt in checks) {
            if (opt.checked) {
                menus.push(opt.value)
            }
        }
        if (menus.length == 0) {
            alert("没有选择任何目录");
            return;
        }
    } else if (type == "rom") {
        var checks = $$(#game_ul li);
        for (var opt in checks) {
            roms.push(opt.attributes["opt"].toInteger());
        }
        if (roms.length == 0) {
            alert("没有选择任何rom");
            return;
        }
    }

    //打开文件保存窗口
    let name = CONF.Platform[platform.toString()].Name;
    let defaultExt = "";
    let initialPath = name + ".zip";
    let exts = "*.zip";
    let filter = "Files (" + exts + ")|" + exts + "|All Files (*.*)|*.*";
    let caption = CONF.Lang.SelectFile;
    let selectFile = view.selectFile(#save, filter, defaultExt, initialPath, caption);
    if (selectFile == undefined) {
        return;
    }

    //调用服务端
    request["save"] = URL.toPath(selectFile);
    request["opt"] = type;
    request["options"] = options;
    request["platform"] = platform.toInteger();
    request["menus"] = menus;
    request["roms"] = roms;
    mainView.OutputRom(JSON.stringify(request));
}