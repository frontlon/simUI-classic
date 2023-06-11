
view.root.on("ready", function(){

    //菜单拖拽排序
    DragDrop{
        what      : "#menulist>dd:not(.fixed)",
        where     : "#menulist",
        container : "#menulist",        
        notBefore : "#menulist>dd.fixed",
        dropped : function(){
            if(ACTIVE_PLATFORM == 0){
                return;
            }
            var lis = $$(#menulist > dd);
            var name = "";
            var newObj = {};
            for(var li in lis) {
                name = li.attributes["opt"];
                newObj[name] = li.index;
            }
            var datastr = JSON.stringify(newObj);
            view.UpdateMenuSort(ACTIVE_PLATFORM.toString(),datastr);
        }
    }

});

//生成菜单列表
function createMenuList(menujson){
    var menuobj = JSON.parse(menujson);
    var menulist = self.$(#menulist);
    var dd = "";
    var name ="";
    var fixed = "";
    menulist.clear();
    dd += "<dt>"+ CONF.Lang.Cate +"</dt>";
    dd += "<dd opt='' class='menuitem fixed'>"+ CONF.Lang.AllGames +"</dd>";
    dd += "<dd opt='favorite' class='menuitem fixed'>"+ CONF.Lang.Favorite +"</dd>";

    //遍历数据，生成dom
    for(var obj in menuobj) {

        if(obj.Name == "_7b9"){
            name = CONF.Lang.Uncate;
            fixed = "fixed";
        }else{
            name = obj.Name;
            fixed = "";
        }

        dd += "<dd opt=\""+  obj.Name +"\" class='menuitem "+ fixed +"'>"+ name + "</dd>";
    }

    menulist.html = dd; //生成dom

    //设置全局变量
    ACTIVE_MENU = CONF.Default.Menu;

   //设置默认激活按钮
    var active =  $(#menulist).select("dd[opt='"+ ACTIVE_MENU +"']");
    if (active != undefined){
        active.state.current = true;
    }else{
        //没有找到默认激活的按钮，则读取全部菜单
        active =  $(#menulist).select("dd[opt='']");
        if (active != undefined){
            active.state.current = true;
        }
    }

}


//目录单击
function changeMenu(obj){

    //设置全局变量
    ACTIVE_MENU = obj.attributes["opt"];

    //重置滚动翻页数据
    resetScroll();

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

    var req = {
        "platform" : ACTIVE_PLATFORM.toInteger(),
        "catname" : ACTIVE_MENU,
        "page" : SCROLL_PAGE,
        "baseType" : filter_type,
        "basePublisher" : filter_publisher,
        "baseYear" : filter_year,
        "baseCountry" : filter_country,
        "baseTranslate" : filter_translate,
        "baseVersion" : filter_version,
        "baseProducer" : filter_producer,
        "score" : filter_score,
        "complete" : filter_complete,
    };
    var request = JSON.stringify(req);

    ROMJSON = view.GetGameList(request);
    createRomList(1);
    
    //统计游戏数量
    var romCount = view.GetGameCount(request); //必须在获取rom的方法下面
    $(#rom_count_num).html = romCount; //初始化在线人数

    //改变样式
    obj.state.current = true;

    //激活第一个字母索引
    $(#num_search li).state.current = true;

    //更新配置文件
    view.UpdateConfig("menu",ACTIVE_MENU);
}

//点击隐藏目录
function changeHideMenu(obj){

    //设置全局变量
    ACTIVE_MENU = "hide";

    //重置滚动翻页数据
    resetScroll();

    var req = {
        "platform" : ACTIVE_PLATFORM.toInteger(),
        "catname" : ACTIVE_MENU,
    };
    var request = JSON.stringify(req);

    ROMJSON = view.GetGameList(request);
    createRomList(1);
    
    //统计游戏数量
    var romCount = view.GetGameCount(request); //必须在获取rom的方法下面
    $(#rom_count_num).html = romCount; //初始化在线人数

    //激活第一个字母索引
    $(#num_search li).state.current = true;

}

//rom搜索功能 - 过滤器搜索
function search(type=1){

    var keyword = "";
    if(type == 1){
        keyword = self.$(#search_input).value;
    }else{
        keyword = self.$(#search_box_input).value;
    }

    //重置滚动翻页数据
    resetScroll();

    //读取搜索条件
    var menu = $(#menulist).select("dd:current").attributes["opt"];
    var filter_type = $(#filter_type).value === undefined ? "" : $(#filter_type).value;
    var filter_publisher = $(#filter_publisher).value === undefined ? "" : $(#filter_publisher).value;
    var filter_year = $(#filter_year).value === undefined ? "" : $(#filter_year).value;
    var filter_country = $(#filter_country).value === undefined ? "" : $(#filter_country).value;
    var filter_translate = $(#filter_translate).value === undefined ? "" : $(#filter_translate).value;
    var filter_version = $(#filter_version).value === undefined ? "" : $(#filter_version).value;
    var filter_producer = $(#filter_producer).value === undefined ? "" : $(#filter_producer).value;
    var filter_score = $(#filter_score).value === undefined ? "" : $(#filter_score).value;
    var filter_complete = $(#filter_complete).value === undefined ? "" : $(#filter_complete).value;

    //生成游戏列表
    var req = {
        "platform" : ACTIVE_PLATFORM.toInteger(),
        "catname" : menu.toString(),
        "keyword" : keyword.toString(),
        "baseType" : filter_type.toString(),
        "basePublisher" : filter_publisher.toString(),
        "baseYear" : filter_year.toString(),
        "baseCountry" : filter_country.toString(),
        "baseTranslate" : filter_translate.toString(),
        "baseVersion" : filter_version.toString(),
        "baseProducer" : filter_producer.toString(),
        "score" : filter_score.toString(),
        "complete" : filter_complete.toString(),
    };
    var request = JSON.stringify(req);

    var romCount = view.GetGameCount(request);
    $(#rom_count_num).html = romCount; //初始化游戏数量

    ROMJSON = view.GetGameList(request);
    createRomList(1);
}

//菜单向下索引
function nextMenu(){

    //如果界面被隐藏，则禁用此功能
    if ($(#left_menu).style["display"] == "none"){
        return;
    }

    var next = $(#menulist).select("dd:current").index + 2;
    var count = $$(#menulist dd).length;
    var nextDom;
    if (next-2 >= count){
        nextDom = $(#menulist).select("dd:nth-child(0)");
    }else{
        nextDom = $(#menulist).select("dd:nth-child("+ next +")");
    }
    if(nextDom != undefined){
        changeMenu(nextDom);
    }
}

//生成菜单列表
function loadMenuList(menujson){
    var menuobj = JSON.parse(menujson);
    var menulist = self.$(#menulist);
    var dd = "";
    var name ="";
    var fixed = "";

    //遍历数据，生成dom
    for(var obj in menuobj) {
        if(obj.Name == "_7b9"){
            continue;
        }

        name = obj.Name;
        fixed = "";
        dd += "<dd opt='"+  obj.Name +"' class='menuitem "+ fixed +"'>"+ name + "</dd>";
    }

    menulist.append(dd); //生成dom

}

//添加菜单
function addMenu(evt){
    var platform = $(#platform_ul).select("li:current").attributes["platform"];
    var data = view.dialog({
        url:self.url(ROOTPATH + "menu_add.html"),
        width:self.toPixels(500dip),
        height:self.toPixels(210dip),
        parameters: {
            platform:platform.toString();
        };
    });

    if(data == undefined){
        return;
    }
    view.CreateRomCache(platform);    
}