package main

import (
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type item struct {
	Src   string
	Dest  string
	Limit int
}

func main() {
	if len(os.Args) > 1 {
		viper.SetConfigFile(os.Args[1])
	} else {
		viper.SetConfigName("gmbconfig") // name of config file (without extension)
		viper.AddConfigPath(".")         // optionally look for config in the working directory
	}
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatalf("fatal error config file: %s", err)
	}

	viper.SetDefault("img_limit", 800)
	viper.SetDefault("menu_limit", 900)
	viper.SetDefault("feature_limit", 1280)

	log.Println("Config", viper.AllKeys())

	workList := getWorkList()

	processWorkList(workList)

	writePost(workList)
}

func writeFrontMatter(f *os.File, workList []item) {
	t := time.Now()

	f.WriteString("---\n")
	f.WriteString("title: Title placeholder\n")
	f.WriteString("date: " + t.Format("2006-01-02") + "\n")
	f.WriteString("description: Description placeholder\n")
	f.WriteString("toplevel: True\n")
	f.WriteString("draft: True\n")
	//f.WriteString("image_feature: " + page_feature + "\n")
	//f.WriteString("image_menu: " + page_menu + "\n")
	f.WriteString("---\n\n")
	f.WriteString("Text placeholder\n\n")
}

func writeImageTag(f *os.File, relPath string) {
	//f.WriteString('{{% fig-l src="/' + filename + '" %}}  {{% /fig-l %}}\n\n')
	f.WriteString("path" + relPath)
	//  #f.WriteString('{{% /fig-l %}}\n\n')
}
func writeAudioTag(f *os.File, relPath string) {
	//f.WriteString('{{% audio src="/' + filename + '" %}}\n')
	//f.WriteString('{{% /audio %}}\n\n')
}

func writeLinkTag(f *os.File, relPath string) {
	//f.WriteString('[Link title](/' + filename + ')\n\n')
}

func writePost(workList []item) {
	outputDir := path.Join(viper.GetString("base_dir"), "content")
	staticDir := path.Join(viper.GetString("base_dir"), "static")
	post := path.Join(outputDir, viper.GetString("slug")+".en.md")

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		log.Println("Creating dir", outputDir)
		os.MkdirAll(outputDir, os.ModePerm)
	}

	if _, err := os.Stat(post); os.IsNotExist(err) {
		log.Println("Creating post", post)
		f, err := os.Create(post)
		check(err)
		defer f.Close()
		writeFrontMatter(f, workList)
	}

	f, err := os.OpenFile(post, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	check(err)
	defer f.Close()

	for _, value := range workList {
		ext := path.Ext(value.Dest)
		relPath, err := filepath.Rel(staticDir, value.Dest)
		check(err)
		switch ext {
		case ".jpg", ".jpeg", ".JPG", "JPEG", ".png", ".PNG":
			writeImageTag(f, relPath)
		case ".mp3", ".MP3", ".ogg", ".OGG":
			writeAudioTag(f, relPath)
		default:
			writeLinkTag(f, relPath)
		}
	}
}

func processWorkList(workList []item) {
	done := make(chan bool)

	for _, value := range workList {
		value := value
		outputDir := path.Dir(value.Dest)

		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			log.Println("Creating dir", outputDir)
			os.MkdirAll(outputDir, os.ModePerm)
		}

		if value.Limit != 0 {
			go func() {
				processImg(value.Src, value.Dest, value.Limit)
				done <- true
			}()
			continue
		} else {
			go func() {
				copyFile(value.Src, value.Dest)
				done <- true
			}()
			continue
		}
	}

	// wait for all goroutines to complete before exiting
	for _ = range workList {
		<-done
	}
}

func copyFile(src, dest string) {
	log.Print("Copy: ", src, dest)
	from, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
}

func getWorkList() (workList []item) {
	inputDir := getInputDir(viper.ConfigFileUsed())
	baseDir := viper.GetString("base_dir")
	slug := viper.GetString("slug")
	menuLimit := viper.GetInt("menu_limit")
	featureLimit := viper.GetInt("feature_limit")
	imgLimit := viper.GetInt("img_limit")
	_, configFile := path.Split(viper.ConfigFileUsed())

	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		} else if file.Name() == configFile {
			continue
		} else if file.Name() == viper.GetString("feature") {
			workList = append(workList, item{path.Join(inputDir, file.Name()),
				path.Join(baseDir, "static/images/feature_images", "m_"+file.Name()), menuLimit})
			workList = append(workList, item{path.Join(inputDir, file.Name()),
				path.Join(baseDir, "static/images/feature_images", "f_"+file.Name()), featureLimit})
		} else {
			ext := path.Ext(file.Name())
			switch ext {
			case ".jpg", ".jpeg", ".JPG", "JPEG":
				workList = append(workList, item{path.Join(inputDir, file.Name()),
					path.Join(baseDir, "static/images", slug, file.Name()), imgLimit})
			case ".png", ".PNG":
				workList = append(workList, item{path.Join(inputDir, file.Name()),
					path.Join(baseDir, "static/images", slug, file.Name()), 0})
			default:
				workList = append(workList, item{path.Join(inputDir, file.Name()),
					path.Join(baseDir, "static/other", slug, file.Name()), 0})
			}
		}
	}
	return workList
}

func processImg(src, dest string, limit int) {
	log.Println("Img: ", src, dest)
	f, err := os.Open(src)
	defer f.Close()
	check(err)
	ot, err := getOrientation(f)
	if err != nil {
		log.Print("No Orientation tag found, move on, nothing to see here")
	}
	img, err := imaging.Open(src)
	check(err)

	switch {
	case ot == 2:
		img = imaging.FlipH(img)
	case ot == 3:
		img = imaging.Rotate180(img)
	case ot == 4:
		img = imaging.FlipH(img)
		img = imaging.Rotate180(img)
	case ot == 5:
		img = imaging.FlipV(img)
		img = imaging.Rotate270(img)
	case ot == 6:
		img = imaging.Rotate270(img)
	case ot == 7:
		img = imaging.FlipV(img)
		img = imaging.Rotate90(img)
	case ot == 8:
		img = imaging.Rotate90(img)
	}

	img = imaging.Fit(img, limit, limit, imaging.Lanczos)
	err = imaging.Save(img, dest)
	check(err)
}

func getOrientation(f *os.File) (ot int64, err error) {
	defer func() {
		if r := recover(); r != nil {
			ot = 1
		}
	}()

	x, err := exif.Decode(f)
	if err != nil {
		return
	}

	tag, err := x.Get(exif.Orientation)
	if err != nil {
		return
	}

	ot, _ = tag.Int64(0)
	return
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
