// upload.go ~ organizes camera uploads into main photo storage folder
// Copyright 2020 Paul Beaudet ~ MIT License
package main

import (
    "fmt"
    "os"
    "io"
    "strings"
    "flag"
    "math/rand"
    "strconv"
    "time"
    "path/filepath"

    "github.com/rwcarlsen/goexif/exif"
)

func main(){
    workingDir, _ := os.Getwd() // get the working directory
    home, _ := os.UserHomeDir() // get $HOME // Defaults will at least work in linux
    sourcePointer := flag.String("src", workingDir + "/", "Source of pictures to sort")
    // Another primary example of src would be $HOME/Dropbox/Camera Uploads
    destinationPointer := flag.String("dest", home + "/Pictures/", "Destination of sorted pictures")
    safemode := flag.Bool("safemode", true, "Keeps a copy of sorted photos in source directory")
    flag.Parse() // get flags that were passed to app
    fmt.Println("Source= " + *sourcePointer + " -> destination= " + *destinationPointer + " Safemode:", *safemode)
    scanAndMove(*sourcePointer, *destinationPointer, *safemode)
}

func scanAndMove(src string, dest string, safemode bool){
    uploads, uErr := os.Open(src)
    if uErr != nil{panic(uErr)}
    files, error := uploads.Readdir(-1)
    uploads.Close()
    if error != nil {panic(error)}
    for _, file := range files {
        fileName := file.Name()
        currentLocation := src + fileName
        taken, isPhoto := timeTakenIfPhoto(currentLocation)
        if isPhoto { // given we are getting a time back from photo w/exif
            fileName = strings.ToLower(fileName) // convert to lower case
            hiarchy := taken.Format("2006") + "/" + taken.Format("01_02_") + "/"
            nextDest := dest + hiarchy
            mkdir(nextDest)
            newName := getValidName(nextDest, taken.Format("15_04_05"), fileName)
            copyFile(currentLocation, nextDest + newName)
            if safemode {
                duplicateDest := src + hiarchy
                mkdir(duplicateDest) //issue if searching folders w/ previously state in same format
                moveFile(currentLocation, duplicateDest + newName);
            } else { rm(currentLocation) } // otherwise remove original
        } // otherwise this is not a photo timeTakenIfPhoto logs out
    }
}

func getValidName(inPath string, newName string, orgName string)(string){
    ext := filepath.Ext(orgName); // get current extention name, e.g. .jpg or .rw2
    // if ext == "" { ext = ".jpg" } // fix past mistake of not including extention
    fullPath := inPath + newName + ext
    _, err := os.Stat(fullPath)
    if os.IsNotExist(err) { // ideally this is a new file in case just do what we we're thinking
        return newName + ext
    } else { // TODO this could cause an infinate loop in cases of +100 duplicates
        psudoRand := strconv.Itoa(rand.Intn(99))
        return getValidName(inPath, newName + "_" + psudoRand , ext)
    }
}

func timeTakenIfPhoto(photoPath string)(time.Time, bool){
    file, err := os.Open(photoPath)
    if err != nil {panic(err)}
    exifData, exifErr := exif.Decode(file)
    if exifErr != nil {
        fmt.Println(photoPath + ": is not a photo with exif data")
        return time.Time{}, false
    }
    taken, dateTimeErr := exifData.DateTime()
    if dateTimeErr != nil {
        fmt.Println(photoPath + ": could not get date and time")
        return time.Time{}, false
    }
    return taken, true
}

func mkdir(dirToCreate string){
    _, error := os.Stat(dirToCreate)
    if os.IsNotExist(error){
        if err := os.MkdirAll(dirToCreate, 0755); err != nil {panic(err)}
    }
}

func copyFile(src string, dest string){
    inputFile, err := os.Open(src)
    if err != nil {fmt.Println("Couldn't open source file: %s", err)}
    outputFile, err := os.Create(dest)
    if err != nil {
        inputFile.Close()
        fmt.Println("Couldn't open dest file: %s", err)
    }
    defer outputFile.Close() // do this when function exits
    _, err = io.Copy(outputFile, inputFile)
    inputFile.Close()
    if err != nil {fmt.Println("Writing to output file failed: %s", err)}
}

func moveFile(oldPath string, newPath string){
    if err := os.Rename(oldPath, newPath); err != nil{panic(err)}
}

func rm(path string){
    if err := os.Remove(path); err != nil {fmt.Println("Could not remove " + path)}
}
