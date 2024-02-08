# Unix ransomware in go

>This project is still in development. For educational purposes.

Ransomware is a type of malware that locks a victim's data or device and demands a ransom to restore access. This is a simple implementation of ransomware written in golang with concurrent full system encryption using AES algorithm in CTR mode.

## How it works

To begin, the program generates a cryptographically secure AES key. It then uses the current user's home directory as the root folder to initiate the encryption of files. 

Encrypted files are temporarily stored in a newly created folder within the Temp directory during the encryption process. This is designed to prevent suspicion, as the victim is unlikely to notice files appearing and disappearing abruptly. Once the encryption is finished, the encrypted files replace their unencrypted counterparts. 

Finally, the AES key itself is encrypted using the RSA algorithm to secure transmission to the server. Server stores encrypted AES key with the victim's ID. In the end it drops a ransomnote on user's desktop with victim's ID. Generally ransom note contains an explanation what happened and contacts for getting decryption key for payment.

>Please note that RSA is an asymmetric encryption algorithm, meaning it operates with two separate keys: a public key and a private key. Typically, the public key is hardcoded into the program because it needs to be distributed to others for encryption purposes. However, in my program, both the public and private keys are hardcoded. This approach allows for easy access to the private key, which is necessary for decryption later on.

Pretty much it already does key ransomware features but it's still far from being complete.

## TODO

- C&C server
- network traverse
- shadow copy deletion
- flags
- proper error handling 
