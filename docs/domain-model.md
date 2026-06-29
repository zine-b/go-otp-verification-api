erDiagram
    OTP {
        string id
        string phone_number
        string code_hash
        string status
        int attempts
        int max_attempts
        datetime expires_at
        datetime created_at
        datetime verified_at
    }

    RATE_LIMIT {
        string key
        int request_count
        datetime window_start
        datetime window_end
    }

    OTP ||--o{ RATE_LIMIT : "controlled_by"