import User from 'App/Models/User'
import UserOAuthAccount from 'App/Models/UserOAuthAccount';
import ApiToken from 'App/Models/ApiToken';
import CustomError from 'App/Exceptions/CustomError';
import { DateTime } from 'luxon'
import { createUniqueToken, calculateExpiryDateFromString } from 'App/Utils/Tokenizer';
import { ApiTokenType } from '../../../../types';
import { v4 as uuid } from 'uuid';
// import axios, { AxiosRequestConfig}  from 'axios';
// import axiosRetry from 'axios-retry';
import Env from '@ioc:Adonis/Core/Env';
// import Ally from '@ioc:Adonis/Addons/Ally';

// Configure Axios to use axios-retry
// axiosRetry(axios, {
//   retries: 3, // Number of retries
//   retryDelay: axiosRetry.exponentialDelay, // Exponential back-off
// });


export default class AuthService{

  private baseEMSURL : string =  Env.get("EMS_URL")
  private baseClientURL : string = Env.get("CLIENT_BASE_URL")

  constructor() {}

  public async createOrFetchUserwithOAuth(oauthUserId: string, email: string, firstName: string, 
    lastName: string, isVerified: boolean, oauthProviderName: string, 
    accessToken: string, refreshToken?: string, oauthLogin?: string): Promise<User> {

      const existingOAuthAccount = await UserOAuthAccount.query()
      .where('oauth_user_id',  oauthUserId)
      .where('oauth_provider_name', oauthProviderName)
      .preload('user')
      .first()

      if (existingOAuthAccount){
        existingOAuthAccount.accessToken = accessToken;
        existingOAuthAccount.refreshToken = refreshToken ?? null;

        return existingOAuthAccount.user
      }

    var user = await User.firstOrCreate({email}, {
        email, firstName, lastName,
        isEmailVerified: isVerified, accountActivated: isVerified,
        githubLogin: oauthProviderName === "github" ? (oauthLogin !== undefined ? oauthLogin : null) : null,
        gitlabLogin: oauthProviderName === "gitlab" ? (oauthLogin !== undefined ? oauthLogin : null) : null,
        githubConnected: oauthProviderName === "github" ? true : null,
        gitlabConnected: oauthProviderName === "gitlab" ? true : null,
      }) 
    
    if (!user.githubConnected &&  oauthProviderName === "github") {
      user.githubConnected = true;
      user.githubLogin = oauthLogin ?? null
      user = await user.save()
    } else if (!user.gitlabConnected &&  oauthProviderName === "gitlab"){
      user.gitlabConnected = true;
      user.gitlabLogin = oauthLogin ?? null
      user = await user.save()
    }

    await UserOAuthAccount.firstOrCreate({oauthProviderName, oauthUserId},{
        userId: user.id, oauthProviderName, oauthUserId, accessToken, refreshToken})


      if (!user.isEmailVerified){
        try{
          await this.resendVerificationEmail(email)
        }catch(e){
        }
      }

    return user
  }
  


  public async signupWithEmail(email: string, password: string, firstName: string, lastName: string): Promise<Object> {
    try { 

      const newUser = await User.create({
        email, password, firstName, lastName,
        isEmailVerified: false, accountActivated: false
      })

      const email_verification_token = await this.generateAPIToken(newUser.id, "2 d", ApiTokenType.mailVerification)
      
      await this.sendMail({
        "confirmationLink":`${this.baseClientURL}/confirm-email/${email_verification_token}`,
        "recipient": email})

    return await newUser.getSafeUser()
    } catch(error){
      if (error.code == '23505') {
            throw new CustomError("ERR_DUPLICATE_EMAIL", "User With Email Already Exists")
        }
      throw error
      }
  }
  

  public async loginWithEmail(email: string, password: string):
   Promise<{ user: Object, 
    accessToken: string, refreshToken: string,
    tokenType: "Bearer"}> {
      const existingUser = await User.findBy("email", email)
      if (!existingUser) {
          throw new CustomError("ERR_INVALID_EMAIL","Email Does Not Exist")
        } else if (! await existingUser.verifyPassword(password)) {
          throw new CustomError("ERR_INVALID_PASSWORD","Password Incorrect")
      }
      return await this.getSessionData(existingUser!)
    }

  public async updateProfile(userId: string, firstName?: string, lastName?: string, email?: string) {
      const user = await User.find(userId)
      if (!user){
        throw new CustomError("ERR_INVALID_EMAIL","Profile Does Not Exist")
      }
      const filteredData = Object.entries({
        firstName,
        lastName,
      })
    .filter(([_, value]) => value !== undefined)
    .reduce((obj, [key, value]) => ({ ...obj, [key]: value }), {});
      user.merge(filteredData)
      if (email && email!= user.email){
        const existingUserWithEmail = await User.findBy("email", email)
        if (existingUserWithEmail){
            throw new CustomError("ERR_DUPLICATE_EMAIL","New Email Provided belongs to an active user")
        }
        const email_verification_token = await this.generateAPIToken(user.id, "2 d", ApiTokenType.mailVerification)
        user.pendingEmail = email
        await user.save()
        await this.sendMail({
          "confirmationLink":`${this.baseClientURL}/confirm-email/${email_verification_token}`,
          "recipient": email})
        const profile = await user.getSafeUser()
        profile.pending_email = email
        return profile
      }
      await user.save()
      return user.getSafeUser()
    }


