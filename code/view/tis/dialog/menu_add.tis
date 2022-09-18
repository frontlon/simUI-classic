view.root.on("ready", function(){
    Init(); //初始化
});

//初始化
function Init(){
   //创建语言
    createLang();
    view.windowCaption = CONF.Lang.AddMenu;

    var platform = view.parameters.platform.toString();

    if(platform > 0 ){
        $(#platform).style["display"] = "none";
        return;
    }

    //生成平台列表
    var platformJson = mainView.GetPlatform();
    var platformList = JSON.parse(platformJson);
    var options = "<option value=''>"+ CONF.Lang.SelectPlatform +"</option>";
    for(var pf in platformList) {
        options += "<option value='"+ pf.Id +"'>"+pf.Name+"</option>";
    }
    $(#platform).html = options;
}
//确定
function confirmDialog(evt){
    var platform = view.parameters.platform;
    if(platform == 0){
        platform = $(#platform).value;
        if(platform == 0 ){
            alert(CONF.Lang.NotSelectPlatform)
            return;
        }
    }

    var menuname = $(#menuname).value;
    if(menuname == ""){
        alert(CONF.Lang.MenuNameCanNotBeEmpty)
        return;
    }

    var virtual = $(#virtual).value.toInteger();

    //创建目录
    mainView.AddMenu(platform,menuname,virtual);

    view.close(menuname); //如果更新成功，则关闭窗口，并把修改后的数据返回
}

