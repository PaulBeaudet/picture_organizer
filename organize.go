// upload.go ~ organizes camera uploads into main photo storage folder
// Copyright 2020 Paul Beaudet ~ MIT License
package main

import (
    "fmt"
    "os"
    "strings"
    "flag"
    "github.com/rwcarlsen/goexif/exif"
)

func main(){
    sourcePointer := flag.String("src", "$HOME/Downloads", "Source of pictures to sort")
    // Another primary example of src would be $HOME/Dropbox/Camera Uploads
    destinationPointer := flag.String("dest", "$HOME/Pictures", "Destination of sorted pictures")
    flag.Parse()
    // Or maybe a location like $HOME/Dropbox/pictures
    fmt.Println("Source= " + *sourcePointer)
    fmt.Println("destination= " + *destinationPointer)
    scanAndMove(*sourcePointer, *destinationPointer)
}

func scanAndMove(uploadDir string, photoDir string){
    uploads, uErr := os.Open(uploadDir)
    if uErr != nil{panic(uErr)}
    files, error := uploads.Readdir(-1)
    uploads.Close()
    if error != nil{panic(error)}
    for _, file := range files {
        if ext := isPhotoExtention(file.Name()); ext != ""{
            moveAndRename(file, uploadDir, photoDir, ext)
        } else {
            fmt.Println(file.Name() + " is not a photo");
        }
    }
}

func moveAndRename(file os.FileInfo, sourceDir string, destDir string, ext string){
    currentLocation := sourceDir + file.Name()
    monthDay, year, hourMinutes := timeTaken(currentLocation)
    nextDest := destDir + year + "/" + monthDay
    mkdir(nextDest)
    moveFile(currentLocation, nextDest + "/" + hourMinutes + ext)
}

func isPhotoExtention(fileName string)(ext string){
    fileName = strings.ToLower(fileName)
    if strings.HasSuffix(fileName, ".rw2"){
        return ".rw2"
    } else if strings.HasSuffix(fileName, ".jpg"){
        return ".jpg"
    } else {
        return ""
    }
}

func timeTaken(photoPath string)(mmdd string, yyyy string, hhmm string){
    file, err := os.Open(photoPath)
    if err != nil {panic(err)}
    exifData, xerr := exif.Decode(file)
    if xerr != nil {panic(xerr)}
    taken, terr := exifData.DateTime()
    if terr != nil {panic(terr)}
    monthDay := taken.Format("01_02_")
    year := taken.Format("2006")
    hourMinutes := taken.Format("15_04_05")
    fmt.Println(photoPath + "\n ---> Taken on month and day: " + monthDay +
        " in year " + year + " at exactly " + hourMinutes)
    return monthDay, year, hourMinutes
}

func getCreation(pathToFile string)(mmdd string, yyyy string, hhmm string){
    fileStat, err := os.Stat(pathToFile)
    if err != nil {panic(err)}
    timeObj := fileStat.ModTime()
    monthDay := timeObj.Format("01_02_")
    year := timeObj.Format("2006")
    hourMinutes := timeObj.Format("15_04_05")
    return monthDay, year, hourMinutes
}

func mkdir(dirToCreate string){
    _, error := os.Stat(dirToCreate)
    if os.IsNotExist(error){
        if err := os.MkdirAll(dirToCreate, 0755); err != nil {panic(err)}
    }
}

func moveFile(oldPath string, newPath string){
    if err := os.Rename(oldPath, newPath); err != nil{panic(err)}
}

func ls(dir string){
    f, err := os.Open(dir)
    if err != nil {panic(err)}
    files, error := f.Readdir(-1)
    f.Close()
    if error != nil {panic(error)}
    fmt.Println("Contents of: " + dir);
    for _, file := range files {
        fmt.Println(file.Name())
    }
}
