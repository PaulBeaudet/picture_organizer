package main

import (
    "testing"
    "os" // Assumptions: Integrity of os.Stat & os.Getwd
    "time"
    "strconv"
)

const TEST_DIR = "/TestFiles/"
const TEST_JPGS = TEST_DIR + "OriginSamples/"

func TestScanAndMove(t *testing.T){ // Options on sorting behaviours expected
    workingDir, _ := os.Getwd()
    testDir := workingDir + TEST_DIR
    dummyDir := testDir + "ignore_dir/"
    createSrcCopy(workingDir + TEST_JPGS, dummyDir)
    daysSinceYoungestPic := getYoungestPic(t, dummyDir)
    eventName := "event"
    scanAndMove(dummyDir, testDir, false, eventName, daysSinceYoungestPic)
    expectedParent := testDir + "2020/"
    expectedFile := expectedParent + "05_04_" + eventName + "/14_36_05.jpg" // This is youngest photo
    // This could be improved to check the three test files to make sure files are being iterated through
    _, err := os.Stat(expectedFile)
    if os.IsNotExist(err){ // Option: Chosen destination created & sorted to
        t.Errorf("Failed to copy first sample photo")
    }
    _, err = os.Stat(dummyDir + "2020")
    if os.IsExist(err){    // Option: False safemode
        t.Errorf("No in source directory copy expeceted in false safemode")
    }
    numberOfDirs := amountOfFilesInDir(t, expectedParent)
    if numberOfDirs != 1 { // Option: Asked for youngest files
        t.Errorf("Should only be one event directory; asked for youngest sample files. Got:" + strconv.Itoa(numberOfDirs))
    }
    // End of test: Housekeeping below
    cleanFolderMess(expectedParent)
    cleanFolderMess(dummyDir)
}
// TODO create a test with copy image that has same timestamp

func TestMain(t *testing.T){    // expected default behaviour without passed flags
    workingDir, _ := os.Getwd() // get the working directory
    home, _ := os.UserHomeDir() // get $HOME // Defaults will at least work in linux
    dummyDir := workingDir + TEST_DIR + "ignore_dir/"
    expectedParent := home + "/Pictures/2020/"
    expectedFile := expectedParent + "04_04_/16_15_35.jpg"
    createSrcCopy(workingDir + TEST_JPGS, dummyDir)
    os.Chdir(dummyDir)          // point app's working director at dummy one
    main()                      // TODO: How would one test flag options?
    _, err := os.Stat(expectedFile)
    if os.IsNotExist(err){ // Default: Sort photos to ~/Pictures
        t.Errorf("Failed to copy first sample photo")
    }
    expectedFile = dummyDir + "/2020/04_04_/16_15_35.jpg"
    _, err = os.Stat(expectedFile)
    if os.IsNotExist(err){ // Default: Safemode = In working dirctory copy of photos
        t.Errorf("Failed to move working directory copy")
    }
    // End of test: Housekeeping below
    cleanFolderMess(expectedParent)
    cleanFolderMess(dummyDir)
    os.Chdir(workingDir)
}

// ----- Test helpler functions ------
func amountOfFilesInDir(t *testing.T, srcDir string)(int){
    t.Helper()
    _, err := os.Stat(srcDir)
    if os.IsNotExist(err){return 0}
    dirContents, err := os.Open(srcDir)
    if err != nil{panic(err)}
    files, err := dirContents.Readdir(-1)
    dirContents.Close()
    if err != nil {panic(err)}
    return len(files)
}

func getYoungestPic(t *testing.T, srcDir string)(int){
    t.Helper()
    dirContents, err := os.Open(srcDir)
    if err != nil{panic(err)}
    files, err := dirContents.Readdir(-1)
    dirContents.Close()
    if err != nil {panic(err)}
    youngest := 0
    firstPhotoIndex := 0;
    now := time.Now()
    for index, file := range files {
        fileName := file.Name()
        taken, isPhoto := timeTakenIfPhoto(srcDir + fileName)
        if !isPhoto{
            if index==0 {firstPhotoIndex = firstPhotoIndex + 1}
            continue
        }
        diff := now.Sub(taken)
        daysSinceTaken := int(diff.Hours()/24)
        if firstPhotoIndex == index || daysSinceTaken < youngest{
            youngest = daysSinceTaken
        }
    }
    return youngest
}

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

func cleanFolderMess(messDir string){ // pay attention not to point this at source origin working directory
    err := os.RemoveAll(messDir)
    if err != nil {panic("failed to clean up folder mess")}
}
