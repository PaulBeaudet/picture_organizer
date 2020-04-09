package main

import (
    "fmt"
    "os"
    "time"
    "strconv"
)

func main(){
    const uploadDefaultFolder = "/home/wooper/Pictures/test/uploads/"
    photosFolderDefualt := "/home/wooper/Pictures/test/personal_pictures/"
    // testMove(uploadDefaultFolder, photosFolderDefualt)
    ls(uploadDefaultFolder)
    scanAndMove(uploadDefaultFolder, photosFolderDefualt)
    ls(uploadDefaultFolder)
    ls(photosFolderDefualt)
}

func scanAndMove(uploadDir string, photoDir string){
    uploads, uErr := os.Open(uploadDir)
    if uErr != nil{panic(uErr)}
    files, error := uploads.Readdir(-1)
    uploads.Close()
    if error != nil{panic(error)}
    fmt.Println("scanning contents of " + uploadDir)
    for _, file := range files {
        move(file, uploadDir, photoDir)
    }
}

func move(file os.FileInfo, sourceDir string, destDir string){
    currentLocation := sourceDir + file.Name()
    monthDay, year := getCreation(currentLocation)
    fmt.Println(file.Name() + " was created in " + year + " on " + monthDay);
    nextDest := destDir + year + "/" + monthDay
    mkdir(nextDest)
    moveFile(currentLocation, nextDest + "/" + file.Name())
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

func getCreation(pathToFile string)(creationDay string, creationYear string){
    fileStat, err := os.Stat(pathToFile)
    if err != nil {
        panic(err)
    }
    timeObj := fileStat.ModTime()
    monthDay := timeObj.Format("01_02_")
    year := timeObj.Format("2006")
    return monthDay, year
}

func mkdir(dirToCreate string){
    // dirToCreate := path + nameOfDir
    _, error := os.Stat(dirToCreate)

    if os.IsNotExist(error){
        if err := os.MkdirAll(dirToCreate, 0755); err != nil {
            panic(err);
        }
    }
}

func moveFile(oldPath string, newPath string){
    if err := os.Rename(oldPath, newPath); err != nil{
        panic(err)
    }
    fmt.Println("reorganizing " + oldPath + " to " + newPath)
}

func writeFile(fileName string, fileString string){
    file, err := os.Create(fileName)
    if err != nil {
        panic(err)
    }
    file.WriteString(fileString + "\n")
    file.Close()
}

func readFile(fileName string){
    file, err := os.Open(fileName)
    if err != nil{
        panic(err)
    }
    data := make([]byte, 30)
    file.Read(data)
    fmt.Printf(string(data) + "\n")
    file.Close()
}


func testMove(uploadDefaultFolder string, photosFolderDefualt string){
    // const wd = "/home/wooper/Dropbox/programs/golang/upload_organizer/"
    year, _, _ := time.Now().Date()
    photosFolderDefualt = photosFolderDefualt + strconv.Itoa(year) + "/"
    const testFile = "testFile1.txt"
    fileDest := uploadDefaultFolder + testFile
    writeFile(fileDest, "this is a test of the erm");
    ls(uploadDefaultFolder)
    modTime, _ := getCreation(fileDest)
    mkdir(photosFolderDefualt + modTime);
    newFolderForFile := photosFolderDefualt + modTime + "/"
    moveFile(fileDest, newFolderForFile + testFile);
    ls(newFolderForFile)
    ls(uploadDefaultFolder)
    readFile(newFolderForFile + testFile);
}
