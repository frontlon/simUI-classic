package main

var constPlatformList = []*Platform{}

type Platform struct{
	Name string
	Value string
}

/**
 * 读取所有可用的平台菜单
 **/
func GetPlatformData()  {
	plat := &Platform{}
	if Config.Fc.Enable == "1"{
		plat = &Platform{
			Name: Config.Fc.Title,
			Value:"fc",
		}
		constPlatformList = append(constPlatformList,plat)
	}

	if Config.Sfc.Enable == "1"{
		plat = &Platform{
			Name: Config.Sfc.Title,
			Value:"sfc",
		}
		constPlatformList = append(constPlatformList,plat)
	}

	if Config.Md.Enable == "1"{
		plat = &Platform{
			Name: Config.Md.Title,
			Value:"md",
		}
		constPlatformList = append(constPlatformList,plat)
	}

	if Config.Pce.Enable == "1"{
		plat := &Platform{
			Name: Config.Pce.Title,
			Value:"pce",
		}
		constPlatformList = append(constPlatformList,plat)
	}

	if Config.Gb.Enable == "1"{
		plat := &Platform{
			Name: Config.Gb.Title,
			Value:"gb",
		}
		constPlatformList = append(constPlatformList,plat)
	}

	if Config.Arcade.Enable == "1"{
		plat := &Platform{
			Name: Config.Arcade.Title,
			Value:"arcade",
		}
		constPlatformList = append(constPlatformList,plat)
	}

}
