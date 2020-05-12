package local

func init()  {
		GetFileSystemManager().Add("local",GetFileLocalSystem())
}
