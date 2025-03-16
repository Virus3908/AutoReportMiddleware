CREATE EXTENSION IF NOT EXISTS "pgcrypto"; -- Включает поддержку UUID

CREATE TABLE Conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_name VARCHAR(255) NOT NULL,
    file_url VARCHAR(255),
    status INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ConversationsParticipant (
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

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_conversations
BEFORE UPDATE ON Conversations
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_update_participant
BEFORE UPDATE ON participants
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_update_conversations_participant
BEFORE UPDATE ON ConversationsParticipant
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_update_convert
BEFORE UPDATE ON Convert
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_update_diarize
BEFORE UPDATE ON Diarize
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_update_transcribe
BEFORE UPDATE ON Transcribe
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_update_report
BEFORE UPDATE ON Report
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_update_promts
BEFORE UPDATE ON Promts
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();