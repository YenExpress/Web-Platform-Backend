import BaseSchema from '@ioc:Adonis/Lucid/Schema'

export default class extends BaseSchema {
  protected tableName = 'user_oauth_accounts'

  public async up () {
    this.schema.createTable(this.tableName, (table) => {
      table.uuid('id').primary()
      table.uuid('user_id').references('id').inTable('users').onDelete('CASCADE') 
      // table.foreign('user_id').references('id').inTable('users').onDelete('CASCADE') 
      table.string("oauth_user_id").notNullable()
      table.string("oauth_provider_name").notNullable()
      table.string("access_token")
      table.string("refresh_token")

      table.timestamp('created_at', { useTz: true })
      table.timestamp('updated_at', { useTz: true })
    })
  }

  public async down () {
    this.schema.dropTable(this.tableName)
  }
}
