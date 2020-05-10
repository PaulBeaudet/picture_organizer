package main

import (
    "testing"
    "os" // Assumptions: Integrity of os.Stat & os.Getwd
)

const TEST_DIR = "/TestFiles/"
const TEST_JPGS = TEST_DIR + "OriginSamples/"

func TestOsAbstractions(t *testing.T){
    // TEST for -> copyFile / mkdir / rm / moveFile # Happy Path
    workingDir, _ := os.Getwd()
    src := workingDir + TEST_JPGS
    dest := workingDir + TEST_DIR
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
    cleanFolderMess(workingDir + "/TestFiles/ignore_parent")
}

func TestTimeTakenIfPhoto(t *testing.T){ // timeTakenIfPhoto -> Happy Path
    workingDir, _ := os.Getwd()
    taken, isPhoto := timeTakenIfPhoto(workingDir + TEST_JPGS + "test0.jpg")
    if isPhoto == false {
        t.Errorf("Returned as happy sample as not a photo with exif")
    }
    fmtString := "Mon Jan 2 15:04:05 -0700 MST 2006"
    if taken.Format(fmtString) != "Mon May 4 14:36:05 -0400 EDT 2020" {
        t.Errorf("Returned incorrect timestamp for happy sample")
    }
}

func TestGetValidName(t *testing.T){
    workingDir, _ := os.Getwd()
    jpgSrc := workingDir + TEST_JPGS
    newName := getValidName(jpgSrc, "test0", jpgSrc + "test0.jpg") // insert rand number for existing file case
    if newName == "test0.jpg" {
        t.Errorf("failed to create a unique name for this file");
    }
    newName = getValidName(jpgSrc, "newName", jpgSrc + "test0.jpg") // retains orginal name case
    if newName != "newName.jpg" {
        t.Errorf("failed to keep original name in absence of a conflict: " + newName)
    }
}

func TestScanAndMove(t *testing.T){
    workingDir, _ := os.Getwd()
    testDir := workingDir + TEST_DIR
    dummyDir := testDir + "ignore_dir/"
    createSrcCopy(workingDir + TEST_JPGS, dummyDir)
    scanAndMove(dummyDir, testDir, false)
    expectedParent := testDir + "2020/"
    expectedFile := expectedParent + "04_04_/16_15_35.jpg"
    // This could be improved to check the three test files to make sure files are being iterated through
    _, err := os.Stat(expectedFile)
    if os.IsNotExist(err){ // This also accounts for the file being renamed
        t.Errorf("Failed to copy first sample photo")
    }
    _, err = os.Stat(dummyDir + "2020")
    if os.IsExist(err){
        t.Errorf("No in source directory copy expeceted in false safemode")
    }
    cleanFolderMess(expectedParent) // Housekeeping: For next test w/out safemode
    // TEST 2 Does safe mode work as expected
    createSrcCopy(workingDir + TEST_JPGS, dummyDir)
    scanAndMove(dummyDir, testDir, true)
    _, err = os.Stat(dummyDir + "2020")
    if os.IsNotExist(err){
        t.Errorf("Source directory copy expeceted in true safemode")
    }
    expectedFile = dummyDir + "/2020/04_04_/16_15_35.jpg"
    _, err = os.Stat(expectedFile)
    if os.IsNotExist(err){
        t.Errorf("Failed to move working directory copy")
    }
    // End of test Housekeeping
    cleanFolderMess(expectedParent)
    cleanFolderMess(dummyDir)
}
// TODO create a test with copy image that has same timestamp

func TestMain(t *testing.T){

}

// ----- Test helpler functions ------
func createSrcCopy(src string, dest string){
    mkdir(dest)
    dir, err := os.Open(src)
    if err != nil{panic(err)}
    files, err := dir.Readdir(-1)
    dir.Close()
    if err != nil {panic(err)}
    for _, file := range files {
        fileName := file.Name()
        copyFile(src + fileName, dest + fileName)
    }
}

func cleanFolderMess(messDir string){
    err := os.RemoveAll(messDir)
    if err != nil {panic("failed to clean up folder mess")}
}
