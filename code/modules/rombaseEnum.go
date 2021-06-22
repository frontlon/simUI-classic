package modules

import (
	"simUI/code/db"
)

//更新rombase枚举
func UpdateRomBaseEnum(t string, data []string) error {

	//先删除记录
	_ = (&db.RombaseEnum{Type: t}).DeleteByType()

	if len(data) == 0 {
		return nil
	}

	create := []*db.RombaseEnum{}
	for _, v := range data {
		if v == "" {
			continue
		}
		c := &db.RombaseEnum{}
		c.Type = t
		c.Name = v
		create = append(create, c)
	}
	if err := (&db.RombaseEnum{}).BatchAdd(create); err != nil {
		return err
	}
	return nil
}
