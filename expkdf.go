package gostikkit

//crypto-js uses openssl's expKDF by default
func evpKDF(passowrd []byte, keySize int, ivSize int, salt []byte, iterations int, hNew  func()hash.Hash, ) (resultKey []byte, resultIv []byte){

    targetKeySize := keySize + ivSize
    derivedBytes := make([]byte,targetKeySize * 4)
    numberOfDerivedWords := 0;

    block = []byte{};

		hasher := hNew()
    for numberOfDerivedWords < targetKeySize {
        if (block != null) {
            hasher.update(block);
        }
        hasher.update(password);
        block = hasher.digest(salt);
        hasher.reset();

        // Iterations
        for (int i = 1; i < iterations; i++) {
            block = hasher.digest(block);
            hasher.reset();
        }

        System.arraycopy(block, 0, derivedBytes, numberOfDerivedWords * 4,
                Math.min(block.length, (targetKeySize - numberOfDerivedWords) * 4));

        numberOfDerivedWords += block.length/4;
    }

    System.arraycopy(derivedBytes, 0, resultKey, 0, keySize * 4);
    System.arraycopy(derivedBytes, keySize * 4, resultIv, 0, ivSize * 4);

    return derivedBytes; // key + iv
}
