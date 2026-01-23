package i18n

func Translate(code, lang string) string {
	if msg, ok := Messages[code]; ok {
		if text, ok := msg[lang]; ok {
			return text
		}
		return msg["en"]
	}
	return "Unknown error"
}
