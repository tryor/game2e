package game2e

import (
	. "github.com/tryor/eui"
	"github.com/tryor/gdiplus"
)

const DefaultClockGeneratorPeriod = 10       //10毫秒
const DefaultRenderClockGeneratorPeriod = 10 //10毫秒, 100/FPS = 1000/10 = 10, //16毫秒,  60/FPS = 1000/60 = 16, //33毫秒,  30/FPS = 1000/30 = 33.333, //83毫秒,  12/FPS = 1000/12 = 83.333

//组件基本信息
type ResFeatureInfo struct {
	Bounds      Rect      `json:"bounds" xml:"bounds"`
	Border      bool      `json:"border" xml:"border"`
	BorderWidth float32   `json:"borderwidth" xml:"borderwidth"`
	BorderColor RGBA      `json:"bordercolor" xml:"bordercolor"`
	Fill        bool      `json:"fill" xml:"fill"`
	FillColor   RGBA      `json:"fillcolor" xml:"fillcolor"`
	Alignment   Alignment `json:"alignment" xml:"alignment"`
	Anchorpoint PointF    `json:"anchorpoint" xml:"anchorpoint"`
	OrderZ      int       `json:"orderz" xml:"orderz"` //Z位置"
}

//文本标签配置信息
type ResFontInfo struct {
	//"name":{"text":"英雄", "font":"Arial", "size":11, "color":{"r":255, "g":0, "b":255, "a":255},
	//"style":0, "format":4096, "textalignment":{"h":1,"v":1}, "bounds":{"x":0,"y":0,"w":85,"h":67},
	//"anchorpoint":{"x":0.5,"y":0}, "alignment":{"h":1,"v":0}},
	Text          string            `json:"text" xml:"text"`
	Font          string            `json:"font" xml:"font"`
	Size          float32           `json:"size" xml:"size"`
	Color         RGBA              `json:"color" xml:"color"`
	Style         gdiplus.FontStyle `json:"style" xml:"style"`
	Multiline     bool              `json:"multiline" xml:"multiline"`
	MaxLimit      int               `json:"maxlimit" xml:"maxlimit"`
	Textalignment Alignment         `json:"textalignment" xml:"textalignment"`
	Feature       ResFeatureInfo    `json:"feature" xml:"feature"`
}

type ResElementInfo struct {
	Id      string         `json:"id" xml:"id"`
	Type    ElementType    `json:"type" xml:"type"`
	RefId   string         `json:"refid" xml:"refid"`
	RefId2  string         `json:"refid2" xml:"refid2"`
	Text    ResFontInfo    `json:"text" xml:"text"`
	Tag     string         `json:"tag" xml:"tag"`
	Desc    string         `json:"desc" xml:"desc"`
	Feature ResFeatureInfo `json:"feature" xml:"feature"`

	Points       []Point `json:"points" xml:"points"` //"points":[{"x":0,"y":0},{"x":100,"y":100}],
	EventEnabled bool    `json:"eventenabled" xml:"eventenabled"`
	Obstacle     int8    `json:"obstacle" xml:"obstacle"`
	Invisible    bool    `json:"invisible" xml:"invisible"`
	Invalid      bool    `json:"invalid" xml:"invalid"`

	//	  "normalrefid":"return_normal","//normalrefid_desc":"必须ElementTypeElement类型",
	//	  "hoveringrefid":"return_hovering","//normalrefid_desc":"必须ElementTypeElement类型",
	//	  "pressdownrefid":"return_pressdown","//normalrefid_desc":"必须ElementTypeElement类型"
	NormalRefid    string `json:"normalrefid" xml:"normalrefid"`       //仅Button有用
	HoveringRefid  string `json:"hoveringrefid" xml:"hoveringrefid"`   //仅Button有用
	PressdownRefid string `json:"pressdownrefid" xml:"pressdownrefid"` //仅Button有用
	NormalColor    RGBA   `json:"normalcolor" xml:"normalcolor"`       //仅Button有用
	HoveringColor  RGBA   `json:"hoveringcolor" xml:"hoveringcolor"`   //仅Button有用
	PressdownColor RGBA   `json:"pressdowncolor" xml:"pressdowncolor"` //仅Button有用

	Children []*ResElementInfo `json:"children" xml:"children"`
}
