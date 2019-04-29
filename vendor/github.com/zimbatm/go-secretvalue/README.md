# go-secretvalue - Don't send secrets to logs

This Go library doesn't do much except encourage your to mark all your
application secrets properly.

Instead of:
```go
secret := os.GetEnv("OAUTH_TOKEN")
```

Write:
```go
secret := secretvalue.New("oauth-token")
secret.SetString(os.GetEnv("OAUTH_TOKEN"))
os.Unsetenv("OAUTH_TOKEN")
```

By doing so, it will prevent the secrets from going to the logs inadvertedly.

The `secret.String()` function exposes the secret name instead of the value,
which avoids sending these into logs by mistake. This happens a lot, trust me.

## StringFlag

This library can also be used with the stdlib flag library. See
string_flag_test.go for an example.

## Companies that have sent passwords to logs by mistake

Remember these are only publicly known instances.

* Twitter: https://arstechnica.com/information-technology/2018/05/twitter-advises-users-to-reset-passwords-after-bug-posts-passwords-to-internal-log/
* GitHub: https://www.zdnet.com/article/github-says-bug-exposed-account-passwords/
* Facebook: https://www.theverge.com/2019/3/21/18275837/facebook-plain-text-password-storage-hundreds-millions-users
* ...

## Missing features

* Optionally use `mlock(2)` on supported systems to prevent the value from
  going to swap.

## Other attacks

This library doesn't prevent the value from going to swap disk. Make sure to
disable swap on all of your servers. `swapoff -a`
