package sql

import "testing"

func Test_Create(t *testing.T) {

	stmts, err := ParseSQL(`CREATE TABLE ` + "`user`" + ` (
		` + "`user_id`" + ` INT NOT NULL,
		` + "`special_role`" + ` VARCHAR DEFAULT NULL,
		` + "`usr_biz_type`" + ` VARCHAR DEFAULT NULL,
		` + "`user_code`" + ` VARCHAR DEFAULT NULL,
		` + "`nickname`" + ` VARCHAR DEFAULT NULL,
		` + "`avatar`" + ` VARCHAR DEFAULT NULL,
		` + "`sex`" + ` INT DEFAULT NULL,
		` + "`division_code`" + ` VARCHAR DEFAULT NULL,
		` + "`detailed_address`" + ` VARCHAR DEFAULT NULL ,
		` + "`is_enabled`" + ` INT NOT NULL DEFAULT '1',
		PRIMARY KEY (` + "`user_id`" + `),
		INDEX user_code_index (` + "`user_code`" + `)
	  );`)

	if err != nil {
		t.Fatalf("%+v", err)
	}
	t.Logf("%+v", stmts)
}
