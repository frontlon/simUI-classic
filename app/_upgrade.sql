DROP TABLE filter;
CREATE TABLE "filter" ("id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,"platform"  INTEGER NOT NULL DEFAULT 0,"type"  TEXT NOT NULL,"name"  TEXT NOT NULL DEFAULT '');
ALTER TABLE menu ADD COLUMN "sort" INTEGER NOT NULL DEFAULT 0
ALTER TABLE platform ADD COLUMN "title_path" TEXT NOT NULL DEFAULT '',
ALTER TABLE platform ADD COLUMN "background_path" TEXT NOT NULL DEFAULT '',
ALTER TABLE platform ADD COLUMN "doc_path" TEXT NOT NULL DEFAULT '',
ALTER TABLE platform ADD COLUMN "rombase" TEXT NOT NULL DEFAULT '',
ALTER TABLE platform ADD COLUMN "cassette_path"  TEXT NOT NULL,
ALTER TABLE platform ADD COLUMN "icon_path"  TEXT NOT NULL,
ALTER TABLE platform ADD COLUMN "gif_path"  TEXT NOT NULL,
ALTER TABLE platform DROP COLUMN romlist;
ALTER TABLE simulator ADD COLUMN "sort" INTEGER NOT NULL DEFAULT 0
ALTER TABLE simulator ADD COLUMN "lua" TEXT NOT NULL DEFAULT ''
DROP TABLE rom;
CREATE TABLE "rom" ("id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,"platform"  TEXT NOT NULL,"menu"  TEXT NOT NULL,"name"  TEXT NOT NULL,"pname"  TEXT NOT NULL,"rom_path"  TEXT NOT NULL,"star"  INTEGER NOT NULL DEFAULT 0,"sim_id"  INTEGER NOT NULL DEFAULT 0,"sim_conf"  TEXT NOT NULL DEFAULT "{}","hide"  INTEGER NOT NULL DEFAULT 0,"run_num"  INTEGER NOT NULL DEFAULT 0,"run_time"  INTEGER NOT NULL DEFAULT 0,"base_type"  TEXT NOT NULL,"base_year"  INTEGER NOT NULL,"base_publisher"  TEXT NOT NULL,"base_country"  TEXT NOT NULL,"base_translate"  TEXT NOT NULL,"pinyin"  TEXT NOT NULL,"file_md5"  TEXT NOT NULL DEFAULT '',"info_md5"  TEXT NOT NULL DEFAULT 0);
CREATE INDEX "idx_info_md5" ON "rom" ("info_md5" ASC);
CREATE INDEX "idx_pf_menu" ON "rom" ("platform" ASC, "menu" ASC);
CREATE INDEX "idx_pf_pname" ON "rom" ("platform" ASC, "pname" ASC);
CREATE INDEX "idx_pinyin" ON "rom" ("pinyin" ASC);
ALTER TABLE rom ADD COLUMN "base_version" TEXT NOT NULL,
DROP TABLE config;
CREATE TABLE "config" ("id" INTEGER NOT NULL, "lang" TEXT NOT NULL, "theme" TEXT NOT NULL, "platform" INTEGER NOT NULL, "menu" TEXT NOT NULL, "thumb" TEXT NOT NULL, "search_engines" TEXT NOT NULL, "root_path" TEXT NOT NULL, "window_width" INTEGER NOT NULL DEFAULT 0, "window_height" INTEGER NOT NULL DEFAULT 0, "window_state" INTEGER NOT NULL, "upgrade_id" INTEGER NOT NULL DEFAULT 1, "soft_name" TEXT NOT NULL DEFAULT 0, "enable_upgrade" INTEGER NOT NULL DEFAULT 1, "panel_platform" INTEGER NOT NULL DEFAULT 1, "panel_menu" INTEGER NOT NULL DEFAULT 1, "panel_sidebar" INTEGER NOT NULL DEFAULT 1, "panel_platform_width" TEXT NOT NULL, "panel_menu_width" TEXT NOT NULL, "panel_sidebar_width" TEXT NOT NULL, "romlist_size" INTEGER NOT NULL DEFAULT 2, "romlist_margin" INTEGER NOT NULL, "romlist_style" INTEGER NOT NULL DEFAULT 1, "romlist_direction" INTEGER NOT NULL, "romlist_font_background" INTEGER NOT NULL, "romlist_font_size" INTEGER NOT NULL DEFAULT 1, "romlist_column" TEXT NOT NULL, "background_image" TEXT NOT NULL, "background_repeat" TEXT NOT NULL, "background_opacity" TEXT NOT NULL, "cursor" TEXT NOT NULL DEFAULT '', "video_volume" INTEGER NOT NULL, PRIMARY KEY ("id") );
INSERT INTO "config" VALUES (1, '简体中文', 'dark', 1, '', 'thumb', 'https://image.baidu.com/search/acjson?tn=resultjson_com&logid=8695915517170565368&ipn=rj&ct=201326592&is=&fp=result&queryWord=fc&cl=2&lm=-1&ie=utf-8&oe=utf-8&adpicid=&st=&z=&ic=&hd=&latest=&copyright=&word={$keyword}&s=&se=&tab=&width=&height=&face=&istype=&qc=&nc=1&fr=&expermode=&force=&pn={$NumIndex}&rn={$pageNum}&gsm=78&1622211780093=', '/Users/frontlon/go/src/simUI/app/', 1388, 873, 1, 2, 'simUI', 1, 1, 1, 1, 88, 156, 268, 2, 2, 1, 0, 2, 2, '1,0,1,1,1,1', '/Users/frontlon/Desktop/20210610895319.jpg', 'no-repeat', 100, '', 1);