### The Empty String Public Key

The address `1HT7xU2Ngenf7D4yocz2SAcnNLW7rK8d4E` is a Bitcoin address derived from an empty `""` string as its public key. Wallets with bugs, including Bitcoin Core[^1], have at times attempted to derive an address using `ripemd160(sha256(""))`, which producees `1HT7xU2Ngenf7D4yocz2SAcnNLW7rK8d4E`.

As it is impossible for any private key to result in an empty public key, this address is provably unspendable.

[^1]: https://github.com/bitcoin/bitcoin/issues/445
