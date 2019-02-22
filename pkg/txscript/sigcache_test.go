package txscript

import (
	"crypto/rand"
	"testing"

	"git.parallelcoin.io/pod/pkg/chaincfg/chainhash"
	"git.parallelcoin.io/pod/pkg/ec"
)

// genRandomSig returns a random message, a signature of the message under the public key and the public key. This function is used to generate randomized test data.
func genRandomSig() (*chainhash.Hash, *ec.Signature, *ec.PublicKey, error) {

	privKey, err := ec.NewPrivateKey(ec.S256())
	if err != nil {
		return nil, nil, nil, err
	}
	var msgHash chainhash.Hash
	if _, err := rand.Read(msgHash[:]); err != nil {
		return nil, nil, nil, err
	}
	sig, err := privKey.Sign(msgHash[:])
	if err != nil {
		return nil, nil, nil, err
	}
	return &msgHash, sig, privKey.PubKey(), nil
}

// TestSigCacheAddExists tests the ability to add, and later check the existence of a signature triplet in the signature cache.
func TestSigCacheAddExists(
	t *testing.T) {

	sigCache := NewSigCache(200)

	// Generate a random sigCache entry triplet.
	msg1, sig1, key1, err := genRandomSig()
	if err != nil {
		t.Errorf("unable to generate random signature test data")
	}

	// Add the triplet to the signature cache.
	sigCache.Add(*msg1, sig1, key1)

	// The previously added triplet should now be found within the sigcache.
	sig1Copy, _ := ec.ParseSignature(sig1.Serialize(), ec.S256())
	key1Copy, _ := ec.ParsePubKey(key1.SerializeCompressed(), ec.S256())
	if !sigCache.Exists(*msg1, sig1Copy, key1Copy) {

		t.Errorf("previously added item not found in signature cache")
	}
}

// TestSigCacheAddEvictEntry tests the eviction case where a new signature triplet is added to a full signature cache which should trigger randomized eviction, followed by adding the new element to the cache.
func TestSigCacheAddEvictEntry(
	t *testing.T) {


	// Create a sigcache that can hold up to 100 entries.
	sigCacheSize := uint(100)
	sigCache := NewSigCache(sigCacheSize)

	// Fill the sigcache up with some random sig triplets.
	for i := uint(0); i < sigCacheSize; i++ {
		msg, sig, key, err := genRandomSig()
		if err != nil {
			t.Fatalf("unable to generate random signature test data")
		}
		sigCache.Add(*msg, sig, key)
		sigCopy, _ := ec.ParseSignature(sig.Serialize(), ec.S256())
		keyCopy, _ := ec.ParsePubKey(key.SerializeCompressed(), ec.S256())
		if !sigCache.Exists(*msg, sigCopy, keyCopy) {

			t.Errorf("previously added item not found in signature" +
				"cache")
		}
	}

	// The sigcache should now have sigCacheSize entries within it.
	if uint(len(sigCache.validSigs)) != sigCacheSize {
		t.Fatalf("sigcache should now have %v entries, instead it has %v",
			sigCacheSize, len(sigCache.validSigs))
	}

	// Add a new entry, this should cause eviction of a randomly chosen previous entry.
	msgNew, sigNew, keyNew, err := genRandomSig()
	if err != nil {
		t.Fatalf("unable to generate random signature test data")
	}
	sigCache.Add(*msgNew, sigNew, keyNew)

	// The sigcache should still have sigCache entries.
	if uint(len(sigCache.validSigs)) != sigCacheSize {
		t.Fatalf("sigcache should now have %v entries, instead it has %v",
			sigCacheSize, len(sigCache.validSigs))
	}

	// The entry added above should be found within the sigcache.
	sigNewCopy, _ := ec.ParseSignature(sigNew.Serialize(), ec.S256())
	keyNewCopy, _ := ec.ParsePubKey(keyNew.SerializeCompressed(), ec.S256())
	if !sigCache.Exists(*msgNew, sigNewCopy, keyNewCopy) {

		t.Fatalf("previously added item not found in signature cache")
	}
}

// TestSigCacheAddMaxEntriesZeroOrNegative tests that if a sigCache is created with a max size <= 0, then no entries are added to the sigcache at all.
func TestSigCacheAddMaxEntriesZeroOrNegative(
	t *testing.T) {


	// Create a sigcache that can hold up to 0 entries.
	sigCache := NewSigCache(0)

	// Generate a random sigCache entry triplet.
	msg1, sig1, key1, err := genRandomSig()
	if err != nil {
		t.Errorf("unable to generate random signature test data")
	}

	// Add the triplet to the signature cache.
	sigCache.Add(*msg1, sig1, key1)

	// The generated triplet should not be found.
	sig1Copy, _ := ec.ParseSignature(sig1.Serialize(), ec.S256())
	key1Copy, _ := ec.ParsePubKey(key1.SerializeCompressed(), ec.S256())
	if sigCache.Exists(*msg1, sig1Copy, key1Copy) {

		t.Errorf("previously added signature found in sigcache, but" +
			"shouldn't have been")
	}

	// There shouldn't be any entries in the sigCache.
	if len(sigCache.validSigs) != 0 {
		t.Errorf("%v items found in sigcache, no items should have"+
			"been added", len(sigCache.validSigs))
	}
}
