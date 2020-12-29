package initializer

import (
	"github.com/TeamZenithy/Araha/lang"
	"github.com/TeamZenithy/Araha/utils"
)

func InitLang() {
	tr := lang.NewTr()
	tr.AddLang(lang.NewTrLocale("lang/ko"))
	tr.AddLang(lang.NewTrLocale("lang/en"))
	utils.TR = tr
}
