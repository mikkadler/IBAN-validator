# IBAN-validator

A custom web service for validating IBANs.

## Coverage

The service validates checksums and IBAN lengths for all countries listed in the [wiki](https://en.wikipedia.org/wiki/International_Bank_Account_Number#Validating_the_IBAN), excluding aspirational countries. It also identifies bank names for Estonian, Latvian, and Lithuanian banks. The Lithuanian bank list is not complete but includes major banks.

## Usage

**Prerequisites**: Ensure you have Go installed or Docker installed.

To run the web service locally:

```
go build .
./web-iban
```

Alternatively, you can use Docker Compose for easy deployment:

```
docker compose up
```

## API Endpoint

The service accepts an array of strings as input, formatted as JSON and labeled with "data". Only the POST method is accepted.

## How it Works

The core function, `ValidateIBAN`, follows the standard steps for validating IBANs:

1. **Input Validation**: Checks for the correct length and valid characters in the IBAN string.
2. **Country Support**: Verifies that the country code is supported and matches the expected IBAN length.
3. **Checksum Calculation**: Moves the initial characters to the end, replaces letters with corresponding digits, and computes the remainder of the resulting number on division by 97.
4. **Checksum Verification**: Ensures that the remainder is equal to 1, indicating a valid IBAN.

## Input and Output

The API endpoint accepts a JSON array of IBAN strings in the following format:

```json
{
  "data": [
    "IBAN1",
    "IBAN2",
    ...
  ]
}
```

It returns a JSON array of objects, each containing the original IBAN string, validation status, bank name (if available), and any errors encountered during validation:

```json
[
  {
    "iban": "IBAN1",
    "valid": true,
    "bank_name": "BankName1",
    "error": "ok"
  },
  {
    "iban": "IBAN2",
    "valid": false,
    "bank_name": "unknown",
    "error": "checksum failure"
  },
  ...
]
```

