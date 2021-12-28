# carousel

Carousel is a modern IRC bouncer written in Go. It's not production ready in the slightest, so take everything with a grain of salt.

## Features

- [x] Connect to multiple networks
- [x] SSL connections
- [x] Buffer messages when disconnect
- [ ] External state store
- [ ] Remote CLI connection

## Configuration

First, add a user. This password will be hashed and stored to authenticate with Carousel, not any network.

``` bash
crsl add user --username user --password pass
```

Then, add a network to the user.
``` bash
crsl add network libera --address irc.libera.chat --port 6667 --user user
```

Finally, set your ident for that network. Only the nickname is required; the password will default empty and the username / realname will match your nick.

``` bash
crls set ident --user user \
    --network libera
    --nickname nick
    --username user
    --realname real
    --password pass
```

You can also edit the configuration directly. It will look something like this:

``` json
{
    "Server": {
        "URI": "127.0.0.1:6667",
        "Verbose": true,
    },
    "Users": [
        {
            "Username": "user",
            "Password": "$2a$04$itD.ZqpW8spWradGqybbAu8asdfF8BrOfGmApUSJcPxo1e.v7A3AYp6",
            "Networks": [
                {
                    "Name": "libera",
                    "URI": "irc.libera.chat:6667",
                    "Ident": {
                        "Username": "user",
                        "Nickname": "nick",
                        "Realname": "name",
                        "Password": "pass"
                    },
                    "Channels": []
                }
            ]
        }
    ]
}
```

## Setup

I've only tested with Weechat so far. To add the Carousel server...

```
/server add crsl 127.0.0.1/6667 -username=user/network -password=pass
```
