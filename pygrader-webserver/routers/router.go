// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"crypto/x509"
	"pygrader-webserver/controllers"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

// Is certificate valid for the intended purpose?
func isValidCert(cert *x509.Certificate, keyUsages ...x509.ExtKeyUsage) bool {
	for _, usage := range cert.ExtKeyUsage {
		for _, keyUsage := range keyUsages {
			if usage == keyUsage {
				return true
			}
		}
	}
	logs.Warning("No key Usage: %v", cert)
	return false
}

func SetServerHeader(ctx *context.Context) {
	ctx.Output.Header("Server", "PyGrader v0.1")

	// Delete the X-Client-Cert-CN unless it is added by a trusted a proxy
	// This so that it cannot be inserted by the client.
	ctx.Request.Header.Del("X-Client-Cert-CN")

	if ctx.Request.TLS != nil && len(ctx.Request.TLS.PeerCertificates) > 0 {
		// The first cert of the chain should be a client cert
		clientCert := ctx.Request.TLS.PeerCertificates[0]

		// As a precaution, verify that this is a client cert
		if !isValidCert(clientCert, x509.ExtKeyUsageClientAuth) {
			ctx.Abort(401, "Invalid client certificate")
			return
		}

		//clientDN := clientCert.Subject.String()
		clientCN := clientCert.Subject.CommonName

		ctx.Request.Header.Set("X-Client-Cert-CN", clientCN)

		logs.Debug("Client CN=%v", clientCN)
	}
}

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/group",
			beego.NSInclude(
				&controllers.GroupController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/module",
			beego.NSInclude(
				&controllers.ModuleController{},
			),
		),
		beego.NSNamespace("/question",
			beego.NSInclude(
				&controllers.QuestionController{},
			),
		),
	)
	beego.InsertFilter("/*", beego.BeforeExec, SetServerHeader)
	beego.AddNamespace(ns)
}
