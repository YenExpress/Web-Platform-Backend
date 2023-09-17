import { HttpContextContract } from '@ioc:Adonis/Core/HttpContext'
import AuthService from '../../Services/AuthService'
import CreateUserValidator from '../../Validators/CreateUserValidator';
import LoginUserValidator from '../../Validators/LoginUserValidator';
import UpdateProfileValidator from '../../Validators/UpdateProfileValidator';
import ChangePasswordValidator from '../../Validators/ChangePasswordValidator';
import ResetPasswordValidator from '../../Validators/ResetPasswordValidator';

export default class AuthController {

  private authService: AuthService

   constructor() {
    this.authService = new AuthService()
  }


  public async signupWithEmail({ request, response }: HttpContextContract) {
    try{
      const payload = await CreateUserValidator.handle(request)
      const newUser = await this.authService.signupWithEmail(payload.email, payload.password, payload.firstName, payload.lastName)
      return response.status(201).json({
        success: true,
        message: "User Account Successfully Created, Verification mail sent to email address!",
        data: newUser,
          })
    }catch(e){
      if (e.code == "ERR_INVALID_PAYLOAD") {
        return response.status(422).json({
          success: false,
          message: "Bad Request Body", 
          error: e.message,
        })
      } else if (e.code == "ERR_DUPLICATE_EMAIL") {
        return response.status(409).json({
          success: false, 
          error: e.message,
          message: "Registration Failed"
        })
      }
      return response.status(500).json({
        success: false,
        message: "Registration Failed",
        error: "Internal Service Error",
      })
    }
  }

  public async loginWithEmail({ request, response}: HttpContextContract) {
    try{
      const payload = await LoginUserValidator.handle(request)
      const authData = await this.authService.loginWithEmail(payload.email, payload.password)
      return response.status(201).json({
        success: true,
        message: "Login Successful",
        data: authData,
          })
    }catch(e){
      if (e.code == "ERR_INVALID_PAYLOAD") {
        return response.status(422).json({
          success: false, 
          error: e.message,
          message: "Bad Request Body"
        })
      } else if (e.code == "ERR_INVALID_EMAIL" || e.code == "ERR_INVALID_PASSWORD"){
        return response.status(400).json({
          success: false, 
          error: "Invalid Login Credentials",
          message: "Login Failed"
        })
      } else if (e.code == "ERR_USER_INACTIVE"){
        return response.status(424).json({
          success: false, 
          error: e.message,
          message: "Login Failed"
        })
      }
      return response.status(500).json({
        success: false, 
        message: "Login Failed",
        error: "Internal Service Error"
      })
    }
  }


  public async updateProfile({ request, auth, params, response}: HttpContextContract) {
      const user = auth.user!
      if (user.id != params.id){
          return response.status(403).json({
            success: false, 
            error: "Cannot Modify Resource Belonging to Another User",
            message: "Profile Update Failed",
          })
        }
      try{
        const payload = await UpdateProfileValidator.handle(request)
        const newProfile = await this.authService.updateProfile(user.id, 
          payload.firstName, payload.lastName, payload.email)
        return response.status(200).json({
            success: true, 
            message: "User Profile Updated",
            data: newProfile
          })
      } catch(e){
        if (e.code == "ERR_INVALID_PAYLOAD") {
          return response.status(422).json({
            success: false, 
            error: e.message,
            message: "Bad Request Body"
          })
        } else if (e.code == "ERR_INVALID_EMAIL" || e.code == "ERR_DUPLICATE_EMAIL") {
          return response.status(400).json({
            success: false, 
            error: e.message,
            message: "Profile Update Failed",
          })
        }
        return response.status(500).json({
          success: false, 
          error: "Internal Service Error",
          message: "Profile Update Failed",
        })
      }
    }

