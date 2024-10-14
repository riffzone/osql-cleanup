package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {

	myOSNewLine := "\n"
	if runtime.GOOS == "windows" {
		myOSNewLine = "\r\n"
	}

	if len(os.Args) <= 1 {
		fmt.Print("Usage : osql-cleanup.exe <folder> [<days>]", myOSNewLine)
		fmt.Print("<folder> is the path of the folder to scan, sub-folders are ignored", myOSNewLine)
		fmt.Print("The .BAK files older than <days> days are deleted, default value is 7 days", myOSNewLine)
		fmt.Print("", myOSNewLine)
		fmt.Print("The .BAK files can be created by the following command :", myOSNewLine)
		fmt.Print("\"[PATH_TO]OSQL.exe\" -S <SQL_SERVER> -U <USER> -P <PASS> -Q \"BACKUP DATABASE <BASE_NAME> TO DISK='[EXPORT_PATH]%DATE:~6,4%%DATE:~3,2%%DATE:~0,2%-<BASE_NAME>.BAK'\"", myOSNewLine)
		return
	}

	myArgScanDirPath := os.Args[1]
	myArgScanDirPath = strings.ReplaceAll(myArgScanDirPath, "\\\\", "\\")
	myArgScanDirPath = strings.ReplaceAll(myArgScanDirPath, "\\", "/")
	myArgScanDirPath = strings.TrimRight(myArgScanDirPath, "/")

	myArgDaysAgo := 7
	if len(os.Args) > 2 {
		myArgDaysAgoStr := os.Args[2]
		if myArgDaysAgoStr != "" {
			myStrConvInt, myStrConvErr := strconv.Atoi(myArgDaysAgoStr)
			if myStrConvErr != nil {
				fmt.Print(myArgDaysAgoStr + " is not an integer : ", myStrConvErr, myOSNewLine)
				return
			}
			myArgDaysAgo = myStrConvInt
		}
	}

	fmt.Print("Scanning directory : ", myArgScanDirPath, myOSNewLine)
	fmt.Print("Keep files for : ", myArgDaysAgo, " days", myOSNewLine)

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
		myDaysAgo := time.Now().AddDate(0, 0, -myArgDaysAgo)

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
