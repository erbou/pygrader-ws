package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {
	beego.GlobalControllerRouter["pygrader-webserver/controllers:FileController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:FileController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "GetGroups",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(
				param.New("name"),
				param.New("page"),
				param.New("pageSize"),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "GetGroup",
			Router:           `/:gid:int`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(
				param.New("gid", param.IsRequired, param.InPath),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "PutGroup",
			Router:           `/:gid:int`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(
				param.New("gid", param.IsRequired, param.InPath),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "DeleteGroup",
			Router:           `/:gid:int`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(
				param.New("gid", param.IsRequired, param.InPath),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "AddSubgroup",
			Router:           `/:gid:int/group/:sgid:int`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(
				param.New("gid", param.IsRequired, param.InPath),
				param.New("sgid", param.IsRequired, param.InPath),
				param.New("secret"),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "RemoveSubgroup",
			Router:           `/:gid:int/group/:sgid:int`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(
				param.New("gid", param.IsRequired, param.InPath),
				param.New("sgid", param.IsRequired, param.InPath),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "GetSubgroup",
			Router:           `/:gid:int/sub/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(
				param.New("gid", param.IsRequired, param.InPath),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "GetSupgroup",
			Router:           `/:gid:int/sup/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(
				param.New("gid", param.IsRequired, param.InPath),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "GetUsers",
			Router:           `/:gid:int/user/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(
				param.New("gid", param.IsRequired, param.InPath),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "AddUser",
			Router:           `/:gid:int/user/:uid:int`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(
				param.New("gid", param.IsRequired, param.InPath),
				param.New("uid", param.IsRequired, param.InPath),
				param.New("secret"),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "RemoveUser",
			Router:           `/:gid:int/user/:uid:int`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(
				param.New("gid", param.IsRequired, param.InPath),
				param.New("uid", param.IsRequired, param.InPath),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:ModuleController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:ModuleController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:ModuleController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:ModuleController"],
		beego.ControllerComments{
			Method:           "GetModules",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(
				param.New("name"),
				param.New("page"),
				param.New("pageSize"),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:ModuleController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:ModuleController"],
		beego.ControllerComments{
			Method:           "GetModule",
			Router:           `/:mid:int`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(
				param.New("mid", param.IsRequired, param.InPath),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:ModuleController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:ModuleController"],
		beego.ControllerComments{
			Method:           "PutModule",
			Router:           `/:mid:int`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(
				param.New("mid", param.IsRequired, param.InPath),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:ModuleController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:ModuleController"],
		beego.ControllerComments{
			Method:           "DeleteModule",
			Router:           `/:mid:int`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(
				param.New("mid", param.IsRequired, param.InPath),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:ModuleController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:ModuleController"],
		beego.ControllerComments{
			Method:           "GetGroupAnswers",
			Router:           `/:mid:int/answer/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(
				param.New("mid", param.IsRequired, param.InPath),
				param.New("question"),
				param.New("group"),
				param.New("poster"),
				param.New("page"),
				param.New("pageSize"),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:ModuleController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:ModuleController"],
		beego.ControllerComments{
			Method:           "PostQuestion",
			Router:           `/:mid:int/question`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(
				param.New("mid", param.IsRequired, param.InPath),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:ModuleController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:ModuleController"],
		beego.ControllerComments{
			Method:           "GetAllQuestions",
			Router:           `/:mid:int/question/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(
				param.New("mid", param.IsRequired, param.InPath),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:QuestionController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:QuestionController"],
		beego.ControllerComments{
			Method:           "GetQuestion",
			Router:           `/:qid:int`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(
				param.New("qid", param.IsRequired, param.InPath),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:QuestionController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:QuestionController"],
		beego.ControllerComments{
			Method:           "PutQuestion",
			Router:           `/:qid:int`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(
				param.New("qid", param.IsRequired, param.InPath),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:QuestionController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:QuestionController"],
		beego.ControllerComments{
			Method:           "DeleteQuestion",
			Router:           `/:qid:int`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(
				param.New("qid", param.IsRequired, param.InPath),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:QuestionController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:QuestionController"],
		beego.ControllerComments{
			Method:           "PostAnswer",
			Router:           `/:qid:int/answer`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(
				param.New("qid", param.IsRequired, param.InPath),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:UserController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:UserController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:UserController"],
		beego.ControllerComments{
			Method:           "GetUsers",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(
				param.New("email"),
				param.New("username"),
				param.New("kid"),
				param.New("page"),
				param.New("pageSize"),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:UserController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:UserController"],
		beego.ControllerComments{
			Method:           "GetUser",
			Router:           `/:uid:int`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(
				param.New("uid", param.IsRequired, param.InPath),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:UserController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:UserController"],
		beego.ControllerComments{
			Method:           "PutUser",
			Router:           `/:uid:int`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(
				param.New("uid", param.IsRequired, param.InPath),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:UserController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:UserController"],
		beego.ControllerComments{
			Method:           "DeleteUser",
			Router:           `/:uid:int`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(
				param.New("uid", param.IsRequired, param.InPath),
			),
			Filters: nil,
			Params:  nil,
		})

	beego.GlobalControllerRouter["pygrader-webserver/controllers:UserController"] = append(beego.GlobalControllerRouter["pygrader-webserver/controllers:UserController"],
		beego.ControllerComments{
			Method:           "GetGroups",
			Router:           `/:uid:int/group/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(
				param.New("uid", param.IsRequired, param.InPath),
			),
			Filters: nil,
			Params:  nil,
		})
}
