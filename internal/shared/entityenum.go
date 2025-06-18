package shared

type Language string

const (
	LanguageEnglish  Language = "ENGLISH"
	LanguageGeorgian Language = "GEORGIAN"
)

func (Language) Values() (kinds []string) {
	for _, s := range []Language{LanguageEnglish, LanguageGeorgian} {
		kinds = append(kinds, string(s))
	}
	return
}
