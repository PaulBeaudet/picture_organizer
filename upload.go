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
    const wd = "/home/wooper/Dropbox/programs/golang/upload_organizer/"
    const testFile = "testFile1.txt"
    year, _, _ := time.Now().Date()
    photosFolderDefualt = photosFolderDefualt + strconv.Itoa(year) + "/"

    fileDest := uploadDefaultFolder + testFile
    writeFile(fileDest, "this is a test of the erm");
    modTime := getCreation(fileDest)
    mkdir(photosFolderDefualt, modTime);
    moveFile(fileDest, photosFolderDefualt + modTime + "/" + testFile);
    // readFile(testFile);

}

func getCreation(pathToFile string)(creationTime string){
    fileStat, err := os.Stat(pathToFile)
    if err != nil {
        panic(err)
    }
    timeObj := fileStat.ModTime()
    modTime := timeObj.Format("01_02_2006")
    return modTime
}

func mkdir(path string, nameOfDir string){
    dirToCreate := path + nameOfDir
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
