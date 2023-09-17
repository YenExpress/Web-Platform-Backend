import { schema, rules, CustomMessages } from '@ioc:Adonis/Core/Validator'
import { HttpContextContract } from '@ioc:Adonis/Core/HttpContext';
import CustomError from 'App/Exceptions/CustomError';

export default class ChangePasswordValidator {
  
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
    oldPassword: schema.string({ trim: true }),
    newPassword: schema.string([ rules.minLength(8),
      rules.regex(/[!@#$%^&*(),.?":{}|<>]/), // Check for at least one special character
      rules.regex(/\d/),  // Check for at least one digit
     ]
    ),
  })

  private static messages: CustomMessages = {
    'oldPassword.required': 'The old password field is required',
    'newPassword.required': 'The new password field is required',
    'newPassword.minLength': 'The new password must be at least 8 characters long',
    'newPassword.regex': 'The new password must contain at least one special character and one digit',

  }
}


