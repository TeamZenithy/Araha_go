package initializer

import (
	"github.com/TeamZenithy/Araha/lang"
	"github.com/TeamZenithy/Araha/utils"
)

func InitLang() {
	tr := lang.NewTr()
	tr.AddLang(lang.NewTrLocale("static/lang/ko"))
	tr.AddLang(lang.NewTrLocale("static/lang/en"))
	utils.TR = tr
}
