# Bastion 🛡️

A modern private and secure messenger.

# TODO

- pkg/ratchet: add docs for each public function. Do not forget to doc that Clone methods can be called with nil value.
- pkg/ratchet: add tests.
- pkg/ratchet: reduce allocations count. For example, reuse slices for HKDF and encryption/decryption. Encrypt/Decrypt to array from stack.
- pkg/ratchet: create benchmarks to increase speed.
