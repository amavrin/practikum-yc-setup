# practikum-yc-setup
Set up student's yc profile on practikum VM

## General info
Command connects to VM, **removes** existing yc profile and set up new one for federated user.

## Prerequisites

1. Prepare your VM ID and Federation ID
2. Set up ssh authorization: generate keys with `ssh-keygen`, place public key in `.ssh/authorized_keys` on VM
3. Authorize in local browser with federation user
4. Execute the command:

   Let `VM_ID` is `std-int-005-10` and `FED_ID=bpfpfctkh7focc85u9sq`
   ```
   go run ./cmd/main.go \
       -server student@${VM_ID}.praktikum-services.tech \
       -federation-id=${FED_ID}
   ```
   If you prefer to use prebuild binaries, download one for your architecture
   from Releases page, and use it:
   ```
   ./vm-profile-setup-linux-amd64 \
       -server student@${VM_ID}.praktikum-services.tech \
       -federation-id=${FED_ID}
   ```
