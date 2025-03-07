# Bastion 🛡️

A modern private and secure messenger.

# TODO

- pkg/ratchet: add docs for each function. Add comments for e.g. ratchetReceivingChain and ratchetSendingChain.
- pkg/ratchet: add tests.
- pkg/ratchet: reduce allocations count. For example, reuse slices for HKDF and encryption/decryption. Encrypt/Decrypt to array from stack.
- pkg/ratchet: create benchmarks to increase speed.
