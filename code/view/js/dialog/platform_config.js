var pfId;
view.root.on("ready", function () {

    //初始化方法
    try {
        Init();
    } catch (e) {
        alert(e);
    }
});

//初始化方法
function Init() {
    //渲染语言
    createLang();

    //定义标题
    view.windowCaption = CONF.Lang.PlatformConfig;

    //生成平台列表
    var pfdata = "";

    var liststr = mainView.GetPlatform();
    var lists = JSON.parse(liststr);

    var platform_list = $(#platform_list);

    //var outdata = "<option value='0'>"+ CONF.Lang.SelectPlatform +"</option>";

    for (var pf in lists) {
        pfdata += '<li value="' + pf.Id + '" title="' + pf.Name + '">' + pf.Name + '</li>';
    }
    platform_list.html = pfdata;
    //$(#out_platform).options.html = outdata;
    //$(#out_platform).value = 0;

    $(#platform_second).style["display"] = "none";

}

//平台列表点击
function GetPlatformItem(evt) {


    //显示选项卡
    $(#platform_second).style["display"] = "block";


    pfId = evt.attributes["value"].toString();

    //激活按钮改变样式
    evt.state.current = true;

    //读取平台信息
    var info = JSON.parse(mainView.GetPlatformById(pfId));
    $(#platform_name).value = info.Name;
    $(#platform_tag).value = info.Tag;
    $(#platform_ico).value = info.Icon;
    $(#platform_exts).value = info.RomExts;
    $(#platform_doc).value = info.DocPath;
    $(#platform_strategy).value = info.StrategyPath;
    $(#platform_rom).value = info.RomPath;
    $(#platform_thumb).value = info.ThumbPath;
    $(#platform_snap).value = info.SnapPath;
    $(#platform_poster).value = info.PosterPath;
    $(#platform_packing).value = info.PackingPath;
    $(#platform_title).value = info.TitlePath;
    $(#platform_cassette).value = info.CassettePath;
    $(#platform_icon).value = info.IconPath;
    $(#platform_gif).value = info.GifPath;
    $(#platform_background).value = info.BackgroundPath;
    $(#platform_wallpaper).value = info.WallpaperPath;
    $(#platform_video).value = info.VideoPath;
    $(#platform_files).value = info.FilesPath;
    $(#platform_audio).value = info.AudioPath;
    $(#platform_rombase).value = info.Rombase;
    $(#platform_optimized).value = info.OptimizedPath;

    //如果没有扩展名，则自动填充大量扩展名
    if ($(#platform_exts).value.trim() == "") {
        var exts = ".20,.40,.48,.58,.60,.78,.2hd,.2mg,.32x,.3ds,.3dsx,.68k,.7z,.7zip,.88d,.98d,.a0,.a26,.a52,.a78,.abs,.adf,.adz,.agb,.app,.arc,.atr,.atx,.aus,.axf,.b0,.b5t,.b6t,.bat,.bin,.bml,.bps,.bs,.bsx,.bwt,.car,.cas,.cbn,.ccd,.cci,.cdi,.cdm,.cdt,.cfg,.cgb,.ch8,.chd,.ciso,.cmd,.cof,.col,.com,.conf,.crt,.cso,.csw,.cue,.cxi,.d13,.d1m,.d2m,.d4m,.d64,.d6z,.d71,.d7z,.d80,.d81,.d82,.d88,.d8z,.d98,.dat,.dci,.dcm,.dff,.dim,.dmg,.dms,.do,.dol,.dsk,.dup,.dx2,.elf,.eur,.exe,.fd,.fdd,.fdi,.fds,.fig,.fm2,.fs-uae,.g41,.g4z,.g64,.g6z,.gb,.gba,.gbc,.gbs,.gbz,.gcm,.gcz,.gd3,.gd7,.gdi,.gen,.gg,.gz,.hdd,.hdf,.hdi,.hdm,.hdn,.hdz,.img,.int,.ipf,.ips,.iso,.isz,.j64,.jag,.jap,.jma,.k7,.kcr,.ldb,.lha,.lnx,.lst,.lzx,.m3u,.m5,.m7,.mb,.md,.mdf,.mds,.mdx,.mgd,.mgh,.mgw,.msa,.mv,.mx1,.mx2,.n64,.nca,.ndd,.nds,.nes,.nez,.nfse,.ngc,.ngp,.ngpc,.nhd,.npc,.nrg,.nro,.nsf,.nsfe,.nso,.nsp,.o,.obx,.p,.p00,.p64,.pal,.pbp,.pc2,.pce,.pdi,.po,.prg,.pro,.prof,.prx,.psexe,.pzx,.qd,.rar,.raw,.ri,.rom,.rpx,.rvz,.rzx,.sap,.sc,.scl,.scp,.sfc,.sg,.sgb,.sgg,.sgx,.sk,.smc,.smd,.sms,.sna,.st,.stx,.swc,.swf,.t64,.t81,.tap,.tar,.tfd,.tgc,.thd,.toc,.trd,.tzx,.u1,.uae,.ufo,.unf,.unif,.ups,.usa,.uze,.v64,.vb,.vboy,.vec,.vms,.voc,.vpk,.vsf,.wad,.wav,.wbfs,.wia,.woz,.ws,.wsc,.wud,.wux,.x64,.x6z,.xbe,.xci,.xdf,.xex,.xfd,.xml,.z64,.z80,.zip,.path";
        $(#platform_exts).value = exts;
    }

    //模拟器数据
    var platform_sim = $(#platform_sim);
    platform_sim.options.clear();

    //读取模拟器列表
    var sims = JSON.parse(mainView.GetSimulatorByPlatform(pfId));

    platform_sim.options.$append(<option value=0 > { CONF.Lang.EditSimulator }</option >);
    for (var sim in sims) {
        platform_sim.options.$append(<option value={sim.Id}>{sim.Name}</option>);
    }
    platform_sim.value = 0;

    //激活提交按钮
    $(#openPlatformFolder).state.disabled = false;
    $(#platform_submit).state.disabled = false;
    $(#add_sim).state.disabled = false;
    $(#platform_sim).state.disabled = false;
    $(#platform_name).state.disabled = false;
    $(#platform_tag).state.disabled = false;
    $(#platform_ico).state.disabled = false;
    $(#platform_exts).state.disabled = false;

    var checkEmpty = info.DocPath + info.StrategyPath + info.RomPath + info.ThumbPath + info.SnapPath + info.PosterPath + info.PackingPath + info.TitlePath + info.BackgroundPath + info.VideoPath + info.Rombase;

    //如果所有目录都为空，说明是新增的，则全部禁用，强制用户指定根目录
    if (checkEmpty.trim() != "") {
        $(#platform_doc).state.disabled = false;
        $(#platform_strategy).state.disabled = false;
        $(#platform_rom).state.disabled = false;
        $(#platform_thumb).state.disabled = false;
        $(#platform_snap).state.disabled = false;
        $(#platform_poster).state.disabled = false;
        $(#platform_packing).state.disabled = false;
        $(#platform_title).state.disabled = false;
        $(#platform_cassette).state.disabled = false;
        $(#platform_icon).state.disabled = false;
        $(#platform_gif).state.disabled = false;
        $(#platform_background).state.disabled = false;
        $(#platform_wallpaper).state.disabled = false;
        $(#platform_video).state.disabled = false;
        $(#platform_files).state.disabled = false;
        $(#platform_audio).state.disabled = false;
        $(#platform_rombase).state.disabled = false;
        $(#platform_optimized).state.disabled = false;
        var folders = $$(.openfolder);
        var files = $$(.openfile);
        for (var f in folders) { f.state.disabled = false; }
        for (var f in files) { f.state.disabled = false; }

    } else {
        $(#platform_doc).state.disabled = true;
        $(#platform_strategy).state.disabled = true;
        $(#platform_rom).state.disabled = true;
        $(#platform_thumb).state.disabled = true;
        $(#platform_snap).state.disabled = true;
        $(#platform_poster).state.disabled = true;
        $(#platform_packing).state.disabled = true;
        $(#platform_title).state.disabled = true;
        $(#platform_cassette).state.disabled = true;
        $(#platform_icon).state.disabled = true;
        $(#platform_gif).state.disabled = true;
        $(#platform_background).state.disabled = true;
        $(#platform_wallpaper).state.disabled = true;
        $(#platform_video).state.disabled = true;
        $(#platform_files).state.disabled = true;
        $(#platform_audio).state.disabled = true;
        $(#platform_rombase).state.disabled = true;
        $(#platform_optimized).state.disabled = true;

        var folders = $$(.openfolder);
        var files = $$(.openfile);
        for (var f in folders) { f.state.disabled = true; }
        for (var f in files) { f.state.disabled = true; }
    }

    $(#platform_desc).value = info.Desc;
}


//更新平台信息
function platformSubmit(evt) {
    var data = {
        id: $(#platform_list).select("li:current").attributes["value"].toString();
        name: $(#platform_name).value,
        tag: $(#platform_tag).value,
        ico: $(#platform_ico).value,
        exts: $(#platform_exts).value,
        rom: $(#platform_rom).value,
        thumb: $(#platform_thumb).value,
        snap: $(#platform_snap).value,
        poster: $(#platform_poster).value,
        packing: $(#platform_packing).value,
        title: $(#platform_title).value,
        cassette: $(#platform_cassette).value,
        icon: $(#platform_icon).value,
        gif: $(#platform_gif).value,
        background: $(#platform_background).value,
        wallpaper: $(#platform_wallpaper).value,
        video: $(#platform_video).value,
        strategy: $(#platform_strategy).value,
        doc: $(#platform_doc).value,
        files: $(#platform_files).value,
        audio: $(#platform_audio).value,
        rombase: $(#platform_rombase).value,
        optimized: $(#platform_optimized).value,
    };
    if (data.id == "") { alert(CONF.Lang.SelectPlatform); return true; }
    if (data.name == "") { alert(CONF.Lang.PlatformNameCanNotBeEmpty); return true; }
    if (data.exts == "") { alert(CONF.Lang.RomTypeCanNotBeEmpty); return true; }
    if (data.rom == "") { alert(CONF.Lang.RomMenuCanNotBeEmpty); return true; }
    //检查并补全扩展名
    data.exts = completeExt(data.exts);
    //更新补全的扩展名到文本框中
    $(#platform_exts).value = data.exts;
    //更新名称到列表中
    $(#platform_list).select("li:current").html = data.name;
    //更新平台信息
    var result = mainView.UpdatePlatform(JSON.stringify(data));
    mainView.CreatePlatform(data.id); //创建资源目录及资料文件

    if (result.toString() == "1") {
        alert(CONF.Lang.UpdateSuccess);
    }

}


//添加模拟器
function addSim(evt) {

    var res = view.dialog({
        url: self.url(ROOTPATH + "edit_sim.html"),
        width: self.toPixels(460dip),
        height: self.toPixels(320dip),
        parameters: {
            id: 0,
            platform: pfId,
        }
    })

    //更新option选项
    if (res != undefined && res != "") {
        var sim = JSON.parse(res);
        $(#platform_sim).options.$append(<option value={sim.Id}>{sim.Name}</option>);
    }

};

//修改模拟器模拟器
function editSim(evt) {
    if (evt.value == 0) { return true; }
    var res = view.dialog({
        url: self.url(ROOTPATH + "edit_sim.html"),
        width: self.toPixels(460dip),
        height: self.toPixels(320dip),
        parameters: {
            id: evt.value,
            platform: pfId,
        };
    });

    //将选择还原回【选择模拟器】项目，防止项目被选定
    evt.value = 0;


    //更新option选项
    if (res != undefined && res != "") {
        var sim = JSON.parse(res);
        if (sim.Opt == undefined) {
            //修改
            $(#platform_sim).select("option[value=" + sim.Id + "]").html = sim.Name;
        } else {
            //删除
            $(#platform_sim).select("option[value=" + sim.Id + "]").remove();
        }
    }
}

//选择文件夹
function openFolder(evt) {
    var url = view.selectFolder(evt.attributes["caption"]);
    var out = self.select("#" + evt.attributes["for"]);
    if (url) {
        url = URL.toPath(url);
        url = url.split("\/").join(SEPARATOR);
        url = url.split(CONF.RootPath.toString()).join("");
        out.value = url;
    }
}

//选择文件
function openFile(evt) {
    const defaultExt = "";
    const initialPath = "";
    const filter = evt.attributes["filter"];
    const caption = evt.attributes["caption"];
    var url = view.selectFile(#open, filter, defaultExt, initialPath, caption);
    var out = self.select("#" + evt.attributes["for"]);
    if (url) {
        url = URL.toPath(url);
        url = url.split("\/").join(SEPARATOR);
        url = url.split(CONF.RootPath.toString()).join("");
        out.value = url;
        out.focus = true;
        out.focus = false;
    }
}

//选择平台文件夹
function openPlatformFolder(evt) {
    var url = view.selectFolder(evt.attributes["caption"]);
    if (url) {
        url = URL.toPath(url);
        url = url.split("\/").join(SEPARATOR);
        url = url.split(CONF.RootPath.toString()).join("");

        //资料文件名
        var rombase = $(#platform_name).value == "" ? "rombase" : $(#platform_name).value;

        $(#platform_doc).value = url + SEPARATOR + "docs";
        $(#platform_strategy).value = url + SEPARATOR + "strategies";
        $(#platform_rom).value = url + SEPARATOR + "roms";
        $(#platform_thumb).value = url + SEPARATOR + "thumbs";
        $(#platform_snap).value = url + SEPARATOR + "snaps";
        $(#platform_poster).value = url + SEPARATOR + "poster";
        $(#platform_cassette).value = url + SEPARATOR + "cassette";
        $(#platform_icon).value = url + SEPARATOR + "icon";
        $(#platform_gif).value = url + SEPARATOR + "gif";
        $(#platform_packing).value = url + SEPARATOR + "packing";
        $(#platform_title).value = url + SEPARATOR + "title";
        $(#platform_background).value = url + SEPARATOR + "background";
        $(#platform_wallpaper).value = url + SEPARATOR + "wallpaper";
        $(#platform_video).value = url + SEPARATOR + "video";
        $(#platform_files).value = url + SEPARATOR + "files";
        $(#platform_audio).value = url + SEPARATOR + "audio";
        $(#platform_rombase).value = url + SEPARATOR + rombase.toLowerCase() + ".csv";
        $(#platform_optimized).value = url + SEPARATOR + "optimized";

        $(#platform_doc).state.disabled = false;
        $(#platform_strategy).state.disabled = false;
        $(#platform_rom).state.disabled = false;
        $(#platform_thumb).state.disabled = false;
        $(#platform_snap).state.disabled = false;
        $(#platform_poster).state.disabled = false;
        $(#platform_packing).state.disabled = false;
        $(#platform_title).state.disabled = false;
        $(#platform_cassette).state.disabled = false;
        $(#platform_icon).state.disabled = false;
        $(#platform_gif).state.disabled = false;
        $(#platform_background).state.disabled = false;
        $(#platform_wallpaper).state.disabled = false;
        $(#platform_video).state.disabled = false;
        $(#platform_files).state.disabled = false;
        $(#platform_audio).state.disabled = false;
        $(#platform_rombase).state.disabled = false;
        $(#platform_optimized).state.disabled = false;
        var folders = $$(.openfolder);
        var files = $$(.openfile);
        for (var f in folders) { f.state.disabled = false; }
        for (var f in files) { f.state.disabled = false; }

    }

}

//添加平台
function platformAdd(evt) {
    var name = view.dialog({
        url: self.url(ROOTPATH + "add_platform.html"),
        width: self.toPixels(300dip),
        height: self.toPixels(180dip),
        parameters: {};
    });

    if (name == "" || name == undefined) {
        return;
    }

    //开始添加平台
    var insertId = mainView.AddPlatform(name);

    if (insertId != "0") {
        var newopt = "<li value='" + insertId + "'>" + name + "</li>";
        $(#platform_list).append(newopt);
        alert(CONF.Lang.AddSuccess);
    }
}

//删除平台
function platformDel(evt) {

    if ($(#platform_list li:current) == undefined) {
        return;
    }

    var pfId = $(#platform_list li:current).attributes["value"];
    if (pfId == undefined) {
        alert(CONF.Lang.NotSelectPlatform);
        return true;
    }

    //确认窗口
    var result = confirm(CONF.Lang.DeletePlatformConfirm, CONF.Lang.DeletePlatform);
    if (result != "yes") {
        return true;
    }

    //删除平台
    result = mainView.DelPlatform(pfId);
    if (result == "1") {
        $(#platform_list li:current).remove();
        alert(CONF.Lang.DeleteSuccess);
    }
}

//更新平台介绍
function updatePlatformDesc() {
    var desc = $(#platform_desc).value;

    if ($(#platform_list).select("li:current") == undefined) {
        alert(CONF.Lang.NotSelectPlatform);
        return;
    }

    var platform = $(#platform_list).select("li:current").value;

    mainView.UpdatePlatformFieldById(platform, "desc", desc);
    alert(CONF.Lang.UpdateSuccess);
}