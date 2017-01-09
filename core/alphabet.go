package core

import (
	"crypto/rand"
	"math/big"
)

// An Alphabet abstract the character set from which the password will be randomly created.
type Alphabet struct {
	// String to display on the prompt
	Display string
	// Entire set of characters
	choices string
}

// Alphas contains all the pre-defined alphabets.
var Alphas []Alphabet = []Alphabet{
	// Well, some websites require secret digicodes. Don't look at me like that.
	Alphabet{
		"[0-9]",
		"0123456789",
	},
	// Useless, but provides a nicer display.
	Alphabet{
		"[a-zA-Z]",
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
	},
	// Should not be used. Though, some very legacy apps may still have trouble with special chars (shame on them)
	Alphabet{
		"[a-zA-Z0-9]",
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
	},
	// Finally some special chars, restreint version
	Alphabet{
		"[a-zA-Z0-9!&()*+,-./?[]~]",
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!&()*+,-./?[]~",
	},
	// Full (or almost full) version, though may be a bit too much
	Alphabet{
		"[a-zA-Z0-9 !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~]",
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~]",
	},
}

// GenPassword creates a random password of given length, using Golang's cryptographically secure PRNG (which is just a wrapper around your OS's cryptographically secure PRNG, so if you find a problem of randomness in it, that's something you can be proud of).
func (a Alphabet) GenPassword(n int) string {
	max := big.NewInt(int64(len(a.choices)))
	var next *big.Int

	pass := make([]byte, n)
	for i := 0; i < n; i++ {
		next, _ = rand.Int(rand.Reader, max)
		pass[i] = a.choices[next.Int64()]
	}

	return string(pass)
}
