# btcsupply

`btcsupply` tracks and monitors burns on the Bitcoin network. Three categories of burns are tracked:

1. Bugs - BTC that is never minted due to miner bugs, such as the [50 BTC from the Genesis Block](https://burned.money/transaction/4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b).
2. `OP_RETURN` - Non-zero `OP_RETURN` outputs are provably unspendable and are considered burns.
3. Script Burns - BTC that is sent to a script that is likely or provably unspendable, such as Satoshi's addresses or known burn addresses.
