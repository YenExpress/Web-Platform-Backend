package providers

import (
	"YenExpress/service/postoffice"
)

var (
	loginMailHandler = postoffice.MailHandler[postoffice.OneTimePassword]{
		Mailer: postoffice.PostMan[postoffice.OneTimePassword]{
			MailTemplatePath: "/postoffice/templates/concurrentLoginValidation.html",
			Subject:          "Login Validation",
		},
	}

	signUpMailHandler = postoffice.MailHandler[postoffice.OneTimePassword]{
		Mailer: postoffice.PostMan[postoffice.OneTimePassword]{
			MailTemplatePath: "/postoffice/templates/newAccountEmailValidation.html",
			Subject:          "New Account Email Validation",
		},
	}
)
