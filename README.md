# Microsoft Graph API Golang SDK PoC

This project is a proof-of-concept that intends to establish a connection to Microsoft Graph API through the utilization of a client certificate instead of a client secret. The primary objective is to retrieve a list of user accounts. The Microsoft Graph API SDK is employed to accomplish the set objectives.

## FPX Certificate Bundle

To generate a password-protected Financial Process Exchange (FPX) certificate bundle using OpenSSL, you can follow these steps.

```bash
# use the `req` command to generate both the certificate and private key in one go
openssl req -x509 \
            -sha256 -days 365 \
            -nodes \
            -newkey rsa:4096 \
            -subj "/C=US/ST=CA/O=Azure/CN=myapp" \
            -keyout private.pem -out certificate.pem


# When you run the below command, OpenSSL will prompt you to enter a password for the PFX file.
# Choose a strong password and make sure to remember it,
# as you will need it later within the Azure authentication process
openssl pkcs12 -export -in certificate.pem -inkey private.pem -out bundle.fpx
```

That's it! You now have a password-protected FPX certificate bundle that can be used to authenticate with Microsoft Graph API.

## References

- [Quickstart: Register an application with the Microsoft identity platform](https://learn.microsoft.com/en-us/azure/active-directory/develop/quickstart-register-app)
- [Microsoft Graph SDK for Go](https://github.com/microsoftgraph/msgraph-sdk-go)
- [GoDotEnv](https://github.com/joho/godotenv)
