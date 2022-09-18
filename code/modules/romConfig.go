package modules

import (
	"simUI/code/db"
)

func SetFavorite(id uint64, star uint8) error {

	//数据库中读取rom详情
	rom, err := (&db.Rom{}).GetById(id)
	if err != nil {
		return err
	}

	//更新数据
	err = (&db.RomSetting{
		FileMd5:  rom.FileMd5,
		Platform: rom.Platform,
		Star:     star,
	}).UpdateStar()
	if err != nil {
		return err
	}

	err = (&db.Rom{
		Id:   id,
		Star: star,
	}).UpdateStar()

	if err != nil {
		return err
	}

	return nil
}

//设为隐藏
func SetHide(id uint64, hide uint8) error {

	//数据库中读取rom详情
	rom, err := (&db.Rom{}).GetById(id)
	if err != nil {
		return err
	}

	//更新数据
	err = (&db.RomSetting{
		FileMd5:  rom.FileMd5,
		Platform: rom.Platform,
		Hide:     hide,
	}).UpdateHide()

	err = (&db.Rom{
		Id:   id,
		Hide: hide,
	}).UpdateHide()

	if err != nil {
		return err
	}

	return nil
}

func SetHideBatch(ids []uint64, hide uint8) error {

	//数据库中读取rom详情
	roms, err := (&db.Rom{}).GetByIds(ids)
	if err != nil {
		return err
	}
	if len(roms) == 0 {
		return nil
	}

	md5s := []string{}
	for _, v := range roms {
		md5s = append(md5s, v.FileMd5)
	}

	//更新数据
	if err := (&db.RomSetting{}).UpdateHideByFileMd5(roms[0].Platform, md5s, hide); err != nil {
		return err
	}

	if err := (&db.Rom{}).UpdateHideByIds(ids, hide); err != nil {
		return err
	}

	return nil
}
