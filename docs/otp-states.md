stateDiagram-v2
    [*] --> Created

    Created --> Sent: SMS sent

    Sent --> Verified: Valid code
    Sent --> Failed: Invalid code
    Sent --> Expired: Expiration time reached

    Failed --> Failed: Retry invalid code
    Failed --> Blocked: Max attempts reached
    Failed --> Expired: Expiration time reached

    Verified --> [*]
    Blocked --> [*]
    Expired --> [*]