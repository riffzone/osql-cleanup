# osql-cleanup

OSQL.exe is a Microsoft SQL Server utility to export the content of a database in .BAK files.  
osql-cleanup.exe is designed to be automatically launched in a daily basis, and deletes all the OSQL .BAK files from an export folder that are older than a specified number of days.  
  
Usage : osql-cleanup.exe --dir=\<folder\> [--days=\<days\>]  
\<folder\> is the path of the folder to scan, sub-folders are ignored  
The .BAK files older than \<days\> days are deleted, default value is 7 days  
  
The .BAK files can be created by the following command :  
````
"[PATH_TO]OSQL.exe" -S <SQL_SERVER> -U <USER> -P <PASS> -Q "BACKUP DATABASE <BASE_NAME> TO DISK='[EXPORT_PATH]%DATE:~6,4%%DATE:~3,2%%DATE:~0,2%-<BASE_NAME>.BAK'"
````
