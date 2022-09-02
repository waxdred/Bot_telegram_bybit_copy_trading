
# Bot bybit

#### Bot for make copy trading from signal telegram and send to Bybit the order
#### * Check the price of currency
#### * Change stop lost when tp
#### * Order can be cancel but not the position yet
# Environment Variables

To run this project, you will need to add the following environment variables to your .env file

## .env file
```bash
API_TELEGRAM = "api of your bot telegram"


# api bybit
API = "bybit api"
API_SECRET = "bybit api_secret"
# need stay in testnet 
URL = "https://api-testnet.bybit.com"


#api telegram app
API_ID = "api_id"
API_HASH = "api_hash"

# your channel id where you have your bot for send to bybirt
MY_CHANNEL = "id channel"

# id of the channel you want listen the signal work with only one for the moment
ID_CHANNEL = "id channel"
```



# Installation 

### Please setting your .env
#### Run app
```bash
  make
```

#### Stop app
```bash
  make down
```
## Signal exemple

![Screenshot](asset/signal.png)

## Authors

- [@waxdred](https://www.github.com/waxdred)

