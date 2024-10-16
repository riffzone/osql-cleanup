package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

const kDefaultArgDays = 7

func main() {

	myOSNewLine := "\n"
	if runtime.GOOS == "windows" {
		myOSNewLine = "\r\n"
	}

	myPtrArgDir := flag.String("dir", "", "the folder to clean")
	myPtrArgDays := flag.Int("days", 0, "the number of days to keep")
    flag.Parse()

	myArgsOK := true
	if *myPtrArgDir == "" {
		myArgsOK = false
	}

	myArgDays := *myPtrArgDays
	if myArgDays == 0 {
		myArgDays = kDefaultArgDays
	}

	if !myArgsOK {
		fmt.Print("Usage : osql-cleanup.exe --dir=<folder> [--days=<days>]", myOSNewLine)
		fmt.Print("<folder> is the path of the folder to scan, sub-folders are ignored", myOSNewLine)
		fmt.Print("The .BAK files older than <days> days are deleted, default value is ", kDefaultArgDays, " days", myOSNewLine)
		fmt.Print("", myOSNewLine)
		fmt.Print("The .BAK files can be created by the following command :", myOSNewLine)
		fmt.Print("\"[PATH_TO]OSQL.exe\" -S <SQL_SERVER> -U <USER> -P <PASS> -Q \"BACKUP DATABASE <BASE_NAME> TO DISK='[EXPORT_PATH]%DATE:~6,4%%DATE:~3,2%%DATE:~0,2%-<BASE_NAME>.BAK'\"", myOSNewLine)
		return
	}

	myArgScanDirPath := *myPtrArgDir
	myArgScanDirPath = strings.ReplaceAll(myArgScanDirPath, "\\\\", "\\")
	myArgScanDirPath = strings.ReplaceAll(myArgScanDirPath, "\\", "/")
	myArgScanDirPath = strings.TrimRight(myArgScanDirPath, "/")

	fmt.Print("Scanning directory : ", myArgScanDirPath, myOSNewLine)
	fmt.Print("Keep files for : ", myArgDays, " days", myOSNewLine)

	myScanFiles, myScanErr := os.ReadDir(myArgScanDirPath)
	if myScanErr != nil {
		fmt.Print("Error reading directory : ", myScanErr, myOSNewLine)
		return
	}

	for _, myScanFile := range myScanFiles {
		if myScanFile.IsDir() {
			continue
		}
		myScanFileName := myScanFile.Name()

		if !strings.HasSuffix(myScanFileName, ".BAK") {
			continue
		}

		myScanFilePath := myArgScanDirPath + "/" +  myScanFileName

		myScanFileInfo, myStatErr := os.Stat(myScanFilePath)
		if myStatErr != nil {
			fmt.Print("Stat error : ", myStatErr, myOSNewLine)
			continue
		}
						
		myScanFileModTime := myScanFileInfo.ModTime()
		myDaysAgo := time.Now().AddDate(0, 0, -myArgDays)

		if !myScanFileModTime.Before(myDaysAgo) {
			continue
		}

		myOSErr := os.Remove(myScanFilePath)
		if myOSErr != nil {
			fmt.Printf("Error deleting file "+myScanFilePath+" : %s"+myOSNewLine, myOSErr)
			return
		}

		fmt.Print("Deleted file : ", myScanFilePath, myOSNewLine)
	}

}
