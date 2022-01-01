<<<<<<< HEAD
# Full FiroPOW verification algorithm implementation 

The FIROPOW is a PROGPOW variant(basic PROGPOW with a twist). The twist happens with tweaking of few parameters.   

Added sequence of random math instructions and random cache reads that are merged into a much larger mixed state
KISS 99 calculations
-TODOs
1. Build a ProgPoW package
   To build the progpow package: There exist a seed by scanning through the block headers up until that point
   From the seed one can compute 16mb pseudorandom cache Light client store the cache
   From the dataset we can generate the 2gb dataset with the property the each item from the cache. Full clients and miners store the dataset.The dataset grows linearly with time
   The first bytes of the dag acts also 16kb cache set with 
   
2. Make appropriate twist and tune init FiroPoW
3. Make appropriate tuning of the values
4. Testing (Build a simple blockchain)


# The Flow of the Algorithm
At the beginning of the algorithm, we use a keccak to hash header and nonce of the current block to create a seed. We use this seed to generate the initial data for a 512 bytes wide "mix". We repeatedly fetch random loads from the dag and the cache, perform random math on them and use fnv1a to combine it to the mix. 256 byte wide sequential accesses are used, which increases the efficiency and reduces the overhead on modern GPU's. After that, we combine the mix into a single 256-bit value again using fnv1a. At the end, we use another keccak hash on this single 256-bit value to generate a result.

If the output of this algorithm is below the desired target, then the nonce is valid. Note that the extra application of keccak at the end ensures that there exists an intermediate nonce which can be provided to prove that at least a small amount of work was done; this quick outer PoW verification can be used for anti-DDoS purposes. It also serves to provide statistical assurance that the result is an unbiased, 256 bit number.
