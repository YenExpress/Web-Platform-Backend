import Route from '@ioc:Adonis/Core/Route'

export default function AuthRoute() {
    Route.group(() => {
    Route.post('/signup', 'AuthController.signupWithEmail');
    Route.post('/login', 'AuthController.loginWithEmail');
    Route.get('/:provider/redirect', 'AuthController.redirect');
    Route.get('/confirm-oauth-access/:provider/:id', 'AuthController.checkOAuthAccess').middleware('auth');
    Route.post('/:provider/callback', 'AuthController.callback');
    Route.patch('/update-profile/:id', 'AuthController.updateProfile').middleware('auth');
    Route.put('/verify-email/:token', 'AuthController.verifyEmail');
    Route.put('/confirm-email-change/:token', 'AuthController.verifyEmailUpdate');
    Route.patch('/change-password/:id', 'AuthController.changePassword').middleware('auth');
    Route.put('/resend-verification-email/:email', 'AuthController.resendVerificationEmail');
    Route.put('/send-password-reset-email/:email', 'AuthController.sendPasswordResetMail');
    Route.put('/reset-password/:token', 'AuthController.resetPassword');
    Route.delete('/logout', 'AuthController.logout').middleware('auth');
    })
}
