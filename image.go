package game2e

import (
	"errors"
	"image"

	"github.com/tryor/game2e/log"

	. "github.com/tryor/eui"
)

type ITexture interface {
	IWidget
	Load() error
	GetImageWidth() int
	GetImageHeight() int
}

type Image struct {
	*Widget
	filename      string
	NativeImage   IBitmap
	VisibleRegion image.Rectangle
}

func NewImage(x, y, w, h int, filename string) *Image {
	img := &Image{Widget: NewWidget()}
	img.Self = img
	img.SetCoordinate(x, y)
	img.SetWidth(w)
	img.SetHeight(h)
	img.filename = filename
	img.VisibleRegion.Max.X = w
	img.VisibleRegion.Max.Y = h

	return img
}

func (this *Image) Destroy() {
	this.Widget.Destroy()
	if this.NativeImage != nil {
		this.NativeImage.Release()
		this.NativeImage = nil
	}
}

func (this *Image) GetImageWidth() int {
	if this.NativeImage != nil {
		return int(this.NativeImage.Width())
	}
	return 0
}

func (this *Image) GetImageHeight() int {
	if this.NativeImage != nil {
		return int(this.NativeImage.Height())
	}
	return 0
}

func (this *Image) Filename() string {
	return this.filename
}

func (this *Image) SetFilename(filename string) {
	if filename == "" {
		return
	}
	if this.NativeImage != nil {
		this.NativeImage.Release()
		this.NativeImage = nil
	}
	this.filename = filename
	this.SetModified(true)
}

func (this *Image) Load() error {
	if this.filename == "" {
		return errors.New("Filename is empty")
	}
	ge := this.GetGraphicsEngine()
	if ge != nil {
		img, err := ge.LoadImage(this.filename)
		if err != nil {
			return err
		} else {
			this.NativeImage = img
		}
	} else {
		return errors.New("GraphicsEngine not exist")
	}
	return nil
}

func (this *Image) CreatePath() {
	this.Widget.CreatePath()
	this.Self.CreateBoundRect()
}

func (this *Image) Draw(ge IGraphicsEngine) {
	this.Widget.Draw(ge)

	if this.NativeImage == nil {
		err := this.Self.(ITexture).Load()
		if err == nil {
			if this.Self.Width() == 0 {
				this.Self.SetWidth(this.GetImageWidth())
			}
			if this.Self.Height() == 0 {
				this.Self.SetHeight(this.GetImageHeight())
			}
			vrw, vrh := this.VisibleRegion.Dx(), this.VisibleRegion.Dy()
			if vrw == 0 {
				this.VisibleRegion.Max.X = this.VisibleRegion.Min.X + this.Self.Width()
			}
			if vrh == 0 {
				this.VisibleRegion.Max.Y = this.VisibleRegion.Min.Y + this.Self.Height()
			}

		} else {
			log.Errorf("Load NativeImage error, this.filename is %v, error:%v", this.filename, err)
		}
	}

	if this.NativeImage != nil {
		//thisX, thisY := this.Self.(IElement).GetWorldCoordinate()
		vr := this.Self.GetBoundRect()
		thisX, thisY := vr.Min.X, vr.Min.Y
		ge.DrawImage(this.NativeImage, thisX, thisY, this.VisibleRegion.Min.X, this.VisibleRegion.Min.Y, this.VisibleRegion.Dx(), this.VisibleRegion.Dy())
	} else {
		log.Error("NativeImage not exist ", this.filename)
	}
}

type ImageByResource struct {
	*Widget
	ImgRes *ImageResDict
	Cached bool
	id     string
	id2    string
}

//@see ImageResData, ImageResDict
func NewImageByResource(x, y int, id, id2 string, imageResDatas ...ImageResDataMap) *ImageByResource {
	var imgresdata *ImageResData
	if len(imageResDatas) > 0 {
		imgresdata = imageResDatas[0][id]
	} else {
		imgresdata = ImageResDatas[id]
	}

	//	var resdict *ImageResDict
	//	if imgresdata != nil {
	if imgresdata == nil {
		log.Error("resource " + id + " not find!")
		return nil
	}
	resdict := imgresdata.ResDicts[id2]
	if resdict == nil {
		log.Error("resource " + id + "." + id2 + " not find!")
		return nil
	}

	if resdict.CachedImage == nil {
		log.Error("CachedImage not exist, resource is " + id + "." + id2)
		return nil
	}
	//	}

	img := &ImageByResource{Widget: NewWidget(), Cached: true, id: id, id2: id2}
	img.Self = img
	img.SetCoordinate(x, y)

	if resdict != nil {
		img.SetWidth(resdict.Frame.W)
		img.SetHeight(resdict.Frame.H)
		img.ImgRes = resdict
	}
	return img
}

func (this *ImageByResource) Load() error {
	//LoadImageResData(ge IGraphicsEngine, filename, id2 string, respaths ...string) error
	if this.ImgRes != nil {
		return nil
	}
	imgresdata := ImageResDatas[this.id]
	if imgresdata == nil {
		ge := this.GetGraphicsEngine()
		if ge == nil {
			return errors.New("GraphicsEngine not exist")
		}

		err := LoadImageResData(ge, this.id)
		if err != nil {
			return err
		}
		imgresdata = ImageResDatas[this.id]
		if imgresdata == nil {
			return errors.New("resource " + this.id + " not find!")
		}
	}

	resdict := imgresdata.ResDicts[this.id2]
	if resdict == nil {
		return errors.New("resource " + this.id + "." + this.id2 + " not find!")
	}

	if resdict.CachedImage == nil {
		return errors.New("CachedImage not exist, resource is " + this.id + "." + this.id2)
	}

	this.SetWidth(resdict.Frame.W)
	this.SetHeight(resdict.Frame.H)
	this.ImgRes = resdict

	return nil
}

func (this *ImageByResource) GetImageWidth() int {
	if this.ImgRes != nil {
		return this.ImgRes.Frame.W
	} else {
		return 0
	}
}

func (this *ImageByResource) GetImageHeight() int {
	if this.ImgRes != nil {
		return this.ImgRes.Frame.H
	} else {
		return 0
	}
}

func (this *ImageByResource) CreatePath() {
	this.Widget.CreatePath()
	if this.ImgRes == nil {
		this.Self.(ITexture).Load()
	}
	this.Self.CreateBoundRect()
}

func (this *ImageByResource) Draw(ge IGraphicsEngine) {
	this.Widget.Draw(ge)

	if this.ImgRes == nil {
		log.Error("ImageResDict not Loaded! resource id:" + this.id + ", id2:" + this.id2)
		return
	}

	vr := this.Self.GetBoundRect()
	thisX, thisY := vr.Min.X, vr.Min.Y
	if this.Cached {
		if this.ImgRes.CachedImage != nil {
			//thisX, thisY := this.Self.(IElement).GetWorldCoordinate()
			ge.DrawImage(this.ImgRes.CachedImage, thisX, thisY, 0, 0, 0, 0)
		} else {
			log.Error("CachedImage not exist")
		}
	} else {
		if this.ImgRes.Image != nil {
			//self := this.Self.(IElement)
			//thisX, thisY := self.GetWorldCoordinate()
			ge.DrawImage(this.ImgRes.Image, thisX, thisY, 0, 0, vr.Dx(), vr.Dy())
		} else {
			log.Error("Image not exist")
		}
	}

}
