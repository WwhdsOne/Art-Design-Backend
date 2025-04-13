package entity

// Permission 定义权限模型
type Permission struct {
	ID       int64   `gorm:"primaryKey;type:bigint;comment:'权限ID'"`
	Name     string  `gorm:"type:varchar(20);not null;unique;comment:'页面名称'"`
	ParentId int64   `gorm:"type:bigint;not null;comment:'父页面ID,根节点父页面ID为-1'"`
	Path     string  `gorm:"type:varchar(256);comment:'访问路径'"`
	ChildIds []int64 `gorm:"-:all;comment:'子页面ID列表(不存储)'"`
}

func (p *Permission) TableName() string {
	return "permission"
}
