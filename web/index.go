package web

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/TeamZenithy/Araha/model"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/gin-gonic/gin"
)

func templateDict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func templateAdd(a, b int) (int, error) {
	return a + b, nil
}

func InitWeb() {
	r := gin.Default()
	funcMap := template.FuncMap{
		"dict": templateDict,
		"add":  templateAdd,
	}
	r.SetFuncMap(funcMap)
	r.LoadHTMLGlob("static/web/public/*")
	r.GET("/:lang/queue/:guildID", Lang(), handleQueue)
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "notFound.html", gin.H{})
	})
	r.Run(":8096")
}

func Lang() gin.HandlerFunc {
	return func(c *gin.Context) {
		langParm := c.Param("lang")
		c.Set("T", utils.TR.GetHandlerFunc(langParm, "en"))
	}
}

func handleQueue(c *gin.Context) {
	// T := c.MustGet("T").(lang.HFType)
	guildID := c.Param("guildID")

	ms, ok := model.Music[guildID]
	c.HTML(http.StatusOK, "queue.html", gin.H{
		"hasData": ok,
		"texts":   gin.H{},
		"queue":   ms,
	})
}
