//生成平台列表
function managerInit(){

    //标题
    view.windowCaption = CONF.Lang.RomManager;

    //初始化主题
    initUiTheme();

    //创建语言
    createLang();

    //显示“请选择平台”按钮
    var dd = "<option value='0'>"+CONF.Lang.SelectPlatform +"</option>";
    
    //遍历数据，生成dom
    for(var pf in CONF.PlatformList) {
        dd += "<option value='"+ pf.Id +"'>" + pf.Name + "</option>";
    }
    self.$(#platform_rombase).html = dd; //资料管理
    //self.$(#platform_media).html = dd; //资源管理
    self.$(#platform_sim).html = dd; //模拟器管理
    self.$(#platform_file).html = dd; //文件管理
    self.$(#platform_repeat).html = dd; //rom去重
    self.$(#platform_zombie).html = dd; //无效资源清理
    self.$(#platform_subgame).html = dd; //子游戏
}


//创建菜单选项 - 公共
function createMenuOption(platformId){
    var menujson = mainView.GetMenuList(platformId,0); //读取分配列表
    var menuobj = JSON.parse(menujson);
    var dd = "<option value=''>"+CONF.Lang.AllGames+"</option>";
    var name ="";
    var fixed = "";
    //遍历数据，生成dom
    for(var obj in menuobj) {
        if(obj.Name == "_7b9"){
            name = CONF.Lang.Uncate;
        }else{
            name = obj.Name;
        }

        dd += "<option value=\""+ obj.Name +"\">"+ name + "</option>";
    }
    return dd;
}

//创建分页数据 - 公共
function managerCreatePages(romCount){
    var num = Math.ceil(romCount.toFloat() / CONST_ROM_LIST_PAGE_SIZE);

    var pages = "";
    for(var i=1;i<=num;i++) {
        pages += "<li>"+i+"</li>";
    }
    return pages;
}