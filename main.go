package main

import (
	"fmt"
	"github.com/disintegration/imaging"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type conf struct {
	Slug    string   `yaml:"slug"`
	BaseDir string   `yaml:"base_dir"`
	Feature string   `yaml:"feature"`
	Ignore  []string `yaml:"ignore"`
}

type item struct {
	Src      string
	Dest     string
	Itemtype string
	Limit    int
}

func main() {
	configFile := os.Args[1]
	var c conf
	c.getConf(configFile)
	inputDir := getInputDir(configFile)
	imageOutputDir := path.Join(c.BaseDir, "static/images", c.Slug)
	featureOutputDir := path.Join(c.BaseDir, "static/images/feature_images")
	audioOutputDir := path.Join(c.BaseDir, "static/audio", c.Slug)
	otherOutputDir := path.Join(c.BaseDir, "static/other", c.Slug)
	pageOutputDir := path.Join(c.BaseDir, "content")
	fmt.Println("Input dir:", inputDir)
	fmt.Println("Image output dir:", imageOutputDir)
	fmt.Println("Feature output dir:", featureOutputDir)
	fmt.Println("Audio output dir:", audioOutputDir)
	fmt.Println("Other output dir:", otherOutputDir)
	fmt.Println("Page output dir:", pageOutputDir)

	workList := c.getWorkList(inputDir)
	fmt.Println(workList)
	os.MkdirAll(imageOutputDir, os.ModePerm)
	getFiles(inputDir, imageOutputDir)

}

func (c *conf) getConf(cf string) *conf {
	yamlFile, err := ioutil.ReadFile(cf)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func getInputDir(cf string) string {
	p := path.Dir(cf)
	if path.IsAbs(p) {
		return p
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cant find current working directory %v", err)
	}
	return path.Join(wd, p)
}

func (c conf) getWorkList(inputDir string) []item {
	var workList []item
	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		ext := path.Ext(file.Name())
		switch ext {
		case ".jpg":
			workList = append(workList, item{path.Join(inputDir, file.Name()), "asd", "jpg", 123})
		case ".png":
			workList = append(workList, item{path.Join(inputDir, file.Name()), "asd", "png", 123})
		default:
			workList = append(workList, item{path.Join(inputDir, file.Name()), "asd", "other", 0})
		}
	}
	return workList
}

func getFiles(inputDir, outputDir string) []os.FileInfo {
	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if path.Ext(file.Name()) == ".jpg" {
			img, err := imaging.Open(path.Join(inputDir, file.Name()))
			if err != nil {
				panic(err)
			}
			img = imaging.Fit(img, 350, 350, imaging.Lanczos)
			err = imaging.Save(img, path.Join(outputDir, file.Name()))
			if err != nil {
				panic(err)
			}
		}
	}
	return files
}
