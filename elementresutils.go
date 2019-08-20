package game2e

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/tryor/game2e/log"
)

func init() {
	ResElementInfos = make(ResElementInfoMap)
}

type ResElementInfoMap map[string]*ResElementInfo

var ResElementInfos ResElementInfoMap

func LoadResElementInfos(filenames ...string) error {
	for _, filename := range filenames {
		reselinfos, err := loadResElementInfos(filename)
		if err != nil {
			return err
		}

		for key, val := range reselinfos {
			if _, ok := ResElementInfos[key]; ok {
				log.Error("resource element is exist! key is ", key)
			}
			ResElementInfos[key] = val
		}

	}

	//	log.Info("ResElementInfos:", ResElementInfos)
	//	for mapkey, mapInfo := range ResElementInfos {
	//		log.Info("LoadResElementInfos:", mapkey, mapInfo)
	//	}

	return nil
}

func loadResElementInfos(filename string) (ResElementInfoMap, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	resElementInfos := make(ResElementInfoMap)
	err = json.Unmarshal(buf, &resElementInfos)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return resElementInfos, nil
}
