package main

type db interface {
	open(url string)
	close()
	printFileList()
	addFile(file monitoredFile)
	removeFile(file monitoredFile)
}
