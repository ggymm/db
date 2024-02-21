INSERT INTO device ('device_id', 'device_name') VALUES ('1~sd\n==dfds', '2'), ('3', '4');
UPDATE device SET device_id   = 1, device_name = 'pname \t\<>12' WHERE device_id = 1;
DELETE FROM device WHERE device_id = 1;
SELECT * FROM device WHERE device_id = 1 OR device_id = 2 AND device_name = 'pname \t\\<>12 ' LIMIT 10, 10;
