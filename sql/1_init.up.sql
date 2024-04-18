-------------- PROVIDER TABLE ----------------

CREATE TABLE IF NOT EXISTS provider (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(16) NOT NULL
);

-------------- SIM TABLE ----------------

CREATE TABLE IF NOT EXISTS sim (
    id INT AUTO_INCREMENT PRIMARY KEY,
    number VARCHAR(15) NOT NULL,
    provider_id INT NOT NULL,
    is_activated BOOLEAN DEFAULT 1,
    activate_until BIGINT DEFAULT 0,
    is_blocked BOOLEAN DEFAULT 0,

    FOREIGN KEY (provider_id) REFERENCES provider(id)
);


-------------- SERVICE TABLE ----------------

CREATE TABLE IF NOT EXISTS service (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64) NOT NULL
);

-------------- USED TABLE ----------------

CREATE TABLE IF NOT EXISTS used_service (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sim_id INT NOT NULL,
    service_id INT NOT NULL,
    is_blocked BOOLEAN DEFAULT 0,
    blocked_info VARCHAR(64) DEFAULT '',
    FOREIGN KEY (sim_id) REFERENCES sim(id),
    FOREIGN KEY (service_id) REFERENCES service(id)
);
