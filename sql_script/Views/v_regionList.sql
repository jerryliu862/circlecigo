CREATE OR REPLACE VIEW v_regionList
AS
    SELECT R.code, R.name, C.code AS currencyCode, C.name AS currencyName, C.format AS currencyFormat
    FROM Region R
    LEFT JOIN Currency C ON R.currencyCode = C.code
    WHERE R.show IS TRUE
