package types

import "github.com/beego/beego/v2/core/logs"

type LanguageTypeEnum int

const (
	LanguageTypeNone    LanguageTypeEnum = 0
	LanguageTypeChinese LanguageTypeEnum = 1
	LanguageTypeEnglish LanguageTypeEnum = 2
)

const (
	LangEnUS string = "en"
	LangZhCN string = "zh"

	LanguageTypeEnglishDisplay = `English`
	LanguageTypeChineseDisplay = `简体中文`
)

func LanguageTypeEnumConf() map[LanguageTypeEnum]string {
	return map[LanguageTypeEnum]string{
		LanguageTypeChinese: LanguageTypeChineseDisplay,
		LanguageTypeEnglish: LanguageTypeEnglishDisplay,
	}
}

func LanguageTypeConf() map[LanguageTypeEnum]string {
	return map[LanguageTypeEnum]string{
		LanguageTypeChinese: LangZhCN,
		LanguageTypeEnglish: LangEnUS,
	}
}

var LangShort2TypeConf = map[string]LanguageTypeEnum{
	LangZhCN: LanguageTypeChinese,
	LangEnUS: LanguageTypeEnglish,
}

func langType2ShortConf() map[LanguageTypeEnum]string {
	var conf = make(map[LanguageTypeEnum]string)
	for s, l := range LangShort2TypeConf {
		conf[l] = s
	}

	return conf
}

func Short2LangTypeConf() map[string]LanguageTypeEnum {
	conf := map[string]LanguageTypeEnum{}
	for t, s := range langType2ShortConf() {
		conf[s] = t
	}

	return conf
}

func Short2LangType(lang string) LanguageTypeEnum {
	conf := Short2LangTypeConf()
	if t, ok := conf[lang]; ok {
		return t
	} else {
		logs.Info("lang is undefined, lang: %s", lang)
		return LanguageTypeNone
	}
}
