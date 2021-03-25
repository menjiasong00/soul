syntax = "proto3";
option go_package = "pb";

package proto;
import "google/api/annotations.proto";
import "validate/validate.proto";

//{{.Name}}
service {{.ServerName}}{

    //{{.Name}}列表
    rpc {{.ModuleName}}List ({{.ModuleName}}ListRequest) returns ({{.ModuleName}}ListResponse) {
        option (google.api.http) = {
            get: "/{{.Url}}"
		};
    }

    //{{.Name}}详情
    rpc {{.ModuleName}}Detail ({{.ModuleName}}IdRequest) returns ({{.ModuleName}}DetailResponse) {
        option (google.api.http) = {
            get: "/{{.Url}}/{id}"
		};
    }

    //{{.Name}}创建
    rpc {{.ModuleName}}Create({{.ModuleName}}OneRequest) returns ({{.ModuleName}}Response) {
        option (google.api.http) = {
            post: "/{{.Url}}"
            body: "*"
		};
    }

    //{{.Name}}修改
    rpc {{.ModuleName}}Motify({{.ModuleName}}OneRequest) returns ({{.ModuleName}}Response) {
        option (google.api.http) = {
            put: "/{{.Url}}"
            body: "*"
		};
    }

    //{{.Name}}删除
    rpc {{.ModuleName}}Delete({{.ModuleName}}IdRequest) returns ({{.ModuleName}}Response) {
        option (google.api.http) = {
            delete: "/{{.Url}}"
		};
    }
}

//{{.Name}}通用
message {{.ModuleName}}Response{
    int32 status = 1;
    string message = 2;
       //bool data = 3;
    int32 code = 4;
    string error = 5;
    bool details = 6;
}

//{{.Name}}列表
message {{.ModuleName}}ListRequest{
    //每页多少个
    int32 pageSize = 100;
    //页码
    int32 pageNumber = 101;
    //排序字段
    string orderKey =102;
    //排序方式
    string orderSort =103;
}

message {{.ModuleName}}ListResponse{
    int32 status = 1;
    string message = 2;
    //{{.ModuleName}}List data = 3;
    int32 code = 4;
    string error = 5;
    {{.ModuleName}}List  details = 6;
}

message {{.ModuleName}}List {
    //总数
    int32 total = 1;
    //列表
    repeated {{.ModuleName}}OneRequest list = 2;
}

//{{.Name}}结构
message {{.ModuleName}}OneRequest {
    {{range .Columns}}//{{.ColumnComment}}
    {{.TypeName}} {{.ColumnNameCamel}} ={{.Id}};
    {{end}}
}

//{{.Name}}详情
message {{.ModuleName}}IdRequest{
    int32 id =1;
}

message {{.ModuleName}}DetailResponse{
    int32 status = 1;
    string message = 2;
    //{{.ModuleName}}OneRequest data = 3;
    int32 code = 4;
    string error = 5;
    {{.ModuleName}}OneRequest  details = 6;
}
