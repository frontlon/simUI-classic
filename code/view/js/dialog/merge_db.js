var INPUT_DATA = [];

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
    view.windowCaption = CONF.Lang.SoftMerge;
}

function openFile(evt){
    const defaultExt = "";
    const initialPath = "";
    const filter = evt.attributes["filter"];
    const caption = evt.attributes["caption"];
    var url = view.selectFile(#open, filter, defaultExt, initialPath, caption );
    var out = self.select("#"+evt.attributes["for"]);  
    if(!url){
        return;
    }

    $(#content).html = "";

    url = URL.toPath(url);

    //读取数据
    var result = mainView.GetMergeDbData(url);
    if(result == undefined){
        return;
    }

    out.value = url;

    INPUT_DATA = JSON.parse(result);
    var platformDom = $(#platform_select);
    var platformHtml = "";
    for (var k in INPUT_DATA) {
        platformHtml +='<li><label><input type="checkbox" class="platform_checkbox" opt="'+ k.Platform.Id +'" value="1" />'+ k.Platform.Name +'</label></li>';
    }
    platformDom.html = platformHtml; 
}

function activePlatform(evt){
    if(evt.value == undefined){
        //隐藏
        evt.parent.parent.attributes.removeClass("active");
        var item = $(#content).select("#item_"+ evt.attributes["opt"]);
        if(item != undefined){
            item.remove();
        }
        return;
    }

    //显示
    evt.parent.parent.attributes.addClass("active");

    //加载页面DOM
    var contentDom = $(#content);
    var itemHtml = "";
    for (var k in INPUT_DATA) {

        if(evt.attributes["opt"] != k.Platform.Id){
            continue;
        }
        
        itemHtml += '<div class="item" id="item_'+ k.Platform.Id +'">';
        itemHtml += '<h2>'+ k.Platform.Name +' ('+ CONF.Lang.GameCount +':'+ k.RomCount +')</h2>';
        //模拟器
        itemHtml += '<h3>'+ CONF.Lang.SoftMergeSelectSimulator +'</h3>';
        itemHtml += '<ul class="sim_select">';
        if(k.Simulators.length == 0){
            itemHtml += '<li>'+ CONF.Lang.PlatformNotSimulator +'</li>';
        }else{
            for (var sim in k.Simulators) {
                itemHtml += '<li class="active"><label><input type="checkbox" class="simulator_checkbox" opt="'+ sim.Id +'" pf="'+ k.Platform.Id +'" value="1" checked />'+ sim.Name +'</label></li>';
            }
        }
        itemHtml += '</ul>';
        //目录检测
        itemHtml += ' <h3>'+ CONF.Lang.ResFileCheck +'</h3>';
        itemHtml += '<table class="res_check">';
        itemHtml += '<thead><tr><th>'+ CONF.Lang.ResType +'</th><th>'+ CONF.Lang.InputTo +'</th><th>'+ CONF.Lang.CheckResult +'</th></tr></thead>';
        itemHtml += '<tbody>';

        for (var pName in k.FolderCheck) {
            itemHtml += '<tr>';
            itemHtml += '<td>'+ CONF.Lang[pName] +'</td>';
            itemHtml += '<td>'+ k.FolderCheck[pName].Path +'</td>';
            itemHtml += '<td class="'+ k.FolderCheck[pName].Status +'">'+ k.FolderCheck[pName].Desc +'</td>';
            itemHtml += '</tr>';
        }
        itemHtml += '</tbody>';
        itemHtml += '</table>';
        itemHtml += '</div>';
    }
    contentDom.append(itemHtml);
}

function activeSimulator(evt){
    if(evt.value == undefined){
        //隐藏
        evt.parent.parent.attributes.removeClass("active");
    }else{
        //显示
        evt.parent.parent.attributes.addClass("active");
    }
}

//开始合并数据
function startMergeDb(){

    var dbFile = $(#db_file).value;
    if(dbFile == ""){
        alert(CONF.Lang.NotSelectDbFile);
        return;
    }

    var platformIds = [];
    var pChecks = $$(.platform_checkbox);
    for(var c in pChecks){
        if(c.checked == true){
            platformIds.push(c.attributes["opt"]);
        }
    }

    if(platformIds.length == 0){
        alert(CONF.Lang.NotSelectPlatform);
        return;
    }

    var simulators = [];
    var sChecks = $$(.simulator_checkbox);
    for(var c in sChecks){
        if(c.checked == true){
            simulators.push(c.attributes["opt"]);
        }
    }

    //合并数据
    var pfstr = JSON.stringify(platformIds);
    var simstr = JSON.stringify(simulators);
    mainView.MergeDB(dbFile,pfstr,simstr);
    view.close();
}
