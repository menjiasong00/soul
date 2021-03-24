package service

import (
	"golang.org/x/net/context"
	"rest/pkg/gmysql"
	"rest/model"
	"rest/pb"
	"rest/pkg/tools"
)

type {{.ServerName}}Server struct{}

//{{.ModuleName}}List {{.Name}}列表
func (s *{{.ServerName}}Server) {{.ModuleName}}List(ctx context.Context, in *pb.{{.ModuleName}}ListRequest) (*pb.{{.ModuleName}}ListResponse, error) {

	resp := &pb.{{.ModuleName}}List{}

	dbModel := gmysql.DB.Model(model.{{.ModelName}}{})

/*	if len(in.Name) > 0 {
		dbModel=dbModel.Where("  name like ?","%"+in.Name+"%")
	}*/

	dbModel.Count(&resp.Total)

	if in.PageSize > 0 {
		//分页
		offset := (in.PageNumber - 1) * in.PageSize
		dbModel = dbModel.Offset(offset).Limit(in.PageSize)
	}
	
	sortStr := " id DESC"
	if in.OrderKey != "" && in.OrderSort != "" {
		sortStr = in.OrderKey + " " + in.OrderSort
	}
	dbModel = dbModel.Order(sortStr)

	dbModel.Scan(&resp.List)

/*	for _, v := range resp.List {

	}*/

	return &pb.{{.ModuleName}}ListResponse{Status: 200, Message: "success", Data: resp}, nil
}

//{{.ModuleName}}Detail {{.Name}}详情
func (s *{{.ServerName}}Server) {{.ModuleName}}Detail(ctx context.Context, in *pb.{{.ModuleName}}IdRequest) (*pb.{{.ModuleName}}DetailResponse, error) {

	resp := &pb.{{.ModuleName}}OneRequest{}
	gmysql.DB.Model(model.{{.ModelName}}{}).Where("id = ?",in.Id).Scan(&resp)

	return &pb.{{.ModuleName}}DetailResponse{Status: 200, Message: "success", Data: resp}, nil
}


//{{.ModuleName}}Create {{.Name}}新建
func (s *{{.ServerName}}Server) {{.ModuleName}}Create(ctx context.Context, in *pb.{{.ModuleName}}OneRequest) (*pb.{{.ModuleName}}Response, error) {
    
	//表单验证
    errValidate := in.Validate()
    	if errValidate != nil {
    		return nil, errValidate
    }

	{{.ModuleName}}One := model.{{.ModelName}}{}
	tools.ScanStuct(in,&{{.ModuleName}}One)

	gmysql.DB.Create(&{{.ModuleName}}One)

	return &pb.{{.ModuleName}}Response{Status: 200, Message: "success", Data:true}, nil
}

//{{.ModuleName}}Motify {{.Name}}修改
func (s *{{.ServerName}}Server) {{.ModuleName}}Motify(ctx context.Context, in *pb.{{.ModuleName}}OneRequest) (*pb.{{.ModuleName}}Response, error) {
    
	//表单验证
    errValidate := in.Validate()
    	if errValidate != nil {
    		return nil, errValidate
    }

	{{.ModuleName}}One := model.{{.ModelName}}{}
	gmysql.DB.Model(model.{{.ModelName}}{}).Where("id = ?",in.Id).Find(&{{.ModuleName}}One)
	tools.ScanStuct(in,&{{.ModuleName}}One)

	gmysql.DB.Model(model.{{.ModelName}}{}).Where(" id = ?",in.Id).Save(&{{.ModuleName}}One)

	return &pb.{{.ModuleName}}Response{Status: 200, Message: "success", Data:true}, nil
}

//{{.ModuleName}}Delete {{.Name}}删除
func (s *{{.ServerName}}Server) {{.ModuleName}}Delete(ctx context.Context, in *pb.{{.ModuleName}}IdRequest) (*pb.{{.ModuleName}}Response, error) {

    {{.ModuleName}}One := model.{{.ModelName}}{}
    gmysql.DB.Model(model.{{.ModelName}}{}).First(&{{.ModuleName}}One,in.Id)
    //{{.ModuleName}}One.Status = 2
    gmysql.DB.Save(&{{.ModuleName}}One)
	
	return &pb.{{.ModuleName}}Response{Status: 200, Message: "success", Data:true}, nil
}






