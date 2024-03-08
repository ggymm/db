SELECT *
FROM device
WHERE device_id = 1
   OR device_id = 2 AND device_name = 'pname \t\\<>12 ' LIMIT 10, 10;