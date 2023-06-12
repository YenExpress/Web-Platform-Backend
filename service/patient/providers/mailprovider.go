package providers

import (
	"YenExpress/service/postoffice"
)

var (
	signUpMailHandler = postoffice.MailHandler[postoffice.OneTimePassword]{
		Mailer: postoffice.PostMan[postoffice.OneTimePassword]{
			MailTemplatePath: "/service/postoffice/templates/newAccountEmailValidation.html",
			Subject:          "New Account Email Validation",
		},
	}
)