  public async generateAPIToken(userId: string, expireIn:string, 
   type: ApiTokenType, sessionId?: string): Promise<string> {
    const value = await createUniqueToken(32)
    const expiry = calculateExpiryDateFromString(expireIn)
    await ApiToken.create({
      userId,
      value,
      type,
      sessionId,
      expiry,
    })

    return value
  }

  public async validateAPIToken(token: string, type: ApiTokenType): Promise<ApiToken | null>{
    const foundToken = await ApiToken.query()
      .where('value', token)
      .where('type', type)
      .where('expiry', '>', DateTime.now().toISO()!)
      .preload('user')
      .first()
    if (foundToken){
        return foundToken
      }
    return null
  }


  public async getSessionData(user: User):
   Promise<{ user: Object, 
    accessToken: string, refreshToken: string,
    tokenType: "Bearer"}> {
    
      if (!user.accountActivated || !user.isEmailVerified || user.accountBlocked) {
          throw new CustomError("ERR_USER_INACTIVE","User Account Inactive Either Because Email Not Verified Or Account Suspended")
      }
      const session_id = uuid()
      return {
            user: await user.getSafeUser(),
            accessToken : await this.generateAPIToken(user.id, '7 d', ApiTokenType.access ,session_id),
            refreshToken: await this.generateAPIToken(user.id, '30 d', ApiTokenType.refresh, session_id),
            tokenType: "Bearer"
          }
    }


  public async changePassword(userId: string, oldPassword:string, newPassword: string) {
      const user = await User.find(userId);
      if (!user){
        throw new CustomError("ERR_INVALID_ID","User Does Not Exist")
      }
      else if (! await user.verifyPassword(oldPassword)) {
        throw new CustomError("ERR_INVALID_PASSWORD","Old Password Provided Incorrect")
      }
      user.password = newPassword;
      await user.save();
    } 


  public async verifyEmail(token: string) {

      const foundToken = await this.validateAPIToken(token, ApiTokenType.mailVerification)
      if (!foundToken) {
        throw new CustomError("ERR_INVALID_TOKEN","Email Verification Token Incorrect or Expired")
        }
      const user = foundToken.user

      user.isEmailVerified = true
      user.accountActivated = true

      await user.save()

      await foundToken.delete()

    } 

  public async verifyEmailChange(token: string) {

        const foundToken = await this.validateAPIToken(token, ApiTokenType.mailVerification)
        if (!foundToken) {
          throw new CustomError("ERR_INVALID_TOKEN", "Email Verification Token Incorrect or Expired")
          }

        const user = foundToken.user
  
        user.email = user.pendingEmail!
        user.pendingEmail = null
  
        await user.save()
        await foundToken.delete()
  
    }

  public async resendVerificationEmail(email: string) {

    const user = await User.findBy('email', email)
    if (!user) {
      throw new CustomError("ERR_INVALID_EMAIL","Email Does Not Exist")
    }
    const newVerificationToken = await createUniqueToken(32);
    const expiry = calculateExpiryDateFromString("2 d");
    const payload = await ApiToken.firstOrCreate({userId: user.id, 
      type: ApiTokenType.mailVerification}, {
        userId: user.id, type: ApiTokenType.mailVerification,
        value: newVerificationToken, expiry
    });
    payload.value = newVerificationToken;
    payload.expiry = expiry;

    await payload.save()

    const sent = await this.sendMail({
      "confirmationLink":`${this.baseClientURL}/confirm-email/${newVerificationToken}`,
      "recipient": email})

      if (!sent){
        throw new CustomError("ERR_MAIL_SERVICE","Failed To Send Verification Email")
    }
} 


  public async sendPasswordResetMail(email: string) {
    const user = await User.findBy('email', email)
    if (!user) {
      throw new CustomError("ERR_INVALID_EMAIL","User Does Not Exist")
    }

    const newResetToken = await createUniqueToken(32);
    const expiry = calculateExpiryDateFromString("20 m");
    const payload = await ApiToken.firstOrCreate({userId: user.id, 
      type: ApiTokenType.mailVerification}, {
        userId: user.id, type: ApiTokenType.passwordReset,
        value: newResetToken, expiry
    });
    payload.value = newResetToken;
    payload.expiry = expiry;

    await payload.save()

    const sent = await this.sendMail({
      "confirmationLink":`${this.baseClientURL}/confirm-password-reset/${newResetToken}`,
      "recipient": email})
    
    if (!sent){
        throw new CustomError("ERR_MAIL_SERVICE","Failed To Send Password Reset Email")
    }
}

  public async resetPassword(token: string, newPassword: string) {

  const foundToken = await this.validateAPIToken(token, ApiTokenType.passwordReset)
  if (!foundToken) {
    throw new CustomError("ERR_INVALID_TOKEN","Email Verification Token Incorrect or Expired")
  }

  const user = foundToken.user;
  user.password = newPassword;
  await user.save()
  await foundToken.delete()

}


public async deleteSessionTokens(accessToken: string) {
  const foundAccessToken = await ApiToken.findBy("value", accessToken)
  const foundRefreshToken =  await ApiToken.query()
  .where('session_id', foundAccessToken!.sessionId!)
  .where('type', ApiTokenType.refresh)
  .first()
  await foundAccessToken!.delete()
  await foundRefreshToken!.delete()
}


public async sendMail(body: object, msURL: string = this.baseEMSURL) {
  console.log(body)
  console.log(msURL)
//   const config: AxiosRequestConfig = {
//     method: 'post',
//     url: msURL,
//     data: body,
//   };

//   try {
//     await axios(config);
//     return true;
//   } catch (e) {
//     return false
// }
}
}
