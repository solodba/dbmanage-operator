package v1beta1

// Origin结构体构造函数
func NewOrigin() *Origin {
	return &Origin{
		Host:     "127.0.0.1",
		Port:     3306,
		Username: "test",
		Password: "test",
	}
}

// Destination结构体构造函数
func NewDestination() *Destination {
	return &Destination{
		Endpoint:     "127.0.0.1",
		AccessKey:    "test",
		AccessSecret: "test",
		BucketName:   "test",
	}
}

// DbManageSpec结构体构造函数
func NewDbManageSpec() *DbManageSpec {
	return &DbManageSpec{
		Origin:      NewOrigin(),
		Destination: NewDestination(),
	}
}

// DbManageStatus结构体构造函数
func NewDbManageStatus() *DbManageStatus {
	return &DbManageStatus{}
}

// DbManage结构体构造函数
func NewDbManage() *DbManage {
	return &DbManage{
		Spec:   *NewDbManageSpec(),
		Status: *NewDbManageStatus(),
	}
}

// DbManageList结构体构造函数
func NewDbManageList() *DbManageList {
	return &DbManageList{
		Items: make([]DbManage, 0),
	}
}

// DbManageList结构体添加方法
func (d *DbManageList) AddItems(items ...DbManage) {
	d.Items = append(d.Items, items...)
}
