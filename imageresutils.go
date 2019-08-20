package game2e

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"tryor/game2e/log"

	. "github.com/tryor/eui"
)

func init() {
	ImageResDicts = make(map[string]ImageResDictMap)
	ImageResDatas = make(ImageResDataMap)
}

//{
//"$legendaryswordsman_1.png":
//{
//	"1": {"frame": {"x":0,"y":0,"w":72,"h":67}}
//	"2": {"frame": {"x":72,"y":0,"w":72,"h":67}}
//	"3": {"frame": {"x":144,"y":0,"w":72,"h":67}}
//},

type Frame struct {
	X int `json:"x" xml:"x"`
	Y int `json:"y" xml:"y"`
	W int `json:"w" xml:"w"`
	H int `json:"h" xml:"h"`
}

type ImageResDict struct {
	Frame       Frame `json:"frame" xml:"frame"`
	Rotated     bool  `json:"rotated" xml:"rotated"`
	Trimmed     bool  `json:"trimmed" xml:"trimmed"`
	FlipX       bool  `json:"flipx" xml:"flipx"` //水平镜像
	FlipY       bool  `json:"flipy" xml:"flipy"` //垂直镜像
	Image       IBitmap
	CachedImage ICachedBitmap
}

type ImageResData struct {
	Src      IBitmap
	ResDicts map[string]*ImageResDict
}

type ImageResDataMap map[string]*ImageResData
type ImageResDictMap map[string]*ImageResDict

var ImageResDicts map[string]ImageResDictMap
var ImageResDatas ImageResDataMap

func LoadImageResDicts(resfilenames ...string) error {
	for _, name := range resfilenames {
		err := loadImageResDict(name)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadImageResDict(resfilename string) error {
	file, err := os.Open(resfilename)
	if err != nil {
		return err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, &ImageResDicts)
	if err != nil {
		return err
	}

	//	for key, dicts := range ImageResDicts {
	//		println("ImageResDicts.key:", key)
	//		for k1, dict := range dicts {
	//			println("ImageResDicts.k1:", k1, dict)
	//		}
	//	}

	return nil
}

func loadImageData(ge IGraphicsEngine, respath, filename string, dictmap ImageResDictMap) error {
	src, err := ge.LoadImage(respath + filename)
	if err != nil {
		log.Error(err)
		return err
	} else {
		imgresdata := &ImageResData{Src: src, ResDicts: dictmap}
		for _, dict := range dictmap {
			//log.Info("dict.FlipX:", dict.FlipX)
			dict.Image = ge.CloneImage(src, dict.Frame.X, dict.Frame.Y, dict.Frame.W, dict.Frame.H)
			if dict.FlipX {
				dict.Image.RotateFlipX()
			}
			if dict.FlipY {
				dict.Image.RotateFlipY()
			}
			if dict.Image != nil {
				dict.CachedImage, err = ge.CacheImage(dict.Image)
				if err != nil {
					log.Error(err)
				}
			}
		}
		if ImageResDatas[filename] != nil {
			log.Warn("duplicate resource names!")
		}
		ImageResDatas[filename] = imgresdata
	}
	return nil
}

func LoadImageResData(ge IGraphicsEngine, filename string, respaths ...string) error {
	respath := Resource.ImageResRoot
	if len(respaths) > 0 && respaths[0] != "" {
		respath = respaths[0]

	}
	dictmap, ok := ImageResDicts[filename]
	if !ok {
		err := errors.New("resource " + filename + " not find!")
		log.Error(err)
		return err
	}
	return loadImageData(ge, respath, filename, dictmap)
}

func LoadImageResDatas(ge IGraphicsEngine, after func(finish bool), respaths ...string) error {
	respath := Resource.ImageResRoot
	if len(respaths) > 0 && respaths[0] != "" {
		respath = respaths[0]
	}
	for filename, dictmap := range ImageResDicts {
		_, ok := ImageResDatas[filename]
		if !ok {
			loadImageData(ge, respath, filename, dictmap)
			//			src, err := ge.LoadImage(respath + filename)
			//			if err != nil {
			//				log.Error(err)
			//			} else {
			//				imgresdata := &ImageResData{Src: src, ResDicts: dictmap}
			//				for _, dict := range dictmap {
			//					dict.Image = ge.CloneImage(src, dict.Frame.X, dict.Frame.Y, dict.Frame.W, dict.Frame.H)
			//					if dict.Image != nil {
			//						dict.CachedImage, err = ge.CacheImage(dict.Image)
			//						if err != nil {
			//							log.Error(err)
			//						}
			//					}
			//				}
			//				ImageResDatas[filename] = imgresdata
			//			}
		}
		if after != nil {
			after(false)
		}
	}

	if after != nil {
		after(true)
	}

	//	for key, data := range ImageResDatas {
	//		println("ImageResDatas.key:", key, data, data.Src)
	//		for k2, dict := range data.ResDicts {
	//			println("ImageResDatas.k2:", k2, dict.CachedImage, dict.Image)
	//		}
	//	}

	//	fmt.Println("LoadImageResData.ImageResDatas:", ImageResDatas["hero/$legendaryswordsman_1.png"].ResDicts["1"])
	return nil
}
