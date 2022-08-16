CREATE OR REPLACE VIEW v_regionTaxList
AS
    SELECT R.code, R.name, C.code AS currencyCode, C.name AS currencyName, C.format AS currencyFormat, T.payType, T.taxFrom, T.taxRate
    FROM Region R
    LEFT JOIN Currency C ON R.currencyCode = C.code
    LEFT JOIN TaxRate T ON R.code = T.regionCode
    WHERE R.show IS TRUE
