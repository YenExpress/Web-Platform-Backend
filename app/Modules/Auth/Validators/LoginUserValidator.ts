import { schema, rules, CustomMessages } from '@ioc:Adonis/Core/Validator'
import { HttpContextContract } from '@ioc:Adonis/Core/HttpContext';
import CustomError from 'App/Exceptions/CustomError';


export default class LoginUserValidator {
  
  public static async handle(request: HttpContextContract["request"]){
    try{
      return await request.validate({
        schema: this.schema,
        messages: this.messages
      })
    } catch(e){
      throw new CustomError("ERR_INVALID_PAYLOAD", e.messages)
    }
  }


  private static schema = schema.create({
    email: schema.string({ trim: true }, [ rules.required(),
      rules.email(),
    ]),
    password: schema.string({}, [rules.required(), 
    ]),
  })

  private static messages: CustomMessages = {
    'email.required': 'The email field is required',
    'email.email': 'Invalid email format',
    'password.required': 'The password field is required',
  }
}



