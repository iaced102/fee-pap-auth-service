package incoming

/*
create by: Hoangnd
create at: 2023-01-01
des: Khai báo và xử lý thông tin nhận từ request
*/

type SystemAccountParams struct {
	Id                  string `json:"id" param:"id" db:"id"  form:"id" query:"id"  primary:"true"`
	Description         string `json:"description" param:"description" db:"description" form:"description" pg:"description"`
	UserRelated         string `json:"userRelated" param:"userRelated" db:"userRelated" query:"userRelated" form:"userRelated" pg:"user_related"`
	SystemAccountTypeId string `json:"systemAccountTypeId" param:"systemAccountTypeId" db:"systemAccountTypeId" query:"systemAccountTypeId" form:"systemAccountTypeId" pg:"system_account_type_id"`
	UserName            string `json:"userName" param:"userName" db:"userName" form:"userName" pg:"user_name"`
	CreatedAt           string `json:"createdAt" param:"createdAt"  db:"createdAt" form:"createdAt" pg:"created_at"`
	UpdatedAt           string `json:"updatedAt" param:"updatedAt" db:"updatedAt" form:"updatedAt" pg:"updated_at"`
	Password            string `json:"password" param:"password"  db:"password" form:"password" pg:"password"`
	SearchKey           string `json:"searchKey" param:"searchKey"  db:"searchKey" form:"searchKey"`
}
