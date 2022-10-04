package utils

//var (
//	fs *FileSystem
//)
//
//type FileSystem interface {
//	GetFS()
//}
//
//type FileSystem struct {
//	AFS  *afero.Afero
//	IOFS *afero.IOFS
//}
//
//func (f *FileSystem) GetFS(fileSystem *afero.Fs) *FileSystem {
//	if fs == nil {
//		fs = &FileSystem{
//			AFS:  nil,
//			IOFS: nil,
//		}
//	}
//
//	return fs
//}

type UtilsI interface {
	SetAndWriteConfig(string, string) error
}

type Utils struct{}
