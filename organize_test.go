package main

import (
    "testing"
    "os"
)

func TestOsAbstractions(t *testing.T){ // Assumptions: Integrity of os.Stat & os.Getwd
    // TEST for -> copyFile / mkdir / rm / moveFile
    workingDir, _ := os.Getwd()
    src := workingDir + "/TestFiles/OriginSamples/"
    dest := workingDir + "/TestFiles/"
    copiedFile := dest + "ignore_copied.jpg"
    copyFile(src + "test0.jpg", copiedFile) // TEST copyFile
    _, err := os.Stat(copiedFile)
    if os.IsNotExist(err){
        t.Errorf("Failed to copy or rename file")
    }
    dirToCreate := workingDir + "/TestFiles/ignore_parent/ignore_child"
    mkdir(dirToCreate) // TEST mkdir
    _, err = os.Stat(dirToCreate)
    if os.IsNotExist(err){
        t.Errorf("failed to create parent and child directory")
    }
    movedFile := dirToCreate + "/ignore_rename.jpg"
    moveFile(copiedFile, movedFile) // TEST moveFile
    _, err = os.Stat(movedFile)
    if os.IsNotExist(err){
        t.Errorf("Failed to move or rename file")
    }
    rm(movedFile) // TEST rm
    _, err = os.Stat(movedFile)
    if os.IsExist(err){
        t.Errorf("Failed to remove file");
    }
    // --- Housekeeping ----
    err = os.RemoveAll(workingDir + "/TestFiles/ignore_parent")
    if err != nil {t.Errorf("failed to clean up folder mess")}
}

func TestTimeTakenIfPhoto(t *testing.T){

}

func TestGetValidName(t *testing.T){

}

func TestScanAndMove(t *testing.T){

}

func TestMain(t *testing.T){

}
