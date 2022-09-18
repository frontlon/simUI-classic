package compoments

import (
	"simUI/code/db"
	"simUI/code/utils"
	"strings"
)

//开始写入rom_subgame数据
func FilterFactory(volist []*db.Rom, dbFilter []*db.Filter) (map[string][]string, map[string][]string) {

	baseType := map[string]bool{}
	baseYear := map[string]bool{}
	baseProducer := map[string]bool{}
	basePublisher := map[string]bool{}
	baseCountry := map[string]bool{}
	baseTranslate := map[string]bool{}
	baseVersion := map[string]bool{}
	score := map[string]bool{}
	complete := map[string]bool{}

	for _, v := range volist {
		if v.BaseYear != "" {
			year := v.BaseYear
			if len(v.BaseYear) > 4 {
				year = v.BaseYear[0:4]
			}
			baseYear[year] = true
		}
		if v.BaseType != "" {
			items := strings.Split(v.BaseType, ";")
			for _, pr := range items {
				baseType[strings.Trim(pr," ")] = true
			}
		}
		if v.BaseProducer != "" {
			items := strings.Split(v.BaseProducer, ";")
			for _, pr := range items {
				baseProducer[strings.Trim(pr," ")] = true
			}
		}
		if v.BasePublisher != "" {
			items := strings.Split(v.BasePublisher, ";")
			for _, pr := range items {
				basePublisher[strings.Trim(pr," ")] = true
			}
		}
		if v.BaseCountry != "" {
			baseCountry[v.BaseCountry] = true
		}
		if v.BaseTranslate != "" {
			baseTranslate[v.BaseTranslate] = true
		}
		if v.Score != "" {
			score[v.Score] = true
		}

		complete[utils.ToString(v.Complete)] = true
	}

	baseTypeList := []string{}
	baseYearList := []string{}
	baseProducerList := []string{}
	basePublisherList := []string{}
	baseCountryList := []string{}
	baseTranslateList := []string{}
	baseVersionList := []string{}
	scoreList := []string{}
	completeList := []string{}

	for k, _ := range baseType {
		baseTypeList = append(baseTypeList, k)
	}
	for k, _ := range baseYear {
		baseYearList = append(baseYearList, k)
	}
	for k, _ := range baseProducer {
		baseProducerList = append(baseProducerList, k)
	}
	for k, _ := range basePublisher {
		basePublisherList = append(basePublisherList, k)
	}
	for k, _ := range baseCountry {
		baseCountryList = append(baseCountryList, k)
	}
	for k, _ := range baseTranslate {
		baseTranslateList = append(baseTranslateList, k)
	}
	for k, _ := range baseVersion {
		baseVersionList = append(baseVersionList, k)
	}
	for k, _ := range score {
		scoreList = append(scoreList, k)
	}
	for k, _ := range complete {
		completeList = append(completeList, k)
	}

	createVo := map[string][]string{}
	createVo["base_type"] = append(createVo["base_type"], baseTypeList...)
	createVo["base_year"] = append(createVo["base_year"], baseYearList...)
	createVo["base_producer"] = append(createVo["base_producer"], baseProducerList...)
	createVo["base_publisher"] = append(createVo["base_publisher"], basePublisherList...)
	createVo["base_country"] = append(createVo["base_country"], baseCountryList...)
	createVo["base_translate"] = append(createVo["base_translate"], baseTranslateList...)
	createVo["base_version"] = append(createVo["base_version"], baseVersionList...)
	createVo["score"] = append(createVo["score"], scoreList...)
	createVo["complete"] = append(createVo["complete"], completeList...)

	//整理过滤器数据
	createDb := map[string][]string{}
	for _, v := range dbFilter {
		if _, ok := createDb[v.Type]; ok {
			createDb[v.Type] = append(createDb[v.Type], v.Name)
		} else {
			createDb[v.Type] = []string{v.Name}
		}
	}

	return createVo, createDb
}
