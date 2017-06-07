# The signedauth middleware
## Purpose
Allows to protect routes with a signature based authentication.
## Features
Quite customizable, including the following custom settings.
* Hash used for signature (e.g. SHA-1), cf `SignedAuthManager.HashFunction`.
* Authorization header prefix (e.g. SAUTH), cf `SignedAuthManager.AuthHeaderPrefix`.
* Access key to secret key logic (e.g. hardcoded strings or database connection), cf `SignedAuthManager.SecretKey`.
* Additional request verifications apart from just the provided header (e.g. a Date header whose value must be in a given format and represent a recent time), cf `SignedAuthManager.SecretKey`.
* Data extraction for HMAC signature (e.g. date header on the first line, first four characters of the body on the second line), cf. `SignedAuthManager.DataToSign`.
* Allow unsigned requests, so they can be intercepted by another middleware for example, cf. `SignedAuthManager.AuthHeaderRequired`.
* Context key and value which can be used in the rest of the calls, cf. `SignedAuthManager.ContextKey` and cf. `SignedAuthManager.ContextValue`.

## Examples
Refer to the [tests](./signatureauth_test.go) and the [example](./example/) directory.
