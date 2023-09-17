import { DateTime } from 'luxon'
import Hash from '@ioc:Adonis/Core/Hash'
import {
  column,
  beforeSave,
  BaseModel, hasMany, HasMany, beforeCreate,
} from '@ioc:Adonis/Lucid/Orm'
import UserOAuthAccount from './UserOAuthAccount';
import ApiToken from './ApiToken';
import { v4 as uuid } from 'uuid';


export default class User extends BaseModel {

  public static selfAssignPrimaryKey = true


  @column({ isPrimary: true })
  public id: string

  @column({})
  public email: string

  @column()
  public pendingEmail: string | null

  @column({ serializeAs: null })
  public password: string | null

  @column()
  public githubLogin: string | null

  @column()
  public  gitlabLogin: string | null

  @column()
  public firstName: string | null

  @column()
  public lastName: string  | null

  @column({columnName: "avatar_url"})
  public avatarURL: string | null


  @column({columnName: "is_email_verified"})
  public isEmailVerified: boolean  | null


  @column({columnName: "is_2fa_enabled"})
  public is2FAEnabled: boolean  | null

  @column({columnName: "two_factor_secret"})
  public twoFactorSecret: string | null

  @hasMany(() => ApiToken)
  public tokens: HasMany<typeof ApiToken>


  @column()
  public accountActivated: boolean  | null


  @column()
  public accountBlocked: boolean  | null


  @column()
  public githubConnected: boolean  | null

  @column()
  public gitlabConnected: boolean  | null

  @column()
  public googleConnected: boolean  | null


  @hasMany(() => UserOAuthAccount, {foreignKey : 'user_id'})
  public oauthAccounts: HasMany<typeof UserOAuthAccount>;


  @column.dateTime({ autoCreate: true, serializeAs: null })
  public createdAt: DateTime

  @column.dateTime({ autoCreate: true, serializeAs: null, autoUpdate: true })
  public updatedAt: DateTime

  @beforeCreate()
  public static generateUUID(user: User) {
    user.id = uuid()
  }

  @beforeSave()
  public static async hashPassword(user: User) {
    if (user.$dirty.password) {
      user.password = await Hash.make(user.password!)
    }
  }

   public async verifyPassword(plainPassword: string) {
    if (!this.email) {
      throw new Error('Cannot verify password for non-existing user')
    }

    return await Hash.verify(this.password!, plainPassword)
   }


 // Method to return a sanitized user object
 public async getSafeUser() {
  const { password, is_email_verified, is_2fa_enabled,
    email_verification_token, reset_token, reset_token_expiry, 
    two_factor_secret, account_activated, account_blocked,
    pending_email,
    ...safeUser } = this.toJSON()
  return safeUser
}

}