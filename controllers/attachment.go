package controllers

import (
		log "github.com/sirupsen/logrus"
		"mime/multipart"
		"path"
)

type ApiAttachmentController interface {
		Upload()
		Uploads()
		Download()
		GetById()
		Remove()
		Lists()
}

type AttachmentController struct {
		BaseController
}

var attachmentController *AttachmentController

func GetAttachmentController() ApiAttachmentController {
		if attachmentController == nil {
				attachmentController = new(AttachmentController)
		}
		return attachmentController
}

func (this *AttachmentController) Upload() {
		file, information, err := this.GetFile("_file")
		if err != nil {
				this.ApiResponse(nil)
				return
		}
		defer this.close(file)
		//将文件信息头的信息赋值给filename变量
		filename := information.Filename
		//保存文件的路径。保存在static/upload中   （文件名）
		err = this.SaveToFile("file", path.Join("static/upload", filename))
		this.ApiResponse(nil)
}

func (this *AttachmentController) close(file multipart.File) {
		err := file.Close()
		log.Error(err)
}

func (this *AttachmentController) Uploads() {
		files,err:=this.GetFiles("_files")
		if err != nil {
				this.ApiResponse(nil)
				return
		}
		//将文件信息头的信息赋值给filename变量
		for _,file:=range files {
				//保存文件的路径。保存在static/upload中   （文件名）
				err = this.SaveToFile("file", path.Join("static/upload", file.Filename))
		}

		this.ApiResponse(nil)
}

func (this *AttachmentController) Download() {
		panic("implement me")
}

func (this *AttachmentController) GetById() {
		panic("implement me")
}

func (this *AttachmentController) Remove() {
		panic("implement me")
}

func (this *AttachmentController) Lists() {
		panic("implement me")
}
