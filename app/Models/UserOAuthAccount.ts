import { BaseModel, column, belongsTo, BelongsTo, beforeCreate } from '@ioc:Adonis/Lucid/Orm';
import User from './User';
import { v4 as uuid } from 'uuid';
import { DateTime } from 'luxon'

export default class UserOAuthAccount extends BaseModel {

  public static selfAssignPrimaryKey = true;
  public static table = 'user_oauth_accounts';


  @column({ isPrimary: true })
  public id: string;

  @column()
  public userId: string;

  @column({columnName: "oauth_provider_name"})
  public oauthProviderName: string;

  @column({columnName: "oauth_user_id"})
  public oauthUserId: string;

  @column()
  public accessToken: string | null ;

  @column()
  public refreshToken: string | null ;

  @column.dateTime({ autoCreate: true, serializeAs: null })
  public createdAt: DateTime

  @column.dateTime({ autoCreate: true, serializeAs: null, autoUpdate: true })
  public updatedAt: DateTime

  @belongsTo(() => User)
  public user: BelongsTo<typeof User>;

  @beforeCreate()
  public static generateUUID(account: UserOAuthAccount) {
    account.id = uuid()
  }

}
