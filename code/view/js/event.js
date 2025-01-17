﻿
view.root.on("ready", function () {

    //禁用系统默认滚动
    self.on("keydown", "#center_content", function (evt) {
        return true;
    });


    //窗口获取焦点
    self.on("focusin", function () {
    });

    //窗口失去焦点
    self.on("focusout", function () {
        videoPause();
    });


    /**
    * 窗口
    **/

    //调整窗口大小
    view.on("size", function () {
        windowSize();
    });

    //视图窗口状态改变时
    view.on("statechange", function () {
        windowStateChange(this);
    });

    //调整平台栏宽度
    self.on("mouseup", "#platform_splitter", function (evt) {
        var width = $(#left_platform).box(#width);
        width = Math.round(width * 0.8).toInteger();
        view.UpdateConfig("panel_platform_width", width);
    });

    //调整菜单栏宽度
    self.on("mouseup", "#menu_splitter", function (evt) {
        var width = $(#left_menu).box(#width);
        width = Math.round(width * 0.8).toInteger();
        view.UpdateConfig("panel_menu_width", width);
    });
    //调整侧边栏宽度
    self.on("mouseup", "#right_splitter", function (evt) {
        sidebarSize();
    });

    /**
     * 导航栏
     **/


    //开启关闭平台栏
    self.on("click", "#close_platform", function (evt) {
        togglePlatform();
    });

    //开启关闭菜单栏
    self.on("click", "#close_menu", function (evt) {
        toggleMenu();
    });

    //开启关闭侧边栏
    self.on("click", "#close_right", function (evt) {
        toggleSidebar();
    });

    //切换rom列表样式
    self.on("click", "#switch_romlist", function (evt) {
        switchRomListStyle(this);
    });


    //回到顶端
    self.on("click", "#istop", function (evt) {
        gotoTop();
    });

    //打开快捷工具
    self.on("click", "#shortcut menu li", function (evt) {
        runShortcut(this);
    });

    //生成缓存
    self.on("click", "#create_cache", function (evt) {
        createCache(this.attributes["opt"]);
    });

    //重建优化图片缓存
    self.on("click", "#create_optimized li", function (evt) {
        createOptimizedCache(this);
    });


    //清理游戏统计信息
    self.on("click", "#clear_game_stat", function (evt) {
        clearGameStat();
    });

    //清理资料文件
    self.on("click", "#clear_rombase", function (evt) {
        clearRombase();
    });

    //清空缓存
    self.on("click", "#clear_db", function (evt) {
        clearDB(this);
    });

    //弹出平台配置窗口
    self.on("click", "#platform_config", function (evt) {
        openPlatformConfig();
    });

    //弹出设置窗口
    self.on("click", "#config", function (evt) {
        openConfig(this);
    });

    //弹出rom管理窗口
    self.on("click", "#rom_manager", function (evt) {
        openRomManager(this);
    });

    //切换主题
    self.on("click", "#theme menu li", function (evt) {
        changeTheme(this);
    });

    //检测外设（手柄、摇杆）
    self.on("click", "#check_joystick", function (evt) {
        checkJoystick(this);
    });



    //检查更新
    self.on("click", "#upgrade", function (evt) {
        checkUpgrade(this);
    });

    //弹出关于窗口
    self.on("click", "#about", function (evt) {
        openAbout(this);
    });

    //设置菜单 帮助按钮
    self.on("click", "#help", function (evt) {
        runHelp();
    });

    //添加游戏 - 打开rom目录
    self.on("click", "#add_rom_btn .add_rom", function (evt) {
        openRomFolder();
    });

    //添加游戏分身
    self.on("click", "#add_rom_btn .add_slnk", function (evt) {
        addSlnk();
    });

    //添加游戏 - 添加pc游戏
    self.on("click", "#add_rom_btn .add_indie", function (evt) {
        var opt = this.attributes["opt"];
        AddIndieGame(opt);
    });

    //备份rom配置(当前平台)
    self.on("click", "#rom_config_backup_platform", function (evt) {
        romConfigBackup(ACTIVE_PLATFORM);
    });

    //还原rom配置(当前平台)
    self.on("click", "#rom_config_restore_platform", function (evt) {
        romConfigRestore(ACTIVE_PLATFORM);
    });

    //备份rom配置(全部平台)
    self.on("click", "#rom_config_backup_all", function (evt) {
        romConfigBackup(0);
    });

    //还原rom配置(全部平台)
    self.on("click", "#rom_config_restore_all", function (evt) {
        romConfigRestore(0);
    });

    //导出rom配置
    self.on("click", "#rom_config_output", function (evt) {
        openOutput();
    });

    //导入rom
    self.on("click", "#rom_config_input", function (evt) {
        openInput();
    });

    //还原rom配置(全部平台)
    self.on("click", "#merge_db", function (evt) {
        openMergeDb();
    });

    /**
     * 左侧平台边栏
     **/

    self.on("click", "#platform_ul > li", function (evt) {
        changePlatform(this);
    });

    //切换平台标签
    self.on("change", "#platform_tags", function (evt) {
        changePlatformTags(this.value);
    });

    /**
     * 搜索
     **/

    //搜索功能 - 实时改变文本
    self.on("change", "#search_input", function (evt) {
        search();
    });

    //搜索框功能 - 实时改变文本
    self.on("change", "#search_box_input", function (evt) {
        search(2);
    });

    //搜索框失去焦点
    self.on("blur", "#search_box_input", function (evt) {
        $(#search_box).style["display"] = "none";
    });

    //禁用搜索文本框中输入Shift
    self.on("keyup", "#search_box_input", function (evt) {
        if (evt.keyCode == Event.VK_SHIFT) {
            searchBox();
            this.state.focus = false;
        }
    });

    /**
     * 状态栏
     **/

    //基本信息搜索

    self.on("change", "#filter_type", function (evt) {
        search();
    });
    self.on("change", "#filter_producer", function (evt) {
        search();
    });
    self.on("change", "#filter_publisher", function (evt) {
        search();
    });
    self.on("change", "#filter_year", function (evt) {
        search();
    });
    self.on("change", "#filter_country", function (evt) {
        search();
    });
    self.on("change", "#filter_translate", function (evt) {
        search();
    });
    self.on("change", "#filter_version", function (evt) {
        search();
    });
    self.on("change", "#filter_score", function (evt) {
        search();
    });
    self.on("change", "#filter_complete", function (evt) {
        search();
    });
    self.on("click", "#filter_clear", function (evt) {
        filterClear();
    });


    //目录单击
    self.on("click", "#menulist > dd", function (evt) {
        changeMenu(this);
    });

    //添加菜单
    self.on("click", "#add_menu", function (evt) {
        addMenu(this);
    });

    //显示隐藏的游戏
    self.on("click", "#romlist_hide", function (evt) {
        changeHideMenu(this);
    });



    /**
     * 游戏列表
     **/

    //双击游戏模块，启动游戏
    self.on("dblclick", "#romlist li[class=romitem]", function (evt) {
        romListRunGame(this);
    });

    //点击模块，打开侧边栏
    self.on("click", "#romlist li[class=romitem]", function (evt) {
        openSidebar(this);
    });

    //rom分页
    $(#center_content).onScroll = function (evt) {
        scrollLoadRom(evt.scrollPos);
    };

    //加载更多按钮
    self.on("click", "#load_more", function (evt) {
        loadPageRom();
    });

    //按字母搜索rom
    self.on("click", "#num_search li", function (evt) {
        numSearch(this);
    });


    /**
     * 右侧边栏
     **/

    //游戏启动（侧边栏）
    self.on("click", "#buttons li", function (evt) {
        sidebarRunGame(this);
    });

    //游戏启动（侧边栏 - 相关游戏）
    self.on("click", "#related li", function (evt) {
        romListRunGame(this);
    });

    //切换模拟器（侧边栏）
    self.on("click", "#sim_select > li", function (evt) {
        switchRomSim(this);
    });

    //设置rom的cmd（侧边栏）
    self.on("click", "#sim_select > li button", function (evt) {
        SetRomCmd(this);
    });

    //缩略图滑动特效（侧边栏）
    self.on("click", "#rotate", function (evt) {
        thumbSlider(this);
    });

    //缩略图点击显示大图
    self.on("click", "#stack img", function (evt) {
        openBigSlider(this);
    });

    //图集点击显示大图
    self.on("click", "#second_thumbs img", function (evt) {
        openBigSlider(this);
    });

    //缩略图点击隐藏大图
    self.on("click", "#big_thumb_content", function (evt) {
        $(#big_thumb_wrapper).style["display"] = "none";
    });

    //控制视频播放
    self.on("click", "#video_play", function (evt) {
        videoPlay();
    });

    //控制视频音量
    self.on("click", "#right_volume", function (evt) {
        videoVolume(this);
    });

    //打开资料文件
    self.on("click", "#third_files li", function (evt) {
        view.OpenStrategyFiles(this.attributes["path"]);
    });

    //播放音频
    self.on("click", "#audio li", function (evt) {
        playAudio(this.attributes["path"]);
    });

    //播放全部音频
    self.on("click", "#play_audio_all", function (evt) {
        playAudio("");
    });

    //图集展开右键菜单
    self.on("contextmenusetup", ".file-drop-zone", function (evt) {
        secondThumbContext(this);
    });

    //设置通关状态
    self.on("click", "#game_complete", function (evt) {
        setGameComplete(this);
    });





    /**
    * 菜单右键菜单
    **/

    //menu展开右键菜单
    self.on("contextmenusetup", "#menulist dd", function (evt) {
        menuContextMenu(this);
    });


    //menu右键菜单 重命名菜单
    self.on("click", "#menucontext .rename", function (evt) {
        renameMenu(this);
    });

    //menu 右键菜单 删除菜单
    self.on("click", "#menucontext .delete", function (evt) {
        deleteMenu(this);
    });



    /**
    * rom右键菜单
    **/

    //rom展开右键菜单
    self.on("contextmenusetup", "#romlist li", function (evt) {
        //如果是th右键，则忽略
        if (this.attributes.hasClass("romth") == false) {
            openSidebar(this);
            romContextMenu(this);
        }
    });

    //右键菜单启动游戏
    self.on("click", "#romcontext .menu_run_game", function (evt) {
        contextRunGame(this);
    });

    //重命名
    self.on("click", "#romcontext .rename", function (evt) {
        rename(this);
    });

    //编辑资料
    self.on("click", "#romcontext .baseinfo", function (evt) {
        baseinfo(this);
    });


    //编辑子游戏
    self.on("click", "#romcontext .subgame", function (evt) {
        var id = this.attributes["rid"];
        editSubGame(id);
    });

    //编辑游戏攻略
    self.on("click", "#romcontext .strategy", function (evt) {
        editStrategy(this);
    });

    //编辑游戏音频
    self.on("click", "#romcontext .audio", function (evt) {
        editAudio(this);
    });

    //右键菜单设置喜爱
    self.on("click", "#romcontext .fav", function (evt) {
        contextSetFavorite(this);
    });

    //右键菜单设置隐藏
    self.on("click", "#romcontext .hide", function (evt) {
        contextSetHide(this);
    });

    //打开文件夹
    self.on("click", "#romcontext .folder", function (evt) {
        openFolder(this);
    });

    //移动rom
    self.on("click", "#romcontext .move", function (evt) {
        romMove(this);
    });

    //删除rom及相关资源文件
    self.on("click", "#romcontext .delete", function (evt) {

        var id = this.attributes["rid"];
        if (id == undefined) {
            id = ACTIVE_ROM_ID;
        }
        deleteRom(id);
    });

    //侧边栏缩略图编辑功能 - 选择缩略图文件
    self.on("click", "#thumbcontext .openfile", function (evt) {
        openThumbFile(this);
    });

    //侧边栏缩略图编辑功能 - 定位文件目录
    self.on("click", "#thumbcontext .openfolder", function (evt) {
        openFolderBySideSecond(this);
    });

    //侧边栏缩略图编辑功能 - 下载网络图片
    self.on("click", "#thumbcontext .thumb_down", function (evt) {
        var type = this.attributes["value"];
        var sid = this.attributes["sid"];
        thumbDown(type,sid,ACTIVE_ROM_ID);
    });

    //侧边栏缩略图编辑功能 - 导出当前图片
    self.on("click", "#thumbcontext .thumb_output", function (evt) {
        thumbOutput(this);
    });

    //侧边栏缩略图编辑功能 - 将图片设为主图
    self.on("click", "#thumbcontext .thumb_master", function (evt) {
        SetMasterThumb(this);
    });

    //侧边栏缩略图编辑功能 - 删除图片
    self.on("click", "#thumbcontext .thumb_delete", function (evt) {
        DeleteThumb(this);
    });


    /**
     * 界面设置
     **/

    //更改rom列表展示图显示方向
    self.on("click", "#config_romlist_direction> menu > li menu > li", function (evt) {
        setThumbDirection(this);
    });

    //更改rom列表展示图类型
    self.on("click", "#config_thumb > menu > li menu > li", function (evt) {
        setThumbType(this);
    });

    //是否显示rom的标题背景
    self.on("click", "#config_font_background menu li", function (evt) {
        setFontBackgrond(this);
    });

    //更改字体大小
    self.on("click", "#config_title_fontsize menu li", function (evt) {
        setFontsize(this);
    });

    //更改rom模块大小
    self.on("click", "#config_romlist_size menu li", function (evt) {
        setRomlistSize(this);
    });

    //更改rom列表展示模块间距
    self.on("click", "#config_romlist_margin menu li", function (evt) {
        setRomMargin(this);
    });

    //列表列显示
    self.on("click", ".romlist_column", function (evt) {
        setRomlistColumn(this);
    });

    //列表名称显示类型
    self.on("click", "#config_show_name_type menu li", function (evt) {
        setShowNameType(this);
    });

    //列表排序方式
    self.on("click", "#config_orders menu li", function (evt) {
        setListSort(this);
    });


    /**
    * 星级控制
    **/

    //鼠标移入
    self.on("mousemove", "#score span", function (e) {
        mousemoveScore(this, e);
    });

    //鼠标点击设置星级
    self.on("click", "#score span", function (e) {
        setScore();
    });

    //鼠标离开
    $(#score).on("mouseleave", function (evt) {
        mouseleaveScore();
    });




    /**
    * 键盘按下
    **/
    self.on("keyup", function (evt) {

        //解决在输入法选字的情况下，移动rom模块
        var ipt = $$(input);
        for (var i in ipt) {
            if (i.state.focus == true) {
                return true;
            }
        }

        var romlistType = $(#switch_romlist).attributes["value"]; //rom列表样式(模块或列表)
        var romContextVisible = $(#romcontext).style["visibility"] == "visible" ? true : false; //rom右键菜单是否显示中
        var currentRom = $(#romlist).select("li:current");

        switch (evt.keyCode) {
            case Event.VK_UP: //上键移动rom
                if (romlistType == "1") { //模块
                    keyRomUp();
                } else if (romlistType == "2") { //列表
                    keyRomLeft();
                }
                scrollUp(romlistType); //div滚动
                break;
            case Event.VK_DOWN: //下键移动rom
                if (romlistType == "1") { //模块
                    keyRomDown();
                } else if (romlistType == "2") { //列表
                    keyRomRight();
                }
                scrollDown(romlistType); //div滚动
                break;
            case Event.VK_LEFT: //左键移动rom
                if (romlistType == "1") { //模块
                    scrollUp(romlistType); //div滚动
                    keyRomLeft();
                }
                break;
            case Event.VK_RIGHT: //右键移动rom
                if (romlistType == "1") { //模块
                    scrollDown(romlistType); //div滚动
                    keyRomRight();
                }
                break;
            case Event.VK_RETURN: //回车键启动游戏
                romListRunGame(currentRom);
                break;
            case Event.VK_F5: //F5刷新缓存
                createCache();
                break;
            case Event.VK_F1: //F1启动帮助
                runHelp();
                break;
            case Event.VK_SHIFT: //shift启动搜索框
                searchBox();
                break;
            case Event.VK_R: //重命名
                if (!romContextVisible) return true;
                rename(currentRom);
                break;
            case Event.VK_D: //删除
                if (!romContextVisible) return true;
                deleteRom(currentRom.attributes["rid"]);
                break;
            case Event.VK_M: //移动
                if (!romContextVisible) return true;
                romMove(currentRom);
                break;
            case Event.VK_F: //设为喜爱
                if (!romContextVisible) return true;
                keycodeSetFavorite();
                break;
            case Event.VK_H: //设为隐藏
                if (!romContextVisible) return true;
                keycodeSetHide();
                break;
            case Event.VK_E: //编辑资料
                if (!romContextVisible) return true;
                baseinfo(currentRom);
                break;
            case Event.VK_G: //编辑简介和攻略
                if (!romContextVisible) return true;
                editStrategy(currentRom);
                break;
            case Event.VK_S: //编辑子游戏
                if (!romContextVisible) return true;
                editSubGame(currentRom.attributes["rid"]);
                break;
            case Event.VK_A: //编辑音频
                if (!romContextVisible) return true;
                editAudio(currentRom);
                break;
            case Event.VK_HOME: //回到顶部
                gotoTop();
                break;
            case Event.VK_PRIOR: //上翻页
                gotoPageUp();
                break;
            case Event.VK_NEXT: //下翻页
                gotoPageDown();
                break;
            case Event.VK_1:
                if (currentRom == undefined) return;
                numRunGame(1);
                break;
            case Event.VK_2:
                if (currentRom == undefined) return;
                numRunGame(2);
                break;
            case Event.VK_3:
                if (currentRom == undefined) return;
                numRunGame(3);
                break;
            case Event.VK_4:
                if (currentRom == undefined) return;
                numRunGame(4);
                break;
            case Event.VK_5:
                if (currentRom == undefined) return;
                numRunGame(5);
                break;
            case Event.VK_6:
                if (currentRom == undefined) return;
                numRunGame(6);
                break;
            case Event.VK_7:
                if (currentRom == undefined) return;
                numRunGame(7);
                break;
            case Event.VK_8:
                if (currentRom == undefined) return;
                numRunGame(8);
                break;
            case Event.VK_9:
                if (currentRom == undefined) return;
                numRunGame(9);
                break;
        }

    });

});



