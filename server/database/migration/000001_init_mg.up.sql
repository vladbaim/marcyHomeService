CREATE TABLE IF NOT EXISTS sensor_data(
   id SERIAL PRIMARY KEY,
   version INT NOT NULL,
   position VARCHAR (50) NOT NULL,
   humidity FLOAT (8) NOT NULL,
   temperature FLOAT (8) NOT NULL,
   carbon_dioxide INT NOT NULL,
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS telegram_chat(
   id SERIAL PRIMARY KEY,
   name VARCHAR (50) NOT NULL,
   chat_id UUID NOT NULL
);