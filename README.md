# WIP: Wafer

> A barebones approach to single-sign-on

* Provides login and logout page.
* Uses LDAP for storing users and passwords.
* Sets a [JWT](http://jwt.io/) on all specified domain.

## JWT Data:

This is data that will be accessible to all applications via headers.

``` json
{
  "id": 123,
  "name": "Ilia Choly",
  "groups": ["admin", "dev"]
}
```

## WebHooks:

Each application can be configured with a wafer webhook.
These are invoked every time a user logs in via wafer.

Example config file (for the wafer server):
```
sub1.domain.com/wafer_webhook.php?jwt={{jwt}}
sub2.domain.com:8888/wafer_webhook/?jwt={{jwt}}
```

Hooks have two responsibilities:

1. Set the provided JWT in the cookie so it's available on that domain.
2. Provision a user account for the user in the JWT if it does not already exist.

## Login page:

After the login page, the user is redirected to a page containing `img` tags pointing to the configured domains. It provisions user data in applications and it sets the JWT cookie.

``` html
<html>
  <body>
    <h1>Login Successfull</h1>
    
    <!-- set cookies on other domains -->
    <img src="sub1.domain.com/wafer_webhook.php?jwt=xxxxx.yyyyy.zzzzz" />
    <img src="sub2.domain.com:8888/wafer_webhook/?jwt=xxxxx.yyyyy.zzzzz" />
  </body>
</html>
```

## Application Integration

1. The JWT is used to identify which user is logged in.
2. If there is no JWT, the application redirects to the `wafer` server.
3. After login, the `wafer` server invokes the application's `/water_hook` route.
4. The user is then redirected back to the application.

## Apache Integration

I'll need to write an apache module which authenticates against JWT. Example:

``` apache
<Directory "/www/dev">
  AuthName "dev group members"
  AuthWaferServer http://wafer.domain.com/
  AuthType Wafer
  Require group dev
</Directory>
```
