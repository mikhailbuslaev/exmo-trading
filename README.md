# exmo-trading
<p>the application allows you to automate trading on the (https://exmo.me/)</p>

## installation

### via git

<p>you need to have installed golang, git.</p>

    git clone https://github.com/mikhailbuslaev/exmo-trading.git

### via docker

<p>you need to have installed docker.</p>

    docker run -it exmo-trading

## run

    cd exmo-trading
    go run exmo-trader.go

## run at server

    cd exmo-trading
    nohup go run exmo-trader.go &

## setting

<p>by default, the application starts 2 data processors and one trade 
for each of these processors, but you can add as many data processors and traders as you like.</p>

## example: add new trader

### 1 stage - make new trader config

    cd exmo-trading/configs/trader-configs
    mkfile trader-config.yaml

<p>trader-config.yaml example:</p>

    TradingPair : "BTC_USDT"
    TradesFile  : "/cache/trades/5min-btc-usdt-trades.csv"
    TradesHistoryFile : "/cache/trades-history/5min-btc-usdt-trades-history.csv"
    CandlesFile       : "/cache/candles/5min-btc-usdt-candles.csv"
    CandlesFileVolume : 250
    TradesFileVolume  : 10
    StopLimitPercent  : 0.015
    TradeAmount : 5
    TraderTimeout : 30

><p>TradingPair examples : "ETH_USDT", "BTC_ETH".</p>
>
><p>TradesFile is a file that will store open trades.</p>
>
><p>TradesHistoryFile is a file that will store history of closed trades.</p>
>
><p>CandlesFile is a file that will be analyzed by the trading strategy.</p>
>
><p>CandlesFileVolume is a volume of records(candles) in CandlesFile, default: 250.</p>
>
><p>TradesFileVolume is a maximum volume of open trades in TradesFile, default: 10.</p>
>
><p>StopLimitPercent is a this is a part of the current deal price, from which the limit is set, 0.015 means 1.5%.</p>
>
><p>TradeAmount is a volume of trade, example: 5 means 5 usdt in "ETH_USDT"</p>
>
><p>TradeTimeout is a reload duration of trader</p>
>

### 2 stage - change exmo-trader.go

    cd exmo-trading
<p>add new line to PrepareApp() after app.Traders = append(... and before fmt.Println("prepared datahandlers: ")</p>

    app.Traders = append(app.Traders, PrepareTrader(path, path+
    	"/configs/trader-configs/trader-config.yaml", &strategies.RSI{}))

<p>&strategies.RSI{} and &strategies.BollingerBands{} is 2 available strategies at this stage of project</p>

### 3 stage

<p>save changes, now you have new trade handler.</p>

## example: add new datahandler

<p>if you want trade another timeframe or another pair, you need add new </p>
<p>data handler for this timeframe or pair.</p>

### 1 stage - make new datahandler config

    cd exmo-trading/configs/datahandler-configs
    mkfile datahandler-config.yaml

<p>datahandler-config.yaml example:</p>

    Symbol : "BTC_USDT"
    Resolution : 5
    CandlesFile : "/cache/candles/5min-btc-usdt-candles.csv"
    CandlesVolume : 250
    DataHandlerTimeout : 60

><p>Symbol examples : "ETH_USDT", "BTC_ETH".</p>
>
><p>Resolution means duration of single candle, 5 means 5 minute candles.</p>
>
><p>CandlesFile is a file of candles.</p>
>
><p>CandlesFileVolume is a volume of candles in CandlesFile.</p>
>
><p>DataHandlerTimeout is a is a reload duration of datahandler.</p>
>

### 2 stage - change exmo-trader.go

    cd exmo-trading

<p>add new line to PrepareApp() after app.Datahandlers = append(... and before return &app)</p>

    app.DataHandlers = append(app.DataHandlers, PrepareDataHandler(path, path+
    	"/configs/datahandler-configs/datahandler-config.yaml"))

### 3 stage

<p>save changes, now you have new data handler.</p>
