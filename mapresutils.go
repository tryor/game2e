package game2e

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"tryor/game2e/log"

	. "github.com/tryor/eui"
)

/*
{
  "id":"map$new_people_thorpe_demo",
  "title":"新人村",
  "desc":"进入游戏后的第一个地图场景",
  "width":4000,
  "height":4000,
  "layers":[
    {
      "type":1, "//layer_type_def":"1:背景层, 2:活动层,即操作层",
      "elements":[
        {
          "type":1, "//type_desc:":"1:图片,2:动画,3:精灵",
          "id":"aaaa",
          "id2":"bbbb","//id2_desc":"如果type is 1",
          "x":0,"y":0,
          "points":[{"x":0,"y":0}],"//points_desc":"点坐标值是相对于x,y位置的值",
          "eventenabled":true, "//eventenabled_desc":"true or false, 是否响应事件",
          "border":false,
          "fill":false,
          "obstacle":false
        },
        {
          "type":2,
          "id":"aaaa",
          "id2":"bbbb",
          "x":100,"y":0,
          "points":[{"x":0,"y":0},{"x":100,"y":100}],
          "eventenabled":true,
          "border":false,
          "fill":false
        },
        {
          "type":7,
          "id":"cccc",
          "id2":"",
          "x":300,"y":300,
          "points":[{"x":0,"y":0},{"x":200,"y":200}],
          "eventenabled":true,
          "border":false,
          "fill":false
        },
        {
          "type":5,
          "id":"cccc",
          "id2":"",
          "x":500,"y":500,
          "points":[{"x":0,"y":0},{"x":200,"y":200},{"x":300,"y":200},{"x":400,"y":300}],
          "eventenabled":true,
          "border":false,
          "fill":false
        }

      ]
    }
  ]


}
*/

func init() {
	MapInfos = make(map[string]*MapInfo)
}

var MapInfos map[string]*MapInfo

type MapInfo struct {
	Id     string      `json:"id" xml:"id"`
	Title  string      `json:"title" xml:"title"`
	Desc   string      `json:"desc" xml:"desc"`
	Width  int         `json:"width" xml:"width"`
	Height int         `json:"height" xml:"height"`
	Layers []*MapLayer `json:"layers" xml:"layers"`
}

type MapLayer struct {
	Id   string    `json:"id" xml:"id"`
	Tag  string    `json:"tag" xml:"tag"`
	Desc string    `json:"desc" xml:"desc"`
	Type LayerType `json:"type" xml:"type"`
	//	Feature FeatureInfo      `json:"feature" xml:"feature"`

	Alignment   Alignment `json:"alignment" xml:"alignment"`
	Anchorpoint PointF    `json:"anchorpoint" xml:"anchorpoint"`
	X           int       `json:"x" xml:"x"` //"层在地图中的位置",
	Y           int       `json:"y" xml:"y"` //"层在地图中的位置",
	Width       int       `json:"width" xml:"width"`
	Height      int       `json:"height" xml:"height"`

	VX          int               `json:"vx" xml:"vx"`           //"可视区域坐标",
	VY          int               `json:"vy" xml:"vy"`           //"可视区域坐标",
	VWidth      int               `json:"vwidth" xml:"vwidth"`   //可视区域宽度高度
	VHeight     int               `json:"vheight" xml:"vheight"` //可视区域宽度高度
	Background  RGBA              `json:"background" xml:"background"`
	Invisible   bool              `json:"invisible" xml:"invisible"`
	ScrollrateX float32           `json:"scrollrate_x" xml:"scrollrate_x"` //滚动速率，值为0-1之间, 相对于活动层",
	ScrollrateY float32           `json:"scrollrate_y" xml:"scrollrate_y"` //滚动速率，值为0-1之间, 相对于活动层",
	ScrollMode  ScrollMode        `json:"scrollmode" xml:"scrollmode"`
	DrawMode    DrawMode          `json:"drawmode" xml:"drawmode"`
	Children    []*ResElementInfo `json:"children" xml:"children"`
}

//type MapElement struct {
//	Id      string             `json:"id" xml:"id"`
//	Type    ElementType `json:"type" xml:"type"`
//	RefId   string             `json:"refid" xml:"refid"`
//	RefId2  string             `json:"refid2" xml:"refid2"`
//	Text    LabelInfo          `json:"text" xml:"text"`
//	Tag     string             `json:"tag" xml:"tag"`
//	Desc    string             `json:"desc" xml:"desc"`
//	Feature FeatureInfo        `json:"feature" xml:"feature"`

//	Points       []Point `json:"points" xml:"points"` //"points":[{"x":0,"y":0},{"x":100,"y":100}],
//	EventEnabled bool           `json:"eventenabled" xml:"eventenabled"`
//	Obstacle     int            `json:"obstacle" xml:"obstacle"`
//	Invisible    bool           `json:"invisible" xml:"invisible"`
//	Invalid      bool           `json:"invalid" xml:"invalid"`
//	Elements     []*MapElement  `json:"elements" xml:"elements"`
//}

func LoadMapInfos(filenames ...string) error {
	for _, filename := range filenames {
		if err := LoadMapInfo(filename); err != nil {
			return err
		}
	}

	//	log.Info("MapInfos:", MapInfos)
	//	for mapkey, mapInfo := range MapInfos {
	//		log.Info("LoadMapInfos:", mapkey, mapInfo)
	//		for _, maplayer := range mapInfo.Layers {
	//			log.Info("maplayer:", maplayer)
	//			for _, mapelement := range maplayer.Elements {
	//				log.Info("mapelement:", mapelement)
	//			}
	//		}
	//	}

	return nil
}

func LoadMapInfo(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Error(err)
		return err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(err)
		return err
	}

	mapInfo := &MapInfo{}
	err = json.Unmarshal(buf, mapInfo)
	if err != nil {
		log.Error(err)
		return err
	}
	if MapInfos[mapInfo.Id] != nil {
		err = errors.New("map is exist, filename is " + filename + ", id is " + mapInfo.Id)
		log.Error(err)
		return err
	}
	MapInfos[mapInfo.Id] = mapInfo

	return nil
}
