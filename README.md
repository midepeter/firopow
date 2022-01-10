<<<<<<< HEAD
# Full FiroPOW verification algorithm implementation 

The FIROPOW is a fork of PROGPOW consensus algorithm(basic PROGPOW with a twist). The twist happens with tweaking of few parameters.   

# The Flow of the Algorithm
At the beginning of the algorithm, we use a keccak to hash header and nonce of the current block to create a seed. We use this seed to generate the initial data for a 512 bytes wide "mix". We repeatedly fetch random loads from the dag and the cache, perform random math on them and use fnv1a to combine it to the mix. 256 byte wide sequential accesses are used, which increases the efficiency and reduces the overhead on modern GPU's. After that, we combine the mix into a single 256-bit value again using fnv1a. At the end, we use another keccak hash on this single 256-bit value to generate a result.

If the output of this algorithm is below the desired target, then the nonce is valid. Note that the extra application of keccak at the end ensures that there exists an intermediate nonce which can be provided to prove that at least a small amount of work was done; this quick outer PoW verification can be used for anti-DDoS purposes. It also serves to provide statistical assurance that the result is an unbiased, 256 bit number.

# Details of the project
variant
