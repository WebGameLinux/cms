package response

import "github.com/WebGameLinux/cms/models"

type LoginRespJson struct {
		User *models.User `json:"user"`
		Auth string       `json:"authorization"`
}

type LoginFilterRespJson struct {
		User *UserFilterRespJson `json:"user"`
		Auth string              `json:"authorization"`
}

type UserFilterRespJson struct {
		Id       int64  `json:"id"`
		Email    string ` json:"email"`
		UserName string `json:"username"`
		Mobile   string `json:"mobile"`
		Gender   int    `json:"gender"`
		models.UniqueSeqKey
		models.CreateDate
}

func (this *LoginFilterRespJson) Init(login *LoginRespJson) *LoginFilterRespJson {
		this.User = new(UserFilterRespJson).Init(login.User)
		this.Auth = login.Auth
		return this
}

func (this *UserFilterRespJson) Init(user *models.User) *UserFilterRespJson {
		this.Id = user.Id
		this.Email = user.Email
		this.UserName = user.UserName
		this.Mobile = user.Mobile
		this.Gender = user.Gender
		this.SeqId = user.SeqId
		this.CreatedAt = user.CreatedAt
		this.UpdatedAt = user.UpdatedAt
		return this
}
