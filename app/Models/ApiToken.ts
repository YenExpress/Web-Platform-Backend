import { DateTime } from 'luxon'
import { BaseModel, column, BelongsTo, belongsTo, beforeCreate} from '@ioc:Adonis/Lucid/Orm'
import User from './User'
import { v4 as uuid } from 'uuid';

export default class ApiToken extends BaseModel {

  public static selfAssignPrimaryKey = true
  
  @column({ isPrimary: true })
  public id: string

  @column()
  public userId: string;

  @column()
  public sessionId: string | null;

  @column()
  public value: string

  @column()
  public type: string

  @column()
  public expiry: string
  
  @column.dateTime({ autoCreate: true, serializeAs: null })
  public createdAt: DateTime

  @column.dateTime({ autoCreate: true, serializeAs: null, autoUpdate: true })
  public updatedAt: DateTime

  // Define a relationship with the User model
  @belongsTo(() => User)
  public user: BelongsTo<typeof User>

  @beforeCreate()
  public static generateUUID(token: ApiToken) {
    token.id = uuid()
  }
}
