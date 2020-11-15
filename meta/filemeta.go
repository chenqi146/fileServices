package meta

import (
	mydb "fileServices/db"
)

//FileMeta 文件元信息结构
type FileMeta struct {
	Id       int32  `json:"id"`
	FileSha1 string `json:"fileSha1"`
	FileName string `json:"fileName"`
	FileSize int64  `json:"fileSize"`
	Location string `json:"location"`
	UploadAt string `json:"uploadAt"`
}

var fileMetas map[int32]FileMeta

func init() {
	fileMetas = make(map[int32]FileMeta)
}

//新增/更新文件元信息
func UpdateFileMeta(meta FileMeta) {
	fileMetas[meta.Id] = meta
}

//新增/更新文件元信息 mysql
func UpdateFileMetaDB(meta FileMeta) int32 {
	return mydb.OnFileUploadFinished(meta.FileSha1, meta.FileName,
		meta.FileSize, meta.Location)
}

//获取文件的元信息
func GetFileMeta(id int32) FileMeta {
	return fileMetas[id]
}

func GetFileMetaDB(id int32) (FileMeta, error) {
	tfile, err := mydb.GetFileMeta(id)
	if err != nil {
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		Id:       tfile.Id,
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return fmeta, nil

}

//删除
func RemoveFileMete(id int32) {
	delete(fileMetas, id)
}
