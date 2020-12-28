package lang

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"
)

type TrText struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type TrLocale struct {
	Locale      string
	Transitions map[string]*[]TrText
}

type Tr struct {
	TrTexts map[string]*TrLocale
}

func NewTrLocale(folder string) *TrLocale {
	var trLocale TrLocale

	trLocale.Locale = filepath.Base(folder)
	trLocale.Transitions = map[string]*[]TrText{}

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		data, _ := ioutil.ReadFile(folder + "/" + file.Name())
		var trTexts []TrText
		err := json.Unmarshal(data, &trTexts)
		if err != nil {
			panic(err)
		}
		fileName := strings.TrimSuffix(file.Name(), path.Ext(file.Name()))
		trLocale.Transitions[fileName] = &trTexts
	}

	return &trLocale
}

func NewTr() *Tr {
	return &Tr{TrTexts: map[string]*TrLocale{}}
}

func (t *Tr) AddLang(locale *TrLocale) {
	t.TrTexts[locale.Locale] = locale
}

func (t *Tr) GetHandlerFunc(lang, fallback string) func(string, ...string) string {

	find := func(lang, file, key string) (string, bool) {
		for _, d := range *t.TrTexts[lang].Transitions[file] {
			if d.Name == key {
				return d.Text, true
			}
		}
		return key, false
	}

	return func(key string, text ...string) string {
		spilited := strings.Split(key, ":")
		if len(spilited) != 2 {
			return key
		}
		trText, exits := find(lang, spilited[0], spilited[1])
		if !exits {
			trText, exits = find(fallback, spilited[0], spilited[1])
			if !exits {
				return key
			}
		}
		for _, d := range text {
			trText = strings.Replace(trText, "{}", d, 1)
		}
		return trText
	}
}
