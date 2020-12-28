package initializer

import (
	"github.com/TeamZenithy/Araha/lang"
	"github.com/TeamZenithy/Araha/utils"
)

func InitLang() {
	tr := lang.NewTr()
	tr.AddLang(lang.NewTrLocale("lang/ko"))
	utils.TR = &tr
}
