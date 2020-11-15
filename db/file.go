package db

import (
	"database/sql"
	mydb "fileServices/db/mysql"
	"fmt"
)

//文件上传完成 保存meta
func OnFileUploadFinished(fileHash string, fileName string,
	fileSize int64, fileAddr string) int32 {
	stmt, err := mydb.DBConn().Prepare(
		"insert into tbl_file(`file_sha1`, `file_name`, `file_size`," +
			"`file_addr`, `status`) values(?,?,?,?,1)")
	if err != nil {
		fmt.Printf("failed to prepare statement, err:%s", err.Error())
		return 0
	}
	defer stmt.Close()

	ret, err := stmt.Exec(fileHash, fileName, fileSize, fileAddr)
	if err != nil {
		fmt.Printf(err.Error())
		return 0
	}
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf(err.Error())
		return 0
	}

	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("File with hash:%s has been uploaded before", fileHash)
			return 0
		}
	}
	return int32(id)
}

type TableFile struct {
	Id       int32
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

//mysql 获取元信息
func GetFileMeta(id int32) (*TableFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select id, file_sha1, file_name, file_size, file_addr from tbl_file" +
			" where id = ? and status = 1 limit 1")
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	defer stmt.Close()

	tfile := TableFile{}
	err = stmt.QueryRow(id).Scan(&tfile.Id, &tfile.FileHash, &tfile.FileName,
		&tfile.FileSize, &tfile.FileAddr)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	return &tfile, nil

}
