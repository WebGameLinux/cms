package local

const DefaultDisk = "local"

func init() {
		GetFileSystemManager().Add(DefaultDisk, GetFileLocalSystem().setName(DefaultDisk))
}
