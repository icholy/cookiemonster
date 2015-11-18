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
These are invoked when a user logs in via wafer.

Example config file (for the wafer server):
```
http://sub1.domain.com/wafer_webhook.php
https://sub2.domain.com:8888/wafer_webhook/
```

Hooks have two responsibilities:

1. Set the provided JWT in the cookie so it's available on that domain.
2. Provision a user account for the user in the JWT if it does not already exist.

## Login page:

After the login page, the user is redirected to a page containing `img` tags pointing to the configured domain.
It provisions user data in applications and it sets the JWT cookie.


``` html
<html>
  <body>
    <h1>Login Successfull</h1>
    <a href="{{ReferrerURL}}">Click here if you're not redirected</a>

    <!-- set jwt on other domain -->
    <script src="{{WebHookURL}}?jwt={{JWT}}"></script>
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
  AuthType Wafer
  AuthName Dev_Group_Members
  AuthWaferServer http://wafer.domain.com/
  AuthWaferRedirect http://apache.domain.com/myapp
  AuthWaferSigningMethod RS256
  AuthWaferKeyFile pubkey.pem
  AuthWaferAppName MyApp
  Require group dev
</Directory>
```
