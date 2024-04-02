minid query simpleswap params
minid query simpleswap pool
minid query simpleswap coin-reserve ETH
minid query simpleswap coin-reserves
# To Add liquidity
minid tx simpleswap add-liquidity mini1q8zckznq0ck8h9khp92tqcttgvhgu42p6n70h6 10000000ETH --from alice --keyring-backend test

minid tx simpleswap add-liquidity mini1jg80x0c6hlnp7yv6hkjylyvq2m9604tqyvy670 10000000WETH --from bob --keyring-backend test

# To Swap Liquidity
minid tx simpleswap swap-liquidity mini1wxd8ktepsu5tnfhh06dm88mdjr0k3jp0e7ww0f 5000000WETH 5000000ETH --from traderB --keyring-backend test

minid tx simpleswap swap-liquidity mini1hvcnhsgdrn3qvx9rs6ev6exknamu6xz0zn7vjw 5000000ETH 5000000WETH --from traderA --keyring-backend test

minid tx simpleswap swap-liquidity mini1hnxfr47u3nltq5t8ffmh5w8dpmxcv83n9fs2aa 5000000stkETH 5000000ETH --from traderC --keyring-backend test
# To Remove Liquidity
minid tx simpleswap remove-liquidity mini17pzs5k8pwejad0rsj0j4lm7dzjqdmjvtec2uzm 5000000ETH --from alice --keyring-backend test
minid tx simpleswap remove-liquidity mini1jg80x0c6hlnp7yv6hkjylyvq2m9604tqyvy670 10000000WETH --from bob --keyring-backend test
# To Check balance of Alice and Bob
minid q bank balances $(minid keys show alice -a)
minid q bank balances $(minid keys show bob -a)

# To CHeck balance for an account
minid q bank balances mini1cl2fqttaw3ps9zn3rq7w3srn8yn6gfpwvshvjq

# To Check balance of Alice in the simpleswap module
minid q simpleswap liquidity-provider mini17pzs5k8pwejad0rsj0j4lm7dzjqdmjvtec2uzm  

# To Check balance of Bob in the simpleswap module
minid q simpleswap liquidity-provider mini1jg80x0c6hlnp7yv6hkjylyvq2m9604tqyvy670