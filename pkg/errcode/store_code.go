package errcode

var (
	FileNotExist = NewError(1001, "File %s not Exist")
	DirNotExist  = NewError(1002, "Directory is not Exist Path %s")
)
