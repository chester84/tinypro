package types

type AstroEnum int

// 十二星座英文名：
// aries 白羊座, taurus 金牛座,
// gemini 双子座, cancer 钜蟹座,
// leo 狮子座, virgo 处女座,
// libra 天平座, scorpio 天蝎座,
// sagittarius 射手座, capricorn 摩羯座,
// aquarius 水瓶座, pisces 双鱼座。

const (
	AstroAries       AstroEnum = 1
	AstroTaurus      AstroEnum = 2
	AstroGemini      AstroEnum = 3
	AstroCancer      AstroEnum = 4
	AstroLeo         AstroEnum = 5
	AstroVirgo       AstroEnum = 6
	AstroLibra       AstroEnum = 7
	AstroScorpio     AstroEnum = 8
	AstroSagittarius AstroEnum = 9
	AstroCapricorn   AstroEnum = 10
	AstroAquarius    AstroEnum = 11
	AstroPisces      AstroEnum = 12
)

type AstroItem struct {
	AstroSN AstroEnum `json:"astro_sn"`

	En string `json:"en"`
	Zh string `json:"zh"`
}

func AstroConfig() map[AstroEnum]AstroItem {
	return map[AstroEnum]AstroItem{
		AstroAries: {
			AstroAries,
			`aries`,
			`白羊座`,
		},
		AstroTaurus: {
			AstroTaurus,
			`taurus`,
			`金牛座`,
		},
		AstroGemini: {
			AstroGemini,
			`gemini`,
			`双子座`,
		},
		AstroCancer: {
			AstroCancer,
			`cancer`,
			`钜蟹座`,
		},
		AstroLeo: {
			AstroLeo,
			`leo`,
			`狮子座`,
		},
		AstroVirgo: {
			AstroVirgo,
			`virgo`,
			`处女座`,
		},
		AstroLibra: {
			AstroLibra,
			`libra`,
			`天平座`,
		},
		AstroScorpio: {
			AstroScorpio,
			`scorpio`,
			`天蝎座`,
		},
		AstroSagittarius: {
			AstroSagittarius,
			`sagittarius`,
			`射手座`,
		},
		AstroCapricorn: {
			AstroCapricorn,
			`capricorn`,
			`摩羯座`,
		},
		AstroAquarius: {
			AstroAquarius,
			`aquarius`,
			`水瓶座`,
		},
		AstroPisces: {
			AstroPisces,
			`pisces`,
			`双鱼座`,
		},
	}
}
