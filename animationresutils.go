package game2e

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func init() {
	AnimationInfos = make(map[string]*AnimationInfo)
}

/*
{
"$legendaryswordsman_1":{
    metadata:{},
    frames:[
        {name:"hero/$legendaryswordsman_1.png",id:"1",delay:200},
        {name:"hero/$legendaryswordsman_1.png",id:"2",delay:200},
        {name:"hero/$legendaryswordsman_1.png",id:"3",delay:200}
    ]
},

"$legendaryswordsman_forward":{}

}
*/

type AnimationMetadata struct {
}

type AnimationFrame struct {
	Name  string `json:"name" xml:"name"`
	Id    string `json:"id" xml:"id"`
	Delay int    `json:"delay" xml:"delay"`

	Feature ResFeatureInfo `json:"feature" xml:"feature"`
	//	Alignment   Alignment `json:"alignment" xml:"alignment"`
	//	Anchorpoint PointF    `json:"anchorpoint" xml:"anchorpoint"`
}

type AnimationInfo struct {
	Width    int               `json:"width" xml:"width"`
	Height   int               `json:"height" xml:"height"`
	Metadata AnimationMetadata `json:"metadata" xml:"metadata"`
	Frames   []AnimationFrame  `json:"frames" xml:"frames"`
}

var AnimationInfos map[string]*AnimationInfo

func LoadAnimationInfos(filenames ...string) error {
	for _, filename := range filenames {
		err := loadAnimationInfos(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadAnimationInfos(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, &AnimationInfos)
	if err != nil {
		return err
	}

	//	for key, val := range AnimationInfos {
	//		log.Info("AnimationInfos.key:", key, ", val:", val)
	//	}

	return nil
}
