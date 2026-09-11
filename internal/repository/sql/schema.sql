CREATE TABLE secrets (
    id VARCHAR(50) PRIMARY KEY, 
    
    encrypted_text TEXT NOT NULL,
    
    is_client_encrypted BOOLEAN NOT NULL DEFAULT FALSE,
    
    password_hash VARCHAR(255), 
    
    notify_email VARCHAR(255),
    
    file_path VARCHAR(500),
    
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);