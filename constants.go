package donutchallenge

// CoinbaseAPIBaseProduction is for use in production
const CoinbaseAPIBaseProduction = "https://api.pro.coinbase.com"

// CoinbaseAPIBaseTest is for use when testing
const CoinbaseAPIBaseTest = "https://api-public.sandbox.pro.coinbase.com"

// CoinbaseOrdersEndpoint is the endpoint for placing an order
const CoinbaseOrdersEndpoint = "/orders"

// Currency is for marking strings as denoting currencies available on Coinbase
type Currency string

const BAT = Currency("BAT")
const BCH = Currency("BCH")
const BTC = Currency("BTC")
const CVC = Currency("CVC")
const DAI = Currency("DAI")
const DNT = Currency("DNT")
const ETC = Currency("ETC")
const ETH = Currency("ETH")
const EUR = Currency("EUR")
const GBP = Currency("GBP")
const GNT = Currency("GNT")
const LOOM = Currency("LOOM")
const LTC = Currency("LTC")
const MANA = Currency("MANA")
const MKR = Currency("MKR")
const USD = Currency("USD")
const USDC = Currency("USDC")
const ZEC = Currency("ZEC")
const ZIL = Currency("ZIL")
const ZRX = Currency("ZRX")

type OrderType string

const LIMIT = OrderType("limit")
const MARKET = OrderType("market")

type OrderSide string

const BUY = OrderSide("buy")
const SELL = OrderSide("sell")

type OrderSTP string

const DC = OrderSTP("dc")
const CO = OrderSTP("co")
const CN = OrderSTP("cn")
const CB = OrderSTP("cb")

type OrderStop string

const LOSS = OrderStop("loss")
const ENTRY = OrderStop("entry")

type OrderTimeInForce string

const GTC = OrderTimeInForce("GTC")
const GTT = OrderTimeInForce("GTT")
const IOC = OrderTimeInForce("IOC")
const FOK = OrderTimeInForce("FOK")

type OrderCancelAfter string

const MIN = OrderCancelAfter("min")
const HOUR = OrderCancelAfter("hour")
const DAY = OrderCancelAfter("day")
