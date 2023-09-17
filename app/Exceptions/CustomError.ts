import { Exception } from '@adonisjs/core/build/standalone';

export default class CustomError extends Exception {

  constructor(code: string, message: string) {
    super(message=message)
    this.code = code

    // Capture the current stack trace
      Error.captureStackTrace(this, CustomError);
  }
}
