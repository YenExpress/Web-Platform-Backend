import BaseSchema from '@ioc:Adonis/Lucid/Schema'

export default class extends BaseSchema {
  protected tableName = 'users'

  public async up () {
    this.schema.createTable(this.tableName, (table) => {
      table.uuid('id').primary()
      table.string('email', 100).notNullable().unique()
      table.string('pending_email', 100).unique()
      table.string('last_name')
      table.string('first_name')
      table.string('github_login')
      table.string('gitlab_login')
      table.string('avatar_url')
      table.boolean('is_email_verified')
      table.boolean('account_activated')
      table.boolean('account_blocked')
      table.boolean('github_connected')
      table.boolean('google_connected')
      table.boolean('gitlab_connected')
      table.string('password')

      /**
       * Uses timestamptz for PostgreSQL and DATETIME2 for MSSQL
       */
      table.timestamp('created_at', { useTz: true })
      table.timestamp('updated_at', { useTz: true })

    })
  }

  public async down () {
    this.schema.dropTable(this.tableName)
  }
}
