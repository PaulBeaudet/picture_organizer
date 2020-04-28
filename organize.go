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
    if error != nil {panic(error)}
    for _, file := range files {
        if ext := isPhotoExtention(file.Name()); ext != ""{
            moveAndRename(file, uploadDir, photoDir, ext)
        }
    }
}

func moveAndRename(file os.FileInfo, sourceDir string, destDir string, ext string){
    currentLocation := sourceDir + file.Name()
    monthDay, year, hourMinutes := timeTaken(currentLocation)
    hiarchy := year + "/" + monthDay + "/"
    nextDest := destDir + hiarchy
    mkdir(nextDest) // This should actually raise a panic if dest is non-existent
    newName := checkForDuplicate(nextDest, hourMinutes, ext)
    fail := copyFile(currentLocation, nextDest + newName)
    if fail != nil {
        fmt.Println(fail) // soft fail, if copy fails keep going
    } else { // if copy is succesfull move into hiarchial directory with in source
        copyDest := sourceDir + hiarchy
        mkdir(copyDest)
        moveFile(currentLocation, copyDest + newName);
    } // This we at least have a backup | TODO: would be an issue if recursively
} // ...searching folders with a previously rendered state in this format would
//   ...cause unnecisary in place overwrites

func checkForDuplicate(inPath string, fileName string, ext string)(okFileName string){
    fullPath := inPath + fileName + ext
    _, err := os.Stat(fullPath)
    if os.IsNotExist(err) { // ideally this is a new file in case just do what we we're thinking
        return fileName + ext
    } else { // TODO this could cause an infinate loop in cases of +100 duplicates
        psudoRand := strconv.Itoa(rand.Intn(99))
        return checkForDuplicate(inPath, fileName + "_" + psudoRand , ext)
    }
}

func isPhotoExtention(fileName string)(ext string){
    fileName = strings.ToLower(fileName)
    if strings.HasSuffix(fileName, ".rw2"){
        return ".rw2"
    } else if strings.HasSuffix(fileName, ".jpg"){
        return ".jpg"
    } else { // TODO: Add more formats that have exif info
        return ""
    }
}

func timeTaken(photoPath string)(mmdd string, yyyy string, hhmm string){
    file, err := os.Open(photoPath)
    if err != nil {panic(err)}
    exifData, xerr := exif.Decode(file)
    if xerr != nil {panic(xerr)} // TODO: maybe just give unsupport file msg when no exif is found?
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

func copyFile(src string, dest string)(error){
    inputFile, err := os.Open(src)
    if err != nil {return fmt.Errorf("Couldn't open source file: %s", err)}
    outputFile, err := os.Create(dest)
    if err != nil {
        inputFile.Close()
        return fmt.Errorf("Couldn't open dest file: %s", err)
    }
    defer outputFile.Close() // do this when function exits
    _, err = io.Copy(outputFile, inputFile)
    inputFile.Close()
    if err != nil {return fmt.Errorf("Writing to output file failed: %s", err)}
    return nil
}

func moveFile(oldPath string, newPath string){
    if err := os.Rename(oldPath, newPath); err != nil{panic(err)}
}
