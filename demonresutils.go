package game2e

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func init() {
	SpiritInfos = make(map[string]*SpiritInfo)
}

/*
{
"spirit_$legendaryswordsman":{
    "width":85,
    "height":87,
    "features":{
        "stands":[{"angles":["0-360"], "animation":"$legendaryswordsman_1"}],
        "walks":[
	        {"angles":["0-90", "270-360"], "animation":"$legendaryswordsman_forward"},
	        {"angles":["90-270"], "animation":"$legendaryswordsman_rollback"}
	]
    }
}
}
*/

type SpiritInfo struct {
	Width    int                   `json:"width" xml:"width"`
	Height   int                   `json:"height" xml:"height"`
	Features map[string][]*Feature `json:"features" xml:"features"`
	Name     ResFontInfo           `json:"name" xml:"name"`
	Default  string                `json:"default" xml:"default"` //默认动画, 一般为站着不动时
	Moved    string                `json:"moved" xml:"moved"`     //被移动时动画

}

type Feature struct {
	Angles      [][2]float32
	Angles_     []string       `json:"angles" xml:"angles"`
	Animation   string         `json:"animation" xml:"animation"`
	Frameindexs []int          `json:"frameindexs" xml:"frameindexs"`
	Feature     ResFeatureInfo `json:"feature" xml:"feature"`
	//	Alignment   Alignment `json:"alignment" xml:"alignment"`
	//	Anchorpoint PointF    `json:"anchorpoint" xml:"anchorpoint"`
}

var SpiritInfos map[string]*SpiritInfo

func LoadSpiritInfos(filenames ...string) error {
	for _, filename := range filenames {
		err := loadSpiritInfos(filename)
		if err != nil {
			return err
		}
	}

	//	log.Println("LoadSpiritInfos:", SpiritInfos)
	for _, spiritInfo := range SpiritInfos {
		//		log.Println("spiritInfo:", spiritInfo)
		for _, features := range spiritInfo.Features {
			//			log.Println("walk.angles_:", key, features)
			for _, feature := range features {
				feature.Angles = make([][2]float32, len(feature.Angles_))
				for i, angles_ := range feature.Angles_ {
					angles := strings.Split(angles_, "-")
					//var err error
					angle, err := strconv.ParseFloat(angles[0], 32)
					if err != nil {
						return err
					}
					feature.Angles[i][0] = float32(angle)
					angle, err = strconv.ParseFloat(angles[1], 32)
					if err != nil {
						return err
					}
					feature.Angles[i][1] = float32(angle)
				}
			}

		}
	}

	return nil
}

//func LoadSpiritInfos(filenames ...string) error {
//	for _, filename := range filenames {
//		err := loadSpiritInfos(filename)
//		if err != nil {
//			return err
//		}
//	}

//	//	log.Println("LoadSpiritInfos:", SpiritInfos)
//	for _, spiritInfo := range SpiritInfos {
//		//		log.Println("spiritInfo:", spiritInfo)
//		for _, walk := range spiritInfo.Feature.Walks {
//			//			log.Println("walk.angles_:", walk.Angles_)
//			walk.Angles = make([][2]float64, len(walk.Angles_))
//			for i, angles_ := range walk.Angles_ {
//				angles := strings.Split(angles_, "-")
//				var err error
//				walk.Angles[i][0], err = strconv.ParseFloat(angles[0], 64)
//				if err != nil {
//					return err
//				}
//				walk.Angles[i][1], err = strconv.ParseFloat(angles[1], 64)
//				if err != nil {
//					return err
//				}
//			}
//			//			log.Println("walk.Angles:", walk.Angles)
//		}
//	}

//	return nil
//}

func loadSpiritInfos(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, &SpiritInfos)
	if err != nil {
		return err
	}

	return nil
}

func GetSpiritInfo(spiritid string) *SpiritInfo {
	spiritInfo := SpiritInfos[spiritid]
	if spiritInfo == nil {
		log.Println("not defined spirit info, spirit id is " + spiritid)
		return nil
	}
	return spiritInfo
}

func GetSpiritFeatures(spiritInfo *SpiritInfo, featureId string) []*Feature {
	features := spiritInfo.Features[featureId]
	if features == nil {
		log.Println("not defined spirit feature info, feature id is " + featureId)
		return nil
	}
	return features
}
