// import { AuthenticationException } from '@adonisjs/auth/build/standalone'
// import type { GuardsList } from '@ioc:Adonis/Addons/Auth'
import type { HttpContextContract } from '@ioc:Adonis/Core/HttpContext'
import CustomError from '../../../Exceptions/CustomError'
import AuthService from '../Services/AuthService'
import { ApiTokenType } from '../../../../types';


/**
 * Auth middleware is meant to restrict un-authenticated access to a given route
 * or a group of routes.
 *
 * You must register this middleware inside `start/kernel.ts` file under the list
 * of named middleware.
 */
export default class AuthMiddleware {

  private service: AuthService
  
  constructor() {
    this.service = new AuthService()
  }

  protected async authenticate(auth: HttpContextContract['auth'], request: HttpContextContract['request']) {

    const token = request.header('Authorization');
      if (!token) {
        throw new CustomError("ERR_INVALID_AUTH","Access token missing from Request Header",)
      }
      const validatedToken= await this.service.validateAPIToken(token.replace(/^Bearer\s+/i, ''), ApiTokenType.access)
      if (!validatedToken){
        throw new CustomError("ERR_INVALID_AUTH","Bearer Token in Request Header Invalid or Expired")
      }
      auth.use('api').user = validatedToken.user

  }

  public async handle (
    { auth, response, request }: HttpContextContract,
    next: () => Promise<void>,
  ) {
    try{
      await this.authenticate(auth, request)
      await next()
    }catch(e){
      return response.status(401).json({
        success: false, 
        message: "Authentication Failed",
        error: e.message
      })
    }
  }
}
