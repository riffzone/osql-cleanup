package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	if len(os.Args) <= 1 {
		fmt.Println("Usage : osql-cleanup.exe <folder> [<days>]")
		fmt.Println("<folder> is the path of the folder to scan, sub-folders are ignored")
		fmt.Println("The .BAK files older than <days> days are deleted, default value is 7 days")
		fmt.Println("")
		fmt.Println("The .BAK files can be created by the following command :")
		fmt.Println("\"[PATH_TO]OSQL.exe\" -S <SQL_SERVER> -U <USER> -P <PASS> -Q \"BACKUP DATABASE <BASE_NAME> TO DISK='[EXPORT_PATH]%DATE:~6,4%%DATE:~3,2%%DATE:~0,2%-<BASE_NAME>.BAK'\"")
		return
	}

	myArgScanDirPath := os.Args[1]
	myArgScanDirPath = strings.TrimRight(myArgScanDirPath, "/")

	myArgDaysAgo := 7
	if len(os.Args) > 2 {
		myArgDaysAgoStr := os.Args[2]
		if myArgDaysAgoStr != "" {
			myStrConvInt, myStrConvErr := strconv.Atoi(myArgDaysAgoStr)
			if myStrConvErr != nil {
				fmt.Println(myArgDaysAgoStr + " is not an integer : ", myStrConvErr)
				return
			}
			myArgDaysAgo = myStrConvInt
		}
	}

	fmt.Println("Scanning directory : ", myArgScanDirPath)
	fmt.Println("Keep files for : ", myArgDaysAgo, " days")

	myScanFiles, myScanErr := os.ReadDir(myArgScanDirPath)
	if myScanErr != nil {
		fmt.Println("Error reading directory : ", myScanErr)
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
			fmt.Println("Stat error : ", myStatErr)
			continue
		}
						
		myScanFileModTime := myScanFileInfo.ModTime()
		myDaysAgo := time.Now().AddDate(0, 0, -myArgDaysAgo)

		if !myScanFileModTime.Before(myDaysAgo) {
			continue
		}

		myOSErr := os.Remove(myScanFilePath)
		if myOSErr != nil {
			fmt.Printf("Error deleting file "+myScanFilePath+" : %s\n", myOSErr)
			return
		}

		fmt.Println("Deleted file : ", myScanFilePath)
	}

}
