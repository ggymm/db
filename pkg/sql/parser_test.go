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
	for i, stmt := range stmts {
		t.Logf("%d %+v", i, stmt)
	}
}

func TestSelect(t *testing.T) {
	//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
	stmts, err := ParseSQL(`SELECT * FROM device WHERE device_id = 1 OR device_id = 2 AND device_name = 'pname \t\\<>12 ' LIMIT 10, 10;`)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	for i, stmt := range stmts {
		t.Logf("%d %+v", i, stmt)
	}
}

func TestInsert(t *testing.T) {
	//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
	stmts, err := ParseSQL("INSERT INTO device ('device_id' , 'device_name' ) VALUES ('1~sd\n==dfds','2'),('3','4');")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	for i, stmt := range stmts {
		t.Logf("%d %+v", i, stmt)
	}
}

func TestUpdate(t *testing.T) {
	//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
	stmts, err := ParseSQL("UPDATE device SET device_id = 1, device_name = 'pname \t\\<>12 ' WHERE device_id = 1;")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	for i, stmt := range stmts {
		t.Logf("%d %+v", i, stmt)
	}
}

func TestDelete(t *testing.T) {
	//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
	stmts, err := ParseSQL("DELETE FROM device WHERE device_id = 1;")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	for i, stmt := range stmts {
		t.Logf("%d %+v", i, stmt)
	}
}
