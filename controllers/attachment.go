package controllers

import (
		"github.com/WebGameLinux/cms/dto/response"
		"github.com/astaxie/beego"
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
		GetAccessToken()
}

type AttachmentController struct {
		BaseController
}

var attachmentController *AttachmentController

func GetAttachmentController() beego.ControllerInterface {
		if attachmentController == nil {
				attachmentController = new(AttachmentController)
		}
		return attachmentController
}

func (this *AttachmentController) URLMapping() {
		this.Mapping("Lists", this.Lists)
		this.Mapping("Upload", this.Upload)
		this.Mapping("Uploads", this.Uploads)
		this.Mapping("GetById", this.GetById)
}

func (this *AttachmentController) Upload() {
		file, information, err := this.GetFile("_file")
		res := response.RespJson{
				Data: nil,
				Msg:  "",
				Code: 0,
		}
		if err != nil {
				res.Msg = err.Error()
				res.Code = -2
				this.Data["json"] = res
				this.ServeJSON()
				return
		}
		defer this.close(file)
		//将文件信息头的信息赋值给filename变量
		filename := information.Filename
		//保存文件的路径。保存在static/upload中   （文件名）
		err = this.SaveToFile("_file", path.Join("static/uploads", filename))
		if err != nil {
				res.Msg = err.Error()
				res.Code = -1
		}
		this.Data["json"] = res
		this.ServeJSON()
}

func (this *AttachmentController) close(file multipart.File) {
		err := file.Close()
		log.Error(err)
}

func (this *AttachmentController) Uploads() {
		files, err := this.GetFiles("_files")
		if err != nil {
				this.ApiResponse(nil)
				return
		}
		//将文件信息头的信息赋值给filename变量
		for _, file := range files {
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

func (this *AttachmentController) GetAccessToken() {
		panic("implement me")
}
