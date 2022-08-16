DROP FUNCTION IF EXISTS UUID_TO_BIN;

DELIMITER //

CREATE FUNCTION UUID_TO_BIN(uuid CHAR(36))
RETURNS BINARY(16) DETERMINISTIC
BEGIN
  RETURN UNHEX(CONCAT(REPLACE(uuid, '-', '')));
END; //

DELIMITER ;
