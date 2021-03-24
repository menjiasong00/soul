package service

import (
	"golang.org/x/net/context"
	"rest/config/db"
	"rest/model"
	"rest/pb"
	"rest/util/tools"
	"rest/util/wlog"
	"rest/util/bas"
)


type {{.ServerName}}Server struct{}

//{{.ModuleName}}List {{.Name}}列表
func (s *{{.ServerName}}Server) {{.ModuleName}}List(ctx context.Context, in *pb.{{.ModuleName}}ListRequest) (*pb.{{.ModuleName}}ListResponse, error) {

	resp := &pb.{{.ModuleName}}List{}

	dbModel := db.WithContext(ctx,db.{{.DbName}}DB).Model(model.{{.ModelName}}{})

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
	db.WithContext(ctx,db.{{.DbName}}DB).Model(model.{{.ModelName}}{}).Where("id = ?",in.Id).Scan(&resp)

	return &pb.{{.ModuleName}}DetailResponse{Status: 200, Message: "success", Data: resp}, nil
}


//{{.ModuleName}}Create {{.Name}}新建
func (s *{{.ServerName}}Server) {{.ModuleName}}Create(ctx context.Context, in *pb.{{.ModuleName}}OneRequest) (*pb.{{.ModuleName}}Response, error) {
    //表单验证
    errValidate := in.Validate()
    	if errValidate != nil {
    		return nil, errValidate
    }
    authId, err := tools.GetAuthID(ctx)
	if err != nil {
		return nil, err
	}
	{{.ModuleName}}One := model.{{.ModelName}}{}
	tools.ScanStuct(in,&{{.ModuleName}}One)

	err= db.WithContext(ctx,db.{{.DbName}}DB).Create(&{{.ModuleName}}One).Error
	wlog.CheckErr(err, "")



    bas.Auditlog(&pb.Auditlog{
		Module:"{{.ServerName}}",
		Object:"{{.ModelName}}",
		ObjectId:int32({{.ModuleName}}One.Id),
		Operation:"added",
		OperatorBadge:authId,
	},&in)

	return &pb.{{.ModuleName}}Response{Status: 200, Message: "success", Data:true}, nil
}

//{{.ModuleName}}Motify {{.Name}}修改
func (s *{{.ServerName}}Server) {{.ModuleName}}Motify(ctx context.Context, in *pb.{{.ModuleName}}OneRequest) (*pb.{{.ModuleName}}Response, error) {
    //表单验证
    errValidate := in.Validate()
    	if errValidate != nil {
    		return nil, errValidate
    }
    authId, err := tools.GetAuthID(ctx)
	if err != nil {
		return nil, err
	}
	{{.ModuleName}}One := model.{{.ModelName}}{}
	db.WithContext(ctx,db.{{.DbName}}DB).Model(model.{{.ModelName}}{}).Where("id = ?",in.Id).Find(&{{.ModuleName}}One)
	tools.ScanStuct(in,&{{.ModuleName}}One)

	err= db.WithContext(ctx,db.{{.DbName}}DB).Model(model.{{.ModelName}}{}).Where(" id = ?",in.Id).Save(&{{.ModuleName}}One).Error
	wlog.CheckErr(err, "")


    bas.Auditlog(&pb.Auditlog{
		Module:"{{.ServerName}}",
		Object:"{{.ModelName}}",
		ObjectId:int32(in.Id),
		Operation:"edited",
		OperatorBadge:authId,
	},&{{.ModuleName}}One)

	return &pb.{{.ModuleName}}Response{Status: 200, Message: "success", Data:true}, nil
}

//{{.ModuleName}}Delete {{.Name}}删除
func (s *{{.ServerName}}Server) {{.ModuleName}}Delete(ctx context.Context, in *pb.{{.ModuleName}}IdRequest) (*pb.{{.ModuleName}}Response, error) {
	authId, err := tools.GetAuthID(ctx)
	if err != nil {
		return nil, err
	}

    {{.ModuleName}}One := model.{{.ModelName}}{}
    db.WithContext(ctx,db.{{.DbName}}DB).Model(model.{{.ModelName}}{}).First(&{{.ModuleName}}One,in.Id)
    {{.ModuleName}}One.Status = 2
    err = db.{{.DbName}}DB.Save(&{{.ModuleName}}One).Error
	wlog.CheckErr(err, "")

    bas.Auditlog(&pb.Auditlog{
        Module:"{{.ServerName}}",
        Object:"{{.ModelName}}",
        ObjectId:int32(in.Id),
        Operation:"deleted",
        OperatorBadge:authId,
    },&{{.ModuleName}}One)

	return &pb.{{.ModuleName}}Response{Status: 200, Message: "success", Data:true}, nil
}






