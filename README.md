# cookiemonster

> A simple approach to single-sign-on

* Uses LDAP for storing users and passwords.
* Provices login and logout page.
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

Each domain must have a "cookiecutter" which sets cookies based on its GET parameters.
This is a Golang template file where `{{jwt}}` gets replaced by the JWT.

```
sub1.domain.com/cookiecutter.php?jwt={{jwt}}
sub2.domain.com:8888/cookiecutter/?jwt={{jwt}}
```

## Login page:

After the login page, the user is redirected to a page containing `img` tags pointing to the configured domains.

``` html
<html>
  <body>
    <h1>Login Successfull</h1>
    
    <!-- set cookies on other domains -->
    <img src="sub1.domain.com/cookiecutter.php?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEyMyIsIm5hbWUiOiJJbGlhIENob2x5IiwiZ3JvdXBzIjpbImFkbWluIiwiZGV2Il19.JbD8pOZbBz5GOkfLakAisWvM-V9WMlWO4EUt3z8FEd0" />
    <img src="sub2.domain.com:8888/cookiecutter/?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEyMyIsIm5hbWUiOiJJbGlhIENob2x5IiwiZ3JvdXBzIjpbImFkbWluIiwiZGV2Il19.JbD8pOZbBz5GOkfLakAisWvM-V9WMlWO4EUt3z8FEd0" />
  </body>
</html>
```
