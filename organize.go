// upload.go ~ organizes camera uploads into main photo storage folder
// Copyright 2020 Paul Beaudet ~ MIT License
package main

import (
    "fmt"
    "os"
    "strings"
    "flag"
    "math/rand"
    "strconv"

    "github.com/rwcarlsen/goexif/exif"
)

func main(){
    home, _ := os.UserHomeDir() // get $HOME // Defaults will at least work in linux
    sourcePointer := flag.String("src", home + "/Downloads/", "Source of pictures to sort")
    // Another primary example of src would be $HOME/Dropbox/Camera Uploads
    destinationPointer := flag.String("dest", home + "/Pictures/", "Destination of sorted pictures")
    // Or maybe a location like $HOME/Dropbox/pictures
    flag.Parse() // get flags that were passed to app
    fmt.Println("Source= " + *sourcePointer + " -> destination= " + *destinationPointer)
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
        }
    }
}

func checkForDuplicate(filePathRename string, extention string)(okayFilePathName string){
    intendedFileName := filePathRename + extention
    _, err := os.Stat(intendedFileName)
    if os.IsNotExist(err) { // ideally this is a new file in case just do what we we're thinking
        return intendedFileName
    } else { // TODO this could cause an infinate loop in cases of +100 duplicates
        psudoRand := strconv.Itoa(rand.Intn(99))
        return checkForDuplicate(filePathRename + "_" + psudoRand , extention)
    }
}

func moveAndRename(file os.FileInfo, sourceDir string, destDir string, ext string){
    currentLocation := sourceDir + file.Name()
    monthDay, year, hourMinutes := timeTaken(currentLocation)
    nextDest := destDir + year + "/" + monthDay
    mkdir(nextDest)
    newFilePathName := checkForDuplicate(nextDest + "/" + hourMinutes, ext)
    moveFile(currentLocation, newFilePathName)
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
