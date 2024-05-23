package utils

import (
	"regexp"
	"strings"

	"github.com/gofiber/contrib/fiberi18n"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// Helper for language/localizer
func L(key string, c *fiber.Ctx) string {
	lang := fiberi18n.MustGetMessage(key)
	if lang == "" {
		return CamelCaseToOrdinaryString(key)
	}
	return lang
}

// Change string from <camelCaseString> to <Camel case string> (Sentence)
// Unicode compatible version
func CamelCaseToOrdinaryString(s string) string {
	pattern := regexp.MustCompile("(\\p{Lu}+\\P{Lu}*)")
	s = pattern.ReplaceAllString(s, "${1} ")
	s, _ = strings.CutSuffix(strings.ToLower(s), " ")
	return FirstToUpper(s)
}

// Change string from <string> to <String> (Upper first letter)
// Unicode compatible version
func FirstToUpper(s string) string {
	byteSlice := []byte(s)
	if len(byteSlice) == 0 {
		return ""
	}
	byteSlice[0] = byte(byteSlice[0] - 32)
	return string(byteSlice)
}

// var Lang string

type Language struct {
	Lang string `json:"lang" xml:"lang" form:"lang"`
}

func SetLocale(c *fiber.Ctx, defaultLang string) string {
	p := new(Language)
	if err := c.BodyParser(p); err != nil {
		c.Locals("lang", defaultLang)
		return defaultLang
	}
	var lang string
	lang = utils.CopyString(p.Lang) // c.Query("lang")
	if lang != "" {
		c.Locals("lang", lang)
		return lang
	}
	// lang = utils.CopyString(c.Get("Accept-Language"))
	// if lang != "" {
	// 	c.Locals("lang", lang)
	// 	return lang
	// }
	c.Locals("Lang", defaultLang)
	return defaultLang
}
