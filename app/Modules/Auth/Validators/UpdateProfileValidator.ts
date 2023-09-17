import { schema, rules, CustomMessages } from '@ioc:Adonis/Core/Validator'
import { HttpContextContract } from '@ioc:Adonis/Core/HttpContext';
import CustomError from 'App/Exceptions/CustomError';

export default class UpdateProfileValidator {

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
    firstName: schema.string.optional({ trim: true }),
    lastName: schema.string.optional({ trim: true }),
    email: schema.string.optional({ trim: true }, [
      rules.email()
    ]),
  })

  private static messages: CustomMessages = {
    'email.email': 'Invalid email format',
  }
}
