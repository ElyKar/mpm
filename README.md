# mpm
My Password Manager

mpm is a small CLI project that I use to handle all of my passwords. A simple CLI password manager.

Actual version is 0.1 and provides the core functionalities as well as cryptographically secure storage: No plain-text password or some personal crypto involved, only community-approved algorithms (bcrypt, AES-256, Go's crypto/rand for randomness).

Each tagged version can be considered stable (at least I do). I have no release plan for the project, only some features that I'd like to add. I'll implement them when I want to, when I feel particularly motivated, or when I'm in absolute need.

# Usage

As often (always ?) mpm lets you order your passwords in sections and gicing them identifiers to easily create and fetch them. All of your data is stored (encrypted) under "$HOME/.mpm".

```
mpm is a CLI password manager made to handle all of your passwords.
You can customise each password by choosing which characters it may contain and its total length.

Usage:
  mpm [command]

Available Commands:
  add         Generates a new password for the section and name
  change      Change the master password
  get         Copy a password to your clipboard
  init        Initialize an empty store for mpm
  list        List the sections and passwords stored

Use "mpm [command] --help" for more information about a command.
```

Well, documentation says it all, you can create a storage, add passwords in it (sections are lazily initialized), copy it to clipboard, list sections, passwords of a section, or all the content of the storage (password do not appear, only the name you gace them) and finally change your master password. KISS to you too.

# Cryptography

A bcrypt hash (cost 10) of your master password is stored in your data file and used to check for validity. The secret key used to encrypt your password is derived from your master password using the SAH512\_256 hash function. This provides a 32 byte strings which is used as a secret key to encrypt your passwords using AES-256 in CTR mode. Encrypted password are finally encoded in base64 to make them printable to your data file.

# Libraries

I purposely use very few libraries in the project, here is the full list (Go's standard library not included).

- github.com/spf13/cobra : My preferred tool to generate CLI apps in Go.
- golang.org/x/crypto/bcrypt : The official Go implementation of bcrypt.
- github.com/atotto/clipboard : A small library to be able to Write and Read To/From clipboard.
- github.com/howeyc/gopass : A small library to correctly handle password prompt on terminals (should not be displayed).

And... that's it. As promised: small.

# Future work

Here are some functionalities I'd like to implement in the (maybe VERY distant) future:

- Possibility to delete an entry.
- Authenticate the data file's content (no unwanted modification), one way or another.
- Make a tiny graphical client.
- Ability to synchronize a storage between several machines (would be great and amazing, but requires good design).
