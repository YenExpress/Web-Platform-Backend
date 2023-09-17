import { schema, rules, CustomMessages } from '@ioc:Adonis/Core/Validator'
import { HttpContextContract } from '@ioc:Adonis/Core/HttpContext'
import CustomError from 'App/Exceptions/CustomError';

export default class CreateUserValidator {
  
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
    firstName: schema.string({ trim: true }),
    lastName: schema.string({ trim: true }),
    email: schema.string({ trim: true }, [
      rules.email(),
      // rules.unique({ table: 'users', column: 'email' }),
    ]),
    password: schema.string([ rules.minLength(8),
      rules.regex(/[!@#$%^&*(),.?":{}|<>]/), // Check for at least one special character
      rules.regex(/\d/),  // Check for at least one digit
     ]
     
    ),
  })

  private static messages: CustomMessages = {
    'firstName.required': 'The firstName field is required',
    'lastName.required': 'The lastName field is required',
    'email.required': 'The email field is required',
    'email.email': 'Invalid email format',
    'password.required': 'The password field is required',
    'password.minLength': 'The password must be at least 8 characters long',
    'password.regex': 'The password must contain at least one special character and one digit',

  }
}


