TransIP API
==================
Small library in Golang that implements:
* DomainService/getDomainNames
* DomainService/getInfo

If you need other methods on the TransIP API be free to fork my
code and pull request once you have it added.

Why use this library instead of writing own:
* 'Correctly' implementing __nonce because there's no uniqid in Golang;
* Correctly implementing signature, RSA Private key signing of cookie signature (unittests included);
* Tokenized XML-parser to keep data structures free of clutter (stripping SOAP envelopes);

Note
=======
As far as I know Golang/crypto doesn't support loading PKCS8-certificates as RSA.PrivateKey. So if you
want to start using this code please convert your privatekey from PKCS8 to PKCS1 with OpenSSL:
```
openssl rsa -in privkeyfromtransip.pem -out privkey.pem
```
