/**
 * 全局变量列表
 */

var CONF; //读取的config配置
var ROMJSON;//rom列表
var MAXZOOM = 6; //模块最大缩放等级

var ACTIVE_PLATFORM = 0; //当前选定的平台
var ACTIVE_MENU = ""; //当前选定的菜单
var ACTIVE_ROM_ID = 0; //当前rom id

var MENU_SCROLL_POS = 0; //默认滚动条位置
var MENU_SCROLL_LOCK = false; //菜单滚动条锁
var MENU_SCROLL_PAGE = 1; //当前滚动条翻页页数

var SCROLL_POS = 0; //默认滚动条位置
var SCROLL_LOCK = false; //rom列表滚动条锁
var SCROLL_PAGE = 0; //当前滚动条翻页页数

var VIDEO_PLAY_STATE = 0; //记录当前侧边栏视频的播放状态，1正在播放，0未播放
var VIDEO_ALLOW_PLAY = 0; //是否允许播放视频，用于解决视频异步播放时的暂停问题

var VIDEO_FILTER = "*.wmv;*.mp4;*.avi;*.flv;*.webm";
var PIC_FILTER = "*.gif;*.png;*.jpg;*.jpeg;*.ico;*.bmp;*.webp";
var THUMB_EMPTY_SID = "_j5D"; //侧边栏图集空模块sid

var ROOMLIST_SIZE_NUM = 19; //模块字体大小的数量
var ROOMLIST_MARGIN_NUM = 15; //模块间距的数量
var ROOMLIST_ROOM_SIZE_NUM = 16; //模块大小的数量