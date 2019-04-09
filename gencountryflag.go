// +build ignore

package main

import (
    "os"
    "path/filepath"
    "io/ioutil"
    "image/png"
    "github.com/disintegration/imaging"
)

func main() {
    rootPath := "png1000px"

    files,err := ioutil.ReadDir(rootPath)
    if err != nil {
        panic(err)
    }

    for _,fe := range files {
        fpath := filepath.Join(rootPath, fe.Name())
        func() {
            f,err := os.Open(fpath)
            if err != nil {
                panic(err)
            }
            img,err := png.Decode(f)
            if err != nil {
                panic(err)
            }
            f.Close()

            w,h := img.Bounds().Max.X, img.Bounds().Max.Y
            if w > (h*5/4) {
                img = imaging.Resize(img, h*5/4, h, imaging.Lanczos)
            }
            img = imaging.Fill(img, 200, 200, imaging.Center, imaging.Lanczos)

            f,err = os.Create(fpath)
            if err != nil {
                panic(err)
            }
            if err := png.Encode(f, img); err != nil {
                panic(err)
            }
            f.Close()
        }()
    }
}
