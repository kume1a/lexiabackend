package schema

type Language string

const (
	LanguageEnglish  Language = "ENGLISH"
	LanguageGeorgian Language = "GEORGIAN"
	LanguageSpanish  Language = "SPANISH"
	LanguageFrench   Language = "FRENCH"
	LanguageGerman   Language = "GERMAN"
	LanguageRussian  Language = "RUSSIAN"
	LanguageJapanese Language = "JAPANESE"
	LanguageChinese  Language = "CHINESE"
)

func (Language) Values() (kinds []string) {
	for _, s := range []Language{
		LanguageEnglish,
		LanguageGeorgian,
		LanguageSpanish,
		LanguageFrench,
		LanguageGerman,
		LanguageRussian,
		LanguageJapanese,
		LanguageChinese,
	} {
		kinds = append(kinds, string(s))
	}
	return
}

type FolderType string

const (
	FolderTypeFolderCollection FolderType = "FOLDER_COLLECTION"
	FolderTypeWordCollection   FolderType = "WORD_COLLECTION"
)

func (FolderType) Values() (kinds []string) {
	for _, s := range []FolderType{FolderTypeFolderCollection, FolderTypeWordCollection} {
		kinds = append(kinds, string(s))
	}
	return
}
