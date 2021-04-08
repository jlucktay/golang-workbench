# Google Sign-In for Websites

## Google Cloud - Cloud Functions

- [Background Functions](https://cloud.google.com/functions/docs/writing/background)
- [HTTP Triggers](https://cloud.google.com/functions/docs/calling/http)
- [Authenticating](https://cloud.google.com/functions/docs/securing/authenticating)
- [Events and Triggers](https://cloud.google.com/functions/docs/concepts/events-triggers)

## [Add Google Sign-In to Your Web App](https://developers.google.com/identity/sign-in/web/)

- [Integrate Google Sign-In](https://developers.google.com/identity/sign-in/web/sign-in)
- [Authenticate with a Backend Server](https://developers.google.com/identity/sign-in/web/backend-auth)

## Some JWT Things

### Header: Algorithm & Token Type

```json
{
  "alg": "RS256",
  "kid": "40 characters, hexadecimal",
  "typ": "JWT"
}
```

### Payload

```json
{
  "iss": "accounts.google.com",
  "azp": "<OAuth app ID>.apps.googleusercontent.com",
  "aud": "<OAuth app ID>.apps.googleusercontent.com",
  "sub": "21 digit integer",
  "email": "an email address",
  "email_verified": true,
  "at_hash": "string of length 22",
  "name": "a full name",
  "picture": "a URL to a profile picture",
  "given_name": "a first name",
  "family_name": "a last name",
  "locale": "en-GB",
  "iat": 1617631086,
  "exp": 1617634686,
  "jti": "40 characters, hexadecimal"
}
```

## Security considerations

- <https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html>
- cookies **must** be HTTPS only (`Secure` flag) with `SameSite` set to strict
