# WIP: Wafer

> A simple approach to single-sign-on

* Uses LDAP for storing users and passwords.
* Provides login and logout page.
* Sets a JWT on all specified domain.

## JWT Data:

This is data that will be accessible to all applications via headers.

``` json
{
  "id": 123,
  "name": "Ilia Choly",
  "groups": ["admin", "dev"]
}
```

## Domain Config File:

Each domain must have a wafter **hook** which sets cookies based on its GET parameters.
This is a template file where `{{jwt}}` gets replaced by the JWT.

```
sub1.domain.com/wafer_hook.php?jwt={{jwt}}
sub2.domain.com:8888/wafer_hook/?jwt={{jwt}}
```

## Endpoints

Wafter endpoints are web

## Login page:

After the login page, the user is redirected to a page containing `img` tags pointing to the configured domains. It provisions user data in applications and it sets the JWT cookie.

``` html
<html>
  <body>
    <h1>Login Successfull</h1>
    
    <!-- set cookies on other domains -->
    <img src="sub1.domain.com/wafer_hook.php?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEyMyIsIm5hbWUiOiJJbGlhIENob2x5IiwiZ3JvdXBzIjpbImFkbWluIiwiZGV2Il19.JbD8pOZbBz5GOkfLakAisWvM-V9WMlWO4EUt3z8FEd0" />
    <img src="sub2.domain.com:8888/wafer_hook/?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEyMyIsIm5hbWUiOiJJbGlhIENob2x5IiwiZ3JvdXBzIjpbImFkbWluIiwiZGV2Il19.JbD8pOZbBz5GOkfLakAisWvM-V9WMlWO4EUt3z8FEd0" />
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
