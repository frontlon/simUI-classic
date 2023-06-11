var menu;
var pfId;
view.root.on("ready", function(){
    Init(); //初始化
});

//初始化
function Init(){
    
    view.windowCaption = CONF.Lang.AddSlnk;


    //创建语言
    createLang();

    //初始化主题
    initUiTheme();

    menu = view.parameters.menu.toString();
    pfId = view.parameters.platform.toString();
    
}

//切换游戏类别选项卡
function changeGameType(type){
    $(#select_text).html = CONF.Lang.SelectFile;
    if(type == "ps3"){
        $(#select_text).html = CONF.Lang.SelectFolder;
    }
    $(#file_path).value = "";
}


//选择文件
function openFile(evt){

    const defaultExt = "";
    const initialPath = "";
    const filter = evt.attributes["filter"];
    const caption = evt.attributes["caption"];
    var urls = "";
    if($(#game_type).value == "ps3"){
        urls = view.selectFolder(caption);
    }else{
        urls = view.selectFile(#open-multiple, filter, defaultExt, initialPath, caption );
    }

    if(urls == undefined){
        return;
    }

    var paths = [];
    var out = self.select("#"+evt.attributes["for"]);

    if(typeof(urls) == "string"){
        urls = URL.toPath(urls);
        urls = urls.split("\/").join(SEPARATOR);
        urls = urls.split(CONF.RootPath.toString()).join("");
        paths[0] = urls;
    }else{
        for(var url in urls){
            url = URL.toPath(url);
            url = url.split("\/").join(SEPARATOR);
            url = url.split(CONF.RootPath.toString()).join("");
            paths.push(url);
        }
    }
    out.value = paths.join("|");
}

function submitSlnk(){
    var filePath = $(#file_path).value;
    if (filePath == ""){
        alert(CONF.Lang.NotSelectGames);
        return;
    }

    var gameType = $(#game_type).value;
    var cmd = $(#cmd).value;
    mainView.AddSlnkGame(pfId,menu,gameType,cmd,filePath);
    alert(CONF.Lang.AddSuccess);
    view.close("ok");
}
