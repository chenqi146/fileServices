package handler

import (
	"fileServices/meta"
	"fileServices/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		//接收文件流存储在本地目录
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data, err:%s\n", err.Error())
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: header.Filename,
			Location: "./file/" + header.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Faild to create file, err:%s\n", err.Error())
			return
		}
		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Faild to sava data info file, err:%s\n", err.Error())
			return
		}

		_, _ = newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		id := meta.UpdateFileMetaDB(fileMeta)
		if id == 0 {
			err := os.Remove(fileMeta.Location)
			if err != nil {
				fmt.Printf("删除文件失败: %s, location: %s", fileMeta.FileName, fileMeta.Location)
			}
			_, _ = w.Write(util.NewRespMsg(500, "error", nil).JSONBytes())
			return
		}
		fileMeta.Id = id
		meta.UpdateFileMeta(fileMeta)
		resp := util.NewRespMsg(200, "success", fileMeta)
		_, _ = w.Write(resp.JSONBytes())
	}

}

//GetFileMetaHandler  获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	idStr := r.Form["id"][0]
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fMeta, err := meta.GetFileMetaDB(int32(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp := util.NewRespMsg(200, "success", fMeta)
	_, _ = w.Write(resp.JSONBytes())
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	idStr := r.Form.Get("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fm := meta.GetFileMeta(int32(id))

	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	//读取文件到内存
	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("content-disposition", "attachment;filename=\""+fm.FileName+"\"")
	_, _ = w.Write(data)
}

func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	idStr := r.Form.Get("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fMeta := meta.GetFileMeta(int32(id))
	_ = os.Remove(fMeta.Location)

	meta.RemoveFileMete(int32(id))
	w.WriteHeader(http.StatusOK)
	resp := util.NewRespMsg(200, "success", nil)
	_, _ = w.Write(resp.JSONBytes())
}
