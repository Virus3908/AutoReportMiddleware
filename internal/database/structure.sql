CREATE EXTENSION IF NOT EXISTS "pgcrypto"; -- Включает поддержку UUID

CREATE TABLE Conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_name VARCHAR(255) NOT NULL,
    file_url VARCHAR(255),
    status INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ConversationsUsers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES Users(id) ON DELETE CASCADE,
    conversation_id UUID REFERENCES Conversations(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Convert (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversations_id UUID REFERENCES Conversations(id) ON DELETE CASCADE,
    file_url VARCHAR(255),
    task_id UUID,  -- UUID передается вручную
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Diarize (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID REFERENCES Conversations(id) ON DELETE CASCADE,
    start_time TIME NOT NULL,
    end_time INTEGER NOT NULL,
    speaker UUID REFERENCES ConversationsUsers(id) ON DELETE CASCADE,
    task_id UUID,  -- UUID передается вручную
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Transcribe (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID REFERENCES Conversations(id) ON DELETE CASCADE,
    segment_id UUID REFERENCES Diarize(id) ON DELETE CASCADE,
    transcription TEXT NOT NULL,
    task_id UUID,  -- UUID передается вручную
    complete BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Report (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID REFERENCES Conversations(id) ON DELETE CASCADE,
    report TEXT NOT NULL,
    promt_id UUID REFERENCES Promts(id) ON DELETE CASCADE,
    task_id UUID,  -- UUID передается вручную
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Promts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    promt TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);