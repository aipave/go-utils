package test_example

import (
	"flag"
	"fmt"
	"image/gif"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/aipave/go-utils/gerr"
	"github.com/aipave/go-utils/glogs/glogrus"
	"github.com/aipave/go-utils/gwarn"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func IMGConvert01(url string, layoutSuffix string) error {
	if !strings.Contains(layoutSuffix, layoutSuffix) && !strings.Contains(layoutSuffix, ".jpeg") && !strings.Contains(layoutSuffix, ".png") {
		return gerr.New(0, "type error")
	}
	rsp, err := http.Get(url)
	if err != nil {
		logrus.Errorf("download fail|url=%v,err:%v", url, err)
		return gerr.New(0, "download fail")
	}
	defer rsp.Body.Close()

	cTp := rsp.Header.Get("Content-Type")
	switch {
	case strings.HasPrefix(cTp, "image/jpeg"), strings.HasPrefix(cTp, "image/jpg"),
		strings.HasPrefix(cTp, "image/png"):
		logrus.Infof("handleVipExpireEvent| no need to deal with %v", url)
		return nil
	case strings.HasPrefix(cTp, "image/gif"), strings.HasPrefix(cTp, "image/webp"):
		logrus.Infof("%v", cTp)
	default:

	}
	/*
		gif.DecodeAll: This function is used to decode a GIF image that may contain multiple frames.
			It reads the entire image data from an io.Reader and returns a pointer to a gif.GIF struct that contains the frames and other metadata.
		gif.Decode: This function is used to decode a single-frame GIF image.
			It reads the image data from an io.Reader and returns a pointer to an image.Image interface that represents the decoded image.
	*/
	gifImage, err := gif.DecodeAll(rsp.Body)
	if err != nil {
		logrus.Errorf("handleVipExpireEvent|decode, err:%v", err)
		return gerr.New(0, fmt.Sprintf("decode gif fail:err=%v", err))
	}

	// down local
	ext := filepath.Ext(url)
	fileName := strings.TrimSuffix(filepath.Base(url), ext)
	outFile, err := os.Create(fileName + layoutSuffix)
	if err != nil {
		logrus.Errorf("handleVipExpireEvent|trim fail, err:%v", err)
	}
	defer outFile.Close()

	//err = jpeg.Encode(outFile, gifImage, nil)
	//err = jpeg.Encode(outFile, gifImage.Image[rand.Int()%len(gifImage.Image)], nil)
	err = jpeg.Encode(outFile, gifImage.Image[2], nil)
	if err != nil {
		logrus.Warningf("handleVipExpireEvent|to %v fail, err:%v", layoutSuffix, err)
	}

	// Open result file in image viewer
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "start", fileName+layoutSuffix)
	} else if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		//cmd = exec.Command("xdg-open", fileName+layoutSuffix)
		// or use "open" command on macOS
		cmd = exec.Command("open", fileName+layoutSuffix)
	}
	err = cmd.Run()
	if err != nil {
		return gerr.New(0, "run error")
	}

	return gerr.New(0, fmt.Sprint(len(gifImage.Image)))
}

func IMGConvert(url string, layoutSuffix string) {

}

func TestImageGit2Jpg(t *testing.T) {
	//url := "https://img.gif8.com/g8/imgs/20201118/fa38518556770ea51bfee5ec335db7ba.gif"
	url := "https://qna.smzdm.com/201911/16/5dcf844ce3d705101.gif_fo742.jpg" // 52 frames
	err := IMGConvert(url, ".png")
	gwarn.AleterGetter("https://open.feishu.cn/open-apis/bot/v2/hook/2f1dc72c-8d2d-4641-bd95-31bbd6fcd2c7", "run error").Notice(
		fmt.Sprintf("%v", err))
	t.Log(err)

}

var imagelogConfigFile = flag.String("f", "config.yaml", "the config file for dbms test")

func GetImagelogConfig() ImagelogConfig {
	return imagelogConfig
}

var imagelogConfig ImagelogConfig

type ImagelogConfig struct {
	Log struct {
		Level string `yaml:"Level"`
		Path  string `yaml:"Path"`
		Mode  string `yaml:"Mode"`
	} `yaml:"Log"`

	Imagelog string `yaml:"Imagelog"`
}

func init() {
	data, err := ioutil.ReadFile(*imagelogConfigFile)
	if err != nil {
		logrus.Errorf("open file err %v", err)
	}
	err = yaml.Unmarshal(data, &imagelogConfig)
	if err != nil {
		logrus.Errorf("unmarshal file err|%v", err)
	}

	glogrus.Init(glogrus.WithAlertUrl("https://open.feishu.cn/open-apis/bot/v2/hook/2f1dc72c-8d2d-4641-bd95-31bbd6fcd2c7"))
	glogrus.MustSetLevel(GetImagelogConfig().Log.Level)
}
