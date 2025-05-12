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

## Example session

```
$ vm-profile-setup-darwin-arm64 -server student@std-int-005-10.praktikum-services.tech -federation-id=bpfpfctkh7focc85u9sq

2025/05/12 10:23:17 running for  student@std-int-005-10.praktikum-services.tech
2025/05/12 10:23:17 setting up xdg-open...
2025/05/12 10:23:19 removing yc profile...
2025/05/12 10:23:20 installing yc...
2025/05/12 10:23:24 configuring YC profile...
2025/05/12 10:23:25 sending  yc init --federation-id=bpfpfctkh7focc85u9sq  command...
2025/05/12 10:23:25 Waiting for response...
2025/05/12 10:23:29 sending Enter...
2025/05/12 10:23:29 waiting for URL...
2025/05/12 10:23:29 setting up SSH tunnel...
2025/05/12 10:23:30 opening browser...
2025/05/12 10:23:30 configuring profile settings...
2025/05/12 10:23:32 using 1-st available folder...
2025/05/12 10:23:32 not choosing default zone...
2025/05/12 10:23:32 waiting for prompt...
2025/05/12 10:23:33 check YC working...
2025/05/12 10:23:34
+----------------------+-----------------------------+----------------------+--------+
|          ID          |            NAME             |   ORGANIZATION ID    | LABELS |
+----------------------+-----------------------------+----------------------+--------+
| b1g3jddf4nv5e9okle7p | cloud-praktikumdevopscourse | bpfuefc0l27kjb4r9mt3 |        |
+----------------------+-----------------------------+----------------------+--------+
2025/05/12 10:23:34 successfully configured YC profile
```
