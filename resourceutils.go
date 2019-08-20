package game2e

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/tryor/game2e/log"
)

type ResourceInfo struct {
	//respath+"material.images.json", respath+"background.images.json", respath+"hero.images.json", respath+"normal.images.json", respath+"lighteffect.images.json"
	//respath + "animations.json"
	//spiritspath + "spirits.json"
	//elementspath + "normal.elements.json"
	//mappath + "normal.map.json"
	ImageResRoot      string   `json:"imgresroot" xml:"imgresroot"`
	ImageResFiles     []string `json:"imagefiles" xml:"imagefiles"`
	AnimationResFiles []string `json:"animationfiles" xml:"animationfiles"`
	SpiritResFiles    []string `json:"spiritfiles" xml:"spiritfiles"`
	ElementResFiles   []string `json:"elementfiles" xml:"elementfiles"`
	MapResFiles       []string `json:"mapfiles" xml:"mapfiles"`
}

var Resource ResourceInfo

func LoadResourceInfo(resconf string) error {
	err := loadResourceInfo(resconf)
	if err != nil {
		return err
	}

	err = LoadImageResDicts(Resource.ImageResFiles...)
	if err != nil {
		log.Error("load *.images.json error!", err)
		return err
	}

	err = LoadAnimationInfos(Resource.AnimationResFiles...)
	if err != nil {
		log.Error("load animations.json error!", err)
		return err
	}

	err = LoadSpiritInfos(Resource.SpiritResFiles...)
	if err != nil {
		log.Error("load spirits.json error!", err)
		return err
	}

	err = LoadResElementInfos(Resource.ElementResFiles...)
	if err != nil {
		log.Error("load elements.json error!", err)
		return err
	}

	err = LoadMapInfos(Resource.MapResFiles...)
	if err != nil {
		log.Error("load map.json error!", err)
		return err
	}

	return nil
}

func loadResourceInfo(resconf string) error {
	file, err := os.Open(resconf)
	if err != nil {
		return err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, &Resource)
	if err != nil {
		return err
	}

	//log.Info(Resource)

	return nil
}
