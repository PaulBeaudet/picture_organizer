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
    destinationPointer := flag.String("dest", home + "/Pictures/", "Destination of sorted pictures")
    safemode := flag.Bool("safemode", true, "Keeps a copy of sorted photos in source directory")
    eventName := flag.String("name", "", "Adds an event name in folder hiarchy")
    daysIntoPast := flag.Int("retro", 0, "How many days into past to retrospectively sort photos")
    flag.Parse() // get flags that were passed to app
    // fmt.Println("src:", *sourcePointer, " dest:",*destinationPointer," sm:", *safemode, " evn:", *eventName, " retro:", *daysIntoPast)
    scanAndMove(*sourcePointer, *destinationPointer, *safemode, *eventName, *daysIntoPast)
}

func scanAndMove(src string, dest string, safemode bool, eventName string, daysIntoPast int){
    // TODO maybe stat source and dest to make sure they exist
    uploads, err := os.Open(src)
    if err != nil{panic(err)}
    files, err := uploads.Readdir(-1)
    uploads.Close()
    if err != nil {panic(err)}
    now := time.Now()
    for _, file := range files {
        fileName := file.Name()
        currentLocation := src + fileName
        taken, isPhoto := timeTakenIfPhoto(currentLocation)
        if !isPhoto{continue} // skip files without exif
        daysSinceTaken := int(now.Sub(taken).Hours()/24)
        if daysIntoPast > 0 && daysSinceTaken > daysIntoPast{continue}
        fileName = strings.ToLower(fileName) // convert to lower case
        hiarchy := taken.Format("2006") + "/" + taken.Format("01_02_") + eventName + "/"
        nextDest := dest + hiarchy
        mkdir(nextDest)
        newName := getValidName(nextDest, taken.Format("15_04_05"), fileName)
        copyFile(currentLocation, nextDest + newName)
        if safemode {
            duplicateDest := src + hiarchy
            mkdir(duplicateDest) //issue if searching folders w/ previously state in same format
            moveFile(currentLocation, duplicateDest + newName);
        } else { rm(currentLocation) } // otherwise remove original
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
    if err != nil {return time.Time{}, false}
    exifData, err := exif.Decode(file)
    if err != nil {return time.Time{}, false}
    taken, err := exifData.DateTime()
    if err != nil {return time.Time{}, false}
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
    if err != nil {fmt.Println("Couldn't open source file: " + err.Error())}
    outputFile, err := os.Create(dest)
    if err != nil {
        inputFile.Close()
        fmt.Println("Couldn't open dest file: " + err.Error())
    }
    defer outputFile.Close() // do this when function exits
    _, err = io.Copy(outputFile, inputFile)
    inputFile.Close()
    if err != nil {fmt.Println("Writing to output file failed: " + err.Error())}
}

func moveFile(oldPath string, newPath string){
    if err := os.Rename(oldPath, newPath); err != nil{panic(err)}
}

func rm(path string){
    if err := os.Remove(path); err != nil {fmt.Println("Could not remove " + path)}
}
