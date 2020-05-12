package local

import (
		"sync"
)

type FileSystemManager struct {
		container sync.Map
}

func (this *FileSystemManager) Add(name string, system FileSystem) *FileSystemManager {
		if system == nil {
				return this
		}
		if _, ok := this.container.Load(name); ok {
				return this
		}
		this.store(name, system)
		return this
}

func (this *FileSystemManager) Get(name string) FileSystem {
		value, ok := this.container.Load(name)
		if !ok {
				return nil
		}
		if sys, ok := value.(FileSystem); ok {
				return sys
		}
		return nil
}

func (this *FileSystemManager) Exists(name string) bool {
		if _, ok := this.container.Load(name); ok {
				return true
		}
		return false
}

func (this *FileSystemManager) store(key string, system FileSystem) bool {
		if system == nil {
				return false
		}
		this.container.Store(key, system)
		return true
}

var manager = new(FileSystemManager)

func GetFileSystemManager() *FileSystemManager {
		return manager
}