  public async changePassword({ request, auth, params, response }: HttpContextContract) {
    const user = auth.user!
    if (user.id != params.id){
        return response.status(403).json({
          success: false,
          error: "Cannot Modify Resource Belonging to Another User",
          message: "Password Change Failed",
        })
      }
    try{
      const payload = await ChangePasswordValidator.handle(request)
      await this.authService.changePassword(user.id, payload.oldPassword,payload.newPassword)
      return response.status(200).json({
        success: true, 
        message: "Password Changed Successfully"
      })
    } catch(e){
      if (e.code == "ERR_INVALID_PAYLOAD") {
        return response.status(422).json({
          success: false, 
          error: e.message,
          message: "Bad Request Body"
        })
    } else if (e.code == "ERR_INVALID_ID" || e.code == "ERR_INVALID_PASSWORD") {
      return response.status(400).json({
        success: false, 
        error: e.message,
        message:"Password Change Failed",
      })
  }
    return response.status(500).json({
      success: false, 
      error: "Internal Service Error",
      message:"Password Change Failed",
    })
  }}


  public async verifyEmail({ params, response }: HttpContextContract) {
    try {
      const token  = params.token
      await this.authService.verifyEmail(token)
      return response.status(200).json({
        success: true,
        message: "Email Verification Successful"
              })
        } catch(e){
          if (e.code == "ERR_INVALID_TOKEN"){
            return response.status(401).json({
              success: false, 
              error: e.message,
              message: "Failed to verify email"
            })
          }
          return response.status(500).json({
            success: false, 
            error: "Internal Service Error",
            message: "Email Verification Failed"
          })
    }
      }
    
    
  public async verifyEmailUpdate({ params, response }: HttpContextContract) {
      try {
        const token  = params.token
        await this.authService.verifyEmailChange(token)
        return response.status(200).json({
          success: true,
          message: "Email verification and Update Successful"
                })
      } catch(e){
        if (e.name == "ERR_INVALID_TOKEN"){
          return response.status(401).json({
            success: false, 
            error: e.message,
            message: "Email Update Failed"
          })
        }
        return response.status(500).json({
          success: false, 
          message: "Email Verification and Update Failed",
          error: "Internal Service Error"
        })
  }
}


  public async resendVerificationEmail({ params, response }: HttpContextContract) {
    try {
      const email  = params.email
      await this.authService.resendVerificationEmail(email)
    } catch(e){
      if (e.code == "ERR_INVALID_EMAIL"){
        return response.status(400).json({
          success: false, 
          error: e.message,
          message: "Failed To Send Verification Mail"
        })
      } else if (e.code == "ERR_MAIL_SERVICE"){
        return response.status(500).json({
          success: false, 
          error: e.message,
          message: "Failed To Send Verification Mail"
        })
      }
      return response.status(500).json({
        success: false, 
        error: "Internal Service Error",
        message: "Failed To Send Verification Mail"
      })
    }
  }
  


  public async logout({ request, response}: HttpContextContract) {
      const accessToken = request.header('Authorization')!.replace(/^Bearer\s+/i, '');
      await this.authService.deleteSessionTokens(accessToken)
      return response.status(200).json({
          success: true,
          message: "User Logged Out of Session Successfully",
           })
        }


  public async resetPassword({ params,request, response }: HttpContextContract) {
    try{
      const token  = params.token
      const payload = await ResetPasswordValidator.handle(request)
      await this.authService.resetPassword(token, payload.newPassword)
      return response.status(200).json({
        success: true,
        message: "User Password Reset Successfully",
         })
    } catch(e){
        if (e.code == "ERR_INVALID_TOKEN"){
          return response.status(401).json({
            success: true,
            error: e.message,
            message: "Password Reset Failed"
             })
        }
      return response.status(500).json({
        success: false,
        message: 'Password Reset Failed',
        error: "Internal service Error"
      })}}
    


  public async sendPasswordResetMail({ params, response }: HttpContextContract) {

      try{
        const email  = params.email
        await this.authService.sendPasswordResetMail(email)
        return response.status(200).json({
          success: true,
          message: "Password Reset Mail Sent",
           })
      }catch(e){
        if (e.code == "ERR_INVALID_EMAIL"){
          return response.status(400).json({
            success: false,
            error: e.message,
            message: "Failed To send Password Reset Mail"
             })
        } else if (e.code == "ERR_MAIL_SERVICE"){
          return  response.status(500).json({
            success: false,
            error: e.message,
            message: "Failed To send Password Reset Mail"
             })
        }
        return  response.status(500).json({
          success: false,
          error: "Internal Service error",
          message: "Failed To send Password Reset Mail"
           })
      }
  }
    }