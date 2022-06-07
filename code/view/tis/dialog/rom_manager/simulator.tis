
/**
 * 模拟器管理
 */

//切换平台 - 模拟器管理
function managerChangeSimulatorPlatform(platformId){

    if (platformId == 0){
        return;
    }

    //创建菜单
    var dd = createMenuOption(platformId);
    var menulist = self.$(#menu_sim);
    menulist.options.clear();
    if(dd != ""){
        menulist.options.html = dd; //生成dom
    }else{
        menulist.options.html = "<option value=''>"+ CONF.Lang.AllGames +"</option>"; //生成dom
    }
    menulist.value = "";

    //创建模拟器列表
    var sim = "<option value='0'>"+CONF.Lang.SelectSimulator +"</option>";
    for(var pf in CONF.Platform[platformId.toString()].SimList) {
        sim += "<option value='"+ CONF.Platform[platformId.toString()].SimList[pf].Id +"'>"+ CONF.Platform[platformId.toString()].SimList[pf].Name +"</option>"
    }
    $(#select_sim).html = sim;

    //生成默认游戏列表
    var req = {
        "platform" : platformId.toInteger(),
    };
    var request = JSON.stringify(req);
    var romjson = mainView.GetGameList(request);
    $(#romlist_sim tbody).html = "";
    managerCreateSimulator(romjson);

    //生成分页数据
    var romCount = mainView.GetGameCount(request); //必须在获取rom的方法下面
    var pages = managerCreatePages(romCount);
    $(#pages_sim).html = pages;
    if (pages != ""){
        $(#pages_sim).select("li:nth-child(1)").state.current = true;
    }
    if($$(#pages_sim li).length == 1){
        $(#pages_sim).style["display"] = "none";
    }
}

//创建rom模拟器列表 - 模拟器管理
function managerCreateSimulator(romjson){

    var romobj = JSON.parse(romjson);
    var romData = "";
    for(var obj in romobj) {
   
        romData += "<tr rid='"+ obj.Id +"'>";
        var romName = obj.RomPath.replace("\\","/");
        var romPath = romName.split("/");
        var filename = romPath[romPath.length - 1].split(".");
        
        var simConf = {};
        var unzipText = CONF.Lang.No;
        var unzip = 0;
        var bootFile = "";
        var simName = "";
        var simId = 0;
        var cmd = "";
        var lua = "";
        
        if(obj.SimId == 0){
            if(CONF.Platform[obj.Platform.toString()].UseSim != undefined){
                simName = CONF.Platform[obj.Platform.toString()].UseSim.Name;
                simId = CONF.Platform[obj.Platform.toString()].UseSim.Id;
            }
        }else{

            if(CONF.Platform[obj.Platform.toString()].SimList[obj.SimId.toString()] != undefined){
                simName = CONF.Platform[obj.Platform.toString()].SimList[obj.SimId.toString()].Name;
                simId = CONF.Platform[obj.Platform.toString()].SimList[obj.SimId.toString()].Id;

            }else{
                if(CONF.Platform[obj.Platform.toString()].UseSim != undefined){
                    simName = CONF.Platform[obj.Platform.toString()].UseSim.Name;
                    simId = CONF.Platform[obj.Platform.toString()].UseSim.Id;
                }
            }
        }

        if(obj.SimConf != "" && obj.SimConf != "{}"){
            simConf = JSON.parse(obj.SimConf);
            simConf = simConf[simId.toString()];
            if(simConf != undefined){
                unzipText = simConf.Unzip == 1 ? CONF.Lang.Yes : CONF.Lang.No;
                bootFile = simConf.File;
                if(simConf.Cmd != ""){
                    cmd = simConf.Cmd;
                }
                unzip = simConf.Unzip;
                lua = simConf.Lua;
            }
        }else{
            simConf = CONF.Platform[obj.Platform.toString()].UseSim;
            unzipText = simConf.Unzip == 1 ? CONF.Lang.Yes : CONF.Lang.No;
            bootFile = "";
            unzip = simConf.Unzip;
            lua = simConf.Lua;
        }

        romData += "<td><input type='checkbox' value='"+ obj.Id +"' class='check_simulator' id='check_simulator_"+ obj.Id +"' /></td>";
        romData += "<td>"+ filename[0] +"</td>";
        romData += "<td>"+ obj.Name.toHtmlString() +"</td>";
        romData += "<td.simulator_name opt='"+ simId +"'>"+ simName +"</td>";
        romData += "<td><input|text.simulator_cmd value='"+ cmd +"'></td>";
        romData += "<td.simulator_unzip opt='"+ unzip +"'>"+ unzipText +"</td>";
        romData += "<td><input|text.simulator_file value='"+ bootFile +"'></td>";
        romData += "<td.simulator_lua>"+ lua +"</td>";
        romData += "</tr>"
    }
    $(#romlist_sim tbody).html = romData;
}

//切换目录 - 模拟器管理
function managerChangeSimulatorMenu(menuName){

    var platformId = $(#platform_sim).value;

    //生成默认游戏列表
    var req = {
        "platform" : platformId.toInteger(),
        "catname" : menuName,
    };
    var request = JSON.stringify(req);
    var romjson = mainView.GetGameList(request);
    $(#romlist_sim tbody).html = "";
    managerCreateSimulator(romjson);

    var romCount = mainView.GetGameCount(request); //必须在获取rom的方法下面
    var pages = managerCreatePages(romCount);
    $(#pages_sim).html = pages;
    if(pages != ""){
        $(#pages_sim).select("li:nth-child(1)").state.current = true;
    }
    if($$(#pages_sim li).length == 1){
        $(#pages_sim).style["display"] = "none";
    }
}

//点击分页按钮 - 模拟器管理
function managerCreateSimulatorByPages(obj){
    if(obj.state.current == true){
        return;
    }
    var platformId = $(#platform_sim).value;
    var id = obj.parent.id;
    var page = obj.html.toInteger() - 1;

    var req = {
        "platform" : platformId.toInteger(),
        "page" : page,
    };
    var request = JSON.stringify(req);

    var romjson = mainView.GetGameList(request);
    $(#romlist_sim tbody).html = "";
    managerCreateSimulator(romjson);
    obj.state.current = true;
    $(body).scrollTo(0,0, false);
}

//保存模拟器参数 - 模拟器管理
function managerSimulatorSave(obj){
    var tr = obj.parent.parent;
    var romId = tr.attributes["rid"];
    var simId = tr.select(".simulator_name").attributes["opt"];
    var unzip = tr.select(".simulator_unzip").attributes["opt"];
    var cmd = tr.select(".simulator_cmd").value;
    var file = tr.select(".simulator_file").value;
    var lua = tr.select(".simulator_lua").text;
    var data = {
        cmd:cmd,
        unzip:unzip.toString(),
        file:file,
        lua:lua,    
    }
    var datastr = JSON.stringify(data);
    mainView.UpdateRomCmd(romId,simId,datastr);

    obj.attributes.removeClass("nosave");

}

//改变rom模拟器 - 模拟器管理
function managerChangeRomSimulator(obj){
 
    var checks = $$(#romlist_sim .check_simulator:checked);
    if(obj.value == 0){
        return;
    }
    if(checks.length == 0){
    alert(CONF.Lang.NotSelectGames);
        return;
    }

    var arr = [];
    for(var c in checks){
        arr.push(c.value);
    }

    var ids = arr.join(",");
    mainView.SetRomSimulator(ids,obj.value);

    //回显列表数据
    var romlist = JSON.parse(mainView.GetGameListByIds(ids));

    var simConf = "";
    var gamelist = new Object();
    for(var r in romlist){
        if(r.SimConf != "" && r.SimConf != "{}"){
            simConf = JSON.parse(r.SimConf);
            gamelist[r.Id] = simConf;
        }
    }

    var platformId = $(#platform_sim).value;
    var conf = CONF.Platform[platformId.toString()].SimList[obj.value.toString()];

    if(conf != undefined){
        for(var c in checks){
            c.parent.parent.select(".simulator_cmd").value = conf.Cmd;
            if(gamelist[c.value.toInteger()] != undefined && gamelist[c.value.toInteger()][obj.value.toString()] != undefined){
                c.parent.parent.select(".simulator_cmd").text = gamelist[c.value.toInteger()][obj.value.toString()].Cmd;
                c.parent.parent.select(".simulator_unzip").text = gamelist[c.value.toInteger()][obj.value.toString()].Unzip;
                c.parent.parent.select(".simulator_file").text = gamelist[c.value.toInteger()][obj.value.toString()].File;
                c.parent.parent.select(".simulator_lua").text = gamelist[c.value.toInteger()][obj.value.toString()].Lua;
            }else{
                c.parent.parent.select(".simulator_cmd").text = "";
                c.parent.parent.select(".simulator_unzip").text = "";
                c.parent.parent.select(".simulator_file").text = "";
                c.parent.parent.select(".simulator_lua").text = "";
            }
            c.parent.parent.select(".simulator_name").text = conf.Name;
            c.parent.parent.select(".simulator_name").attributes["opt"] = conf.Id;
        }
    }

    $(#select_sim).value = 0;
}

//解压后运行 - 模拟器管理
function managerChangeRomUnzipRunGame(obj){

    var checks = $$(#romlist_sim .check_simulator:checked);
    if(obj.value.toString() == ""){
        return;
    }
    if(checks.length == 0){
    alert(CONF.Lang.NotSelectGames);
        return;
    }

    var text = obj.value == 1 ? CONF.Lang.Yes:CONF.Lang.No;
    for(var c in checks){

    var tr = c.parent.parent;
    var romId = tr.attributes["rid"];
    var simId = tr.select(".simulator_name").attributes["opt"];
    var unzip = obj.value;
    var cmd = tr.select(".simulator_cmd").value;
    var file = tr.select(".simulator_file").value;
    var lua = tr.select(".simulator_lua").text;
    var data = {
        cmd:cmd,
        unzip:unzip.toString(),
        file:file,
        lua:lua,
    }
    var datastr = JSON.stringify(data);

    mainView.UpdateRomCmd(romId,simId,datastr);

    c.parent.parent.select(".simulator_unzip").text = text;
    $(#select_unzip).value = "";

    }

}

//选择文件 - 模拟器管理
function openFileLua(evt) {

    var checks = $$(#romlist_sim .check_simulator:checked);
    if(checks.length == 0){
    alert(CONF.Lang.NotSelectGames);
        return;
    }

    var filter = "Lua Files (*.lua)|*.lua|All Files (*.*)|*.*";

    const defaultExt = "";
    const initialPath = "";
    const caption = CONF.Lang.SelectFile;
    var url = view.selectFile(#open, filter, defaultExt, initialPath, caption );
    if(url){
        url = URL.toPath(url);
        url = url.split("\/").join(SEPARATOR);
        url = url.split(CONF.RootPath.toString()).join("");
        evt.html = url;
        for(var c in checks){
            var tr = c.parent.parent;
            var romId = tr.attributes["rid"];
            var simId = tr.select(".simulator_name").attributes["opt"];
            var unzip = tr.select(".simulator_unzip").attributes["opt"];
            var cmd = tr.select(".simulator_cmd").value;
            var file = tr.select(".simulator_file").value;
            var lua = url;
            var data = {
                cmd:cmd,
                unzip:unzip.toString(),
                file:file,
                lua:lua,
            }
            var datastr = JSON.stringify(data);
            mainView.UpdateRomCmd(romId,simId,datastr);
            c.parent.parent.select(".simulator_lua").text = url;
        }
    }
}