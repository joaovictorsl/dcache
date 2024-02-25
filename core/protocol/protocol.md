## How the protocol works

- SET Command
    - Index 0 byte is 0
    - Index 1 byte is key length (**_KL_**)
    - Bytes in index range [2, **_KL_** + 1] are the key
    - Byte in index range [**_KL_** + 2, **_KL_** + 5] is the value length **_VL_**
    - Bytes in index range [**_KL_** + 3, **_KL_** + 2 + **_VL_**] are the value
    - The  bytes in index range [**_KL_** + 3 + **_VL_**, **_KL_** + 10 + **_VL_**] are the expiration time

- GET Command
    - Index 0 byte is 1
    - Index 1 byte is key length **_KL_**
    - Bytes in index range [2, **_KL_** + 1] are the key

- HAS Command
    - Index 0 byte is 3
    - Index 1 byte is key length **_KL_**
    - Bytes in index range [2, **_KL_** + 1] are the key

- DELETE Command
    - Index 0 byte is 4
    - Index 1 byte is key length **_KL_**
    - Bytes in index range [2, **_KL_** + 1] are the key
