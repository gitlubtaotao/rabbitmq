package datamodels

type Product struct {
	ID           int64  `json:"id" sql:"ID" imooc:"ID" `
	ProductName  string `json:"ProductName" sql:"productName" imooc:"ProductName" validate:"required"`
	ProductNum   int64  `json:"ProductNum" sql:"productNum" imooc:"ProductNum" validate:"required"`
	ProductImage string `json:"ProductImage" sql:"productImage" imooc:"ProductImage" validate:"required"`
	ProductUrl   string `json:"ProductUrl" sql:"productUrl" imooc:"ProductUrl" validate:"required"`
}

