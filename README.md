# Cloudflare DNS Records Updater (cfupdater)

DNS updater for Cloudflare API v4.

## Config example

Use the symbol "%" if you want to automatically retrieve your public IP every time you run this program. Otherwise, specify a certain IP address here.

```
{
    "api_key": "cloudflare-api-key",
    "zones": [
        {
            "name": "example.com",
            "records": [
                {
                    "record": "www",
                    "type": "A",
                    "ip": "%"
                }
            ]
        }
    ]
}
```

## Using

Place `cfupdater` and your `config.json` file in a directory of your choice.

Add the following line to your crontab:

```
*/10 * * * *    cd /opt/cfupdater/ && ./cfupdater
```

This will run cfupdater every 10 minutes.